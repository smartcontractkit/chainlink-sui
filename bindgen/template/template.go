// nolint
package template

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"slices"
	"strings"
	"text/template"

	"github.com/smartcontractkit/chainlink-sui/bindgen/parse"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

type FunctionInfo struct {
	Package    string              `json:"package"`
	Module     string              `json:"module"`
	Name       string              `json:"name"`
	Parameters []FunctionParameter `json:"parameters"`
}

type FunctionParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func ParseFunctionInfo(info ...string) ([]FunctionInfo, error) {
	var result []FunctionInfo
	for _, s := range info {
		var temp []FunctionInfo
		if err := json.Unmarshal([]byte(s), &temp); err != nil {
			return nil, err
		}
		result = append(result, temp...)
	}

	return result, nil
}

func MustParseFunctionInfo(info ...string) []FunctionInfo {
	result, err := ParseFunctionInfo(info...)
	if err != nil {
		panic(err)
	}

	return result
}

//go:embed go.tmpl
var tmpl string

type tmplData struct {
	Package      string
	Module       string
	FunctionInfo string
	Structs      []*tmplStruct
	Funcs        []*tmplFunc
	Imports      []*tmplImport
	Artifact     bind.PackageArtifact
}

type tmplStruct struct {
	Name   string
	Fields []*tmplField
}

type tmplOption struct {
	UnderlyingGoType string
}

type tmplImport struct {
	Path        string // The import path, e.g. github.com/smartcontractkit/chainlink-aptos/path/etc
	PackageName string // The package name to import this import with, e.g. module_ocr3_base
}

type tmplType struct {
	GoType   string
	MoveType string

	Option *tmplOption
	Import *tmplImport // Optional go import to add for this type
}

type tmplField struct {
	Name string
	Type tmplType
}

type tmplFunc struct {
	Name     string
	MoveName string
	Params   []*tmplField
	Returns  []tmplType
}

func Convert(pkg, mod string, structs []parse.Struct, functions []parse.Func, externalStructs []parse.ExternalStruct) (tmplData, error) {
	data := tmplData{
		Package: pkg,
		Module:  mod,
	}
	structMap := make(map[string]parse.Struct)
	importMap := make(map[string]*tmplImport)
	for _, s := range structs {
		out := &tmplStruct{
			Name:   s.Name,
			Fields: nil,
		}
		structMap[s.Name] = s
		data.Structs = append(data.Structs, out)
	}
	for i, s := range data.Structs {
		parsedStruct := structMap[s.Name]
		for _, field := range parsedStruct.Fields {
			goType, err := createGoTypeFromMove(field.Type, structMap, externalStructs)
			if err != nil {
				log.Printf("WARNING: Ignoring unknown type of struct %q: %v\n", s.Name, field.Type)
				continue
			}
			data.Structs[i].Fields = append(data.Structs[i].Fields, &tmplField{
				Type: goType,
				Name: ToUpperCamelCase(field.Name),
			})
			if goType.Import != nil {
				importMap[goType.Import.Path] = goType.Import
			}
		}
	}

	var functionInfos []FunctionInfo

	for _, f := range functions {
		if f.Name == "init_module" {
			continue
		}
		out := &tmplFunc{
			Name:     ToUpperCamelCase(f.Name),
			MoveName: f.Name,
			Params:   nil,
			Returns:  nil,
		}
		functionInfo := FunctionInfo{
			Package:    pkg,
			Module:     mod,
			Name:       f.Name,
			Parameters: nil,
		}
		skip := false
		for _, param := range f.Params {
			// Strip the string "mut" from param.Type
			param.Type = strings.ReplaceAll(param.Type, "&mut", "")
			param.Type = strings.TrimSpace(param.Type)

			param.Type = strings.ReplaceAll(param.Type, "&", "")
			param.Type = strings.TrimSpace(param.Type)

			if param.Type == "TxContext" {
				// Ignore the context parameter
				continue
			}
			// external types aren't supported as parameters, therefore passing no externalStructs
			typ, err := createGoTypeFromMove(param.Type, structMap, externalStructs)
			if err != nil {
				if f.IsEntry {
					panic(fmt.Sprintf("Function %v has unsupported parameter %v, type %v", f.Name, param.Name, param.Type))
				} else {
					log.Printf("WARNING: Ignoring function %v due to unknown parameter type %v: %v\n", f.Name, param.Name, param.Type)
					skip = true

					break
				}
			}
			if typ.Option != nil {
				if f.IsEntry {
					panic(fmt.Sprintf("Function %v has unsupported option::Option parameter %q: %v", f.Name, param.Name, typ.MoveType))
				} else {
					log.Printf("WARNING: Ignoring function %v due to unsupported option::Option parameter %q: %v", f.Name, param.Name, typ.MoveType)
					skip = true

					break
				}
			}
			if typ.Import != nil {
				importMap[typ.Import.Path] = typ.Import
			}
			name := ToLowerCamelCase(param.Name)
			if name == "" {
				panic(fmt.Sprintf("Function %v has unsupported parameter name %v, type %v", f.Name, param.Name, param.Type))
			}
			out.Params = append(out.Params, &tmplField{
				Type: typ,
				Name: name,
			})
			functionInfo.Parameters = append(functionInfo.Parameters, FunctionParameter{
				Name: param.Name,
				Type: typ.MoveType,
			})
		}
		for _, returnType := range f.ReturnTypes {
			typ, err := createGoTypeFromMove(returnType, structMap, externalStructs)
			if err != nil {
				if f.IsView {
					// If the function is a view function and has an unknown return type, panic
					panic(fmt.Sprintf("Function %v has an unknown return type: %v: %v", f.Name, returnType, err))
				} else {
					log.Printf("WARNING: Function %v has an unknown return type: %v", f.Name, returnType)
					// skip = true

					break
				}
			}
			out.Returns = append(out.Returns, typ)
			if typ.Import != nil {
				importMap[typ.Import.Path] = typ.Import
			}
		}
		if f.HasTypeParams {
			typeArgField := &tmplField{
				Type: tmplType{
					GoType:   "string",
					MoveType: "string",
				},
				Name: "typeArgs",
			}
			out.Params = append([]*tmplField{typeArgField}, out.Params...)
		}
		if skip {
			continue
		}
		data.Funcs = append(data.Funcs, out)
		functionInfos = append(functionInfos, functionInfo)
	}
	slices.SortFunc(functionInfos, func(a, b FunctionInfo) int {
		return strings.Compare(a.Name, b.Name)
	})
	marshalledInfo, err := json.Marshal(functionInfos)
	if err != nil {
		return tmplData{}, err
	}
	data.FunctionInfo = string(marshalledInfo)
	for _, v := range importMap {
		data.Imports = append(data.Imports, v)
	}

	return data, nil
}

func Generate(data tmplData) (string, error) {
	funcs := template.FuncMap{
		"toLowerCamel": ToLowerCamelCase,
		"toUpperCamel": ToUpperCamelCase,
		"toJSON":       FormatStructToJSON,
	}

	tpl := template.Must(template.New("").Funcs(funcs).Parse(tmpl))
	buffer := new(bytes.Buffer)
	if err := tpl.Execute(buffer, data); err != nil {
		return "", err
	}
	bb := buffer.Bytes()
	formatted, err := format.Source(bb)
	if err == nil {
		return string(formatted), nil
	}

	return string(bb), nil
}

var UppercaseWords []string

// ToUpperCamelCase converts an under-score string to a camel-case string
func ToUpperCamelCase(input string) string {
	// Remove the first underscore if it exists
	if strings.HasPrefix(input, "_") {
		input = input[1:]
	}
	parts := strings.Split(input, "_")
	for i, s := range parts {
		if len(s) > 0 {
			for _, word := range UppercaseWords {
				if strings.EqualFold(word, s) {
					s = word
				}
			}
			parts[i] = strings.ToUpper(s[:1]) + s[1:]
		}
	}

	return strings.Join(parts, "")
}

func ToLowerCamelCase(input string) string {
	// Remove the first underscore if it exists
	if strings.HasPrefix(input, "_") {
		input = input[1:]
	}
	parts := strings.Split(input, "_")
	for i, s := range parts {
		if len(s) > 0 {
			if i != 0 {
				for _, word := range UppercaseWords {
					if strings.EqualFold(word, s) {
						s = word
					}
				}
				parts[i] = strings.ToUpper(s[:1]) + s[1:]
			}
		}
	}

	param := strings.Join(parts, "")
	if len(param) == 0 { // Give a default name if empty, mostly for `_` named params
		param = "param"
	}
	if param == "c" {
		param = "c_"
	}
	return param
}

// FormatStructToJSON converts any struct into a pretty-printed JSON string.
func FormatStructToJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}

	return string(b), nil
}

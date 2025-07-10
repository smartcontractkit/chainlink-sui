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
	Package  string
	Module   string
	Structs  []*tmplStruct
	Funcs    []*tmplFunc
	Imports  []*tmplImport
	Artifact bind.PackageArtifact
}

func (d *tmplData) BuildStructMap() map[string]*tmplStruct {
	structMap := make(map[string]*tmplStruct)
	for _, s := range d.Structs {
		structMap[s.Name] = s
	}
	return structMap
}

type tmplStruct struct {
	Name   string
	Fields []*tmplField
}

func (s *tmplStruct) NeedsCustomDecoder(allStructs map[string]*tmplStruct) bool {
	for _, field := range s.Fields {
		// TODO: recursively handle address decoding
		if field.Type.MoveType == "address" || field.Type.MoveType == "vector<address>" || field.Type.MoveType == "vector<vector<address>>" {
			return true
		}

		if nestedStruct, ok := allStructs[field.Type.MoveType]; ok {
			if nestedStruct.NeedsCustomDecoder(allStructs) {
				return true
			}
		}
	}
	return false
}

func GetBCSType(field *tmplField, allStructs map[string]*tmplStruct) string {
	// TODO: recursively handle address decoding
	switch field.Type.MoveType {
	case "address":
		return "[32]byte"
	case "vector<address>":
		return "[][32]byte"
	case "vector<vector<address>>":
		return "[][][32]byte"
	default:
		if nestedStruct, ok := allStructs[field.Type.MoveType]; ok {
			if nestedStruct.NeedsCustomDecoder(allStructs) {
				return "bcs" + field.Type.MoveType
			}
		}
		return field.Type.GoType
	}
}

type tmplOption struct {
	UnderlyingGoType string
}

type tmplImport struct {
	Path        string // The import path, e.g. github.com/smartcontractkit/chainlink-aptos/path/etc
	PackageName string // The package name to import this import with, e.g. module_ocr3_base
}

type tmplType struct {
	GoType       string
	MoveType     string
	OriginalType string // original move type with modifiers

	Option *tmplOption
	Import *tmplImport // Optional go import to add for this type
}

type tmplField struct {
	Name string
	Type tmplType
}

type tmplFunc struct {
	Name            string
	MoveName        string
	Params          []*tmplField
	Returns         []tmplType
	IsEntry         bool
	HasReturnValues bool
	HasTypeParams   bool
	TypeParams      []string
}

func (f *tmplFunc) HasSingleReturn() bool {
	return len(f.Returns) == 1
}

func (f *tmplFunc) HasMultipleReturns() bool {
	return len(f.Returns) > 1
}

func (f *tmplFunc) GetSingleReturnGoType() string {
	if !f.HasSingleReturn() {
		return ""
	}
	return f.Returns[0].GoType
}

func (f *tmplFunc) HasGenericReturns() bool {
	for _, ret := range f.Returns {
		if containsGenericTypeParam(ret.MoveType, f.TypeParams) {
			return true
		}
	}
	return false
}

func (f *tmplFunc) GetSingleReturnGoTypeForDevInspect() string {
	if !f.HasSingleReturn() {
		return ""
	}

	if containsGenericTypeParam(f.Returns[0].MoveType, f.TypeParams) {
		return "any"
	}

	return f.Returns[0].GoType
}

func Convert(pkg, mod string, structs []parse.Struct, functions []parse.Func) (tmplData, error) {
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
			goType, err := createGoTypeFromMove(field.Type, structMap)
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
			Name:            ToUpperCamelCase(f.Name),
			MoveName:        f.Name,
			Params:          nil,
			Returns:         nil,
			IsEntry:         f.IsEntry,
			HasReturnValues: len(f.ReturnTypes) > 0,
			HasTypeParams:   f.HasTypeParams,
			TypeParams:      f.TypeParams,
		}
		functionInfo := FunctionInfo{
			Package:    pkg,
			Module:     mod,
			Name:       f.Name,
			Parameters: nil,
		}
		skip := false
		for _, param := range f.Params {
			originalType := param.Type

			// strip modifiers for go type generation
			cleanType := strings.ReplaceAll(param.Type, "&mut", "")
			cleanType = strings.TrimSpace(cleanType)
			cleanType = strings.ReplaceAll(cleanType, "&", "")
			cleanType = strings.TrimSpace(cleanType)

			if cleanType == "TxContext" {
				continue
			}
			typ, err := createGoTypeFromMove(cleanType, structMap)
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
			typ.OriginalType = originalType
			out.Params = append(out.Params, &tmplField{
				Type: typ,
				Name: name,
			})
			functionInfo.Parameters = append(functionInfo.Parameters, FunctionParameter{
				Name: param.Name,
				Type: originalType,
			})
		}
		for _, returnType := range f.ReturnTypes {
			typ, err := createGoTypeFromMove(returnType, structMap)
			if err != nil {
				log.Printf("WARNING: Function %v has an unknown return type: %v", f.Name, returnType)
				// skip = true

				break
			}
			out.Returns = append(out.Returns, typ)
			if typ.Import != nil {
				importMap[typ.Import.Path] = typ.Import
			}
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

	for _, v := range importMap {
		data.Imports = append(data.Imports, v)
	}

	return data, nil
}

func getZeroValue(goType string) string {
	switch goType {
	case "string":
		return `""`
	case "bool":
		return "false"
	case "byte", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64":
		return "0"
	case "*big.Int":
		return "nil"
	case "[]byte":
		return "nil"
	case "any":
		return "nil"
	default:
		if len(goType) > 2 && goType[:2] == "[]" {
			return "nil"
		}
		if len(goType) > 1 && goType[:1] == "*" {
			return "nil"
		}
		return goType + "{}"
	}
}

func Generate(data tmplData) (string, error) {
	structMap := data.BuildStructMap()

	funcs := template.FuncMap{
		"toLowerCamel": ToLowerCamelCase,
		"toUpperCamel": ToUpperCamelCase,
		"getZeroValue": getZeroValue,
		"needsCustomDecoder": func(structName string) bool {
			if s, ok := structMap[structName]; ok {
				return s.NeedsCustomDecoder(structMap)
			}
			return false
		},
		"getBCSType": func(field *tmplField) string {
			return GetBCSType(field, structMap)
		},
		"isNestedStructWithDecoder": func(moveType string) bool {
			if s, ok := structMap[moveType]; ok {
				return s.NeedsCustomDecoder(structMap)
			}
			return false
		},
		"getFullyQualifiedType": func(moveType string, packageName string, moduleName string) string {
			// check if this is a struct defined in this module, else return as-is
			if _, ok := structMap[moveType]; ok {
				return packageName + "::" + moduleName + "::" + moveType
			}
			return moveType
		},
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

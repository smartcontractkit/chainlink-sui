// nolint
package parse

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	tree_sitter_move_on_aptos "github.com/aptos-labs/tree-sitter-move-on-aptos/bindings/go"
	tree_sitter "github.com/smacker/go-tree-sitter"
)

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Func struct {
	IsEntry bool `json:"is_entry"`

	Name          string   `json:"name"`
	Params        []Param  `json:"params"`
	HasTypeParams bool     `json:"has_type_params"`
	TypeParams    []string `json:"type_params"`
	ReturnTypes   []string `json:"return_types"`
}

type Struct struct {
	IsEvent bool

	Name   string
	Fields []Param
}

func ParseModule(module []byte) (pkg string, mod string, err error) {
	// Try regex first since it's more reliable for this simple case
	moduleContent := string(module)

	// Look for "module package::module;" pattern anywhere in the file
	// but prioritize early occurrences (likely to be the actual module declaration)
	re := regexp.MustCompile(`(?m)^\s*module\s+(\w+)::(\w+)\s*;`)
	matches := re.FindStringSubmatch(moduleContent)
	if len(matches) == 3 {
		pkg = matches[1]
		mod = matches[2]
		return pkg, mod, nil
	}

	// Fallback to tree-sitter if regex didn't work
	lang := tree_sitter.NewLanguage(tree_sitter_move_on_aptos.Language())
	n, err := tree_sitter.ParseCtx(context.Background(), module, lang)
	if err != nil {
		return "", "", fmt.Errorf("parsing AST: %w", err)
	}

	query, err := tree_sitter.NewQuery([]byte(`
(module
  (identifier) @package
  "::"
  (identifier) @module
)
	`), lang)

	if err != nil {
		return "", "", fmt.Errorf("creating query: %w", err)
	}

	queryCursor := tree_sitter.NewQueryCursor()
	queryCursor.Exec(query, n)

	// Take the first match - this should be the module declaration
	// since it appears before use statements in the file
	for {
		m, ok := queryCursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range m.Captures {
			switch capture.Index {
			case 0:
				// @package
				pkg = capture.Node.Content(module)
			case 1:
				// @module
				mod = capture.Node.Content(module)
			}
		}

		// Take the first match only
		if pkg != "" && mod != "" {
			break
		}
	}

	return
}

func ParseFunctions(module []byte) ([]Func, error) {
	lang := tree_sitter.NewLanguage(tree_sitter_move_on_aptos.Language())
	n, err := tree_sitter.ParseCtx(context.Background(), module, lang)
	if err != nil {
		return nil, fmt.Errorf("parsing AST: %w", err)
	}

	// query to select all public functions
	queryViewFunctions, err := tree_sitter.NewQuery([]byte(`
(declaration
  (attributes
    (attribute) @attribute
  )?
  (module_member_modifier
  	(visibility) @viz
  )
  (module_member_modifier)? @modifier
  (function_decl
  	name: (identifier) @function_name
	type_parameters: (type_params)? @type_params
    return_type: (type)? @returnType
  ) @function
)
	`), lang)
	if err != nil {
		return nil, fmt.Errorf("error creating query: %w", err)
	}

	// For each function_decl (returned by the previous query), retrieve all parameter names and types
	queryParameters, err := tree_sitter.NewQuery([]byte(`
(function_decl
  name: (identifier)
  (parameters
    (parameter
     variable: (identifier) @parameterName
     (type) @type
    )
  )
)
	`), lang)
	if err != nil {
		return nil, fmt.Errorf("error creating query: %w", err)
	}

	functionCursor := tree_sitter.NewQueryCursor()
	functionCursor.Exec(queryViewFunctions, n)

	var functions []Func
	for {
		m, ok := functionCursor.NextMatch()
		if !ok {
			break
		}
		m = functionCursor.FilterPredicates(m, module)
		if len(m.Captures) == 0 {
			continue
		}
		f := Func{}
		testFunc := false
		for _, capture := range m.Captures {
			switch capture.Index {
			case 0:
				// @attribute
				if capture.Node.Content(module) == "test" || capture.Node.Content(module) == "test_only" {
					testFunc = true
				}
			case 2:
				// @modifier
				if capture.Node.Content(module) == "entry" {
					f.IsEntry = true
				}
			case 3:
				// @function_name
				if strings.Contains(capture.Node.Content(module), "test") {
					testFunc = true
				}
				f.Name = capture.Node.Content(module)
			case 4: // @type_params
				f.HasTypeParams = true
				// type_params is something like < T, U >
				childCount := int(capture.Node.NamedChildCount())
				for i := 0; i < childCount; i++ {
					child := capture.Node.NamedChild(i)
					if child.Type() == "type_param" {
						// first child of type_param is the identifier
						paramIdent := child.NamedChild(0)
						if paramIdent != nil {
							f.TypeParams = append(f.TypeParams, paramIdent.Content(module))
						}
					}
				}
			case 5:
				// @returnType
				switch capture.Node.Child(0).Type() {
				case "tuple_type":
					for i := range capture.Node.Child(0).ChildCount() {
						if capture.Node.Child(0).Child(int(i)).Type() == "type" {
							f.ReturnTypes = append(f.ReturnTypes, capture.Node.Child(0).Child(int(i)).Content(module))
						}
					}
				default:
					f.ReturnTypes = append(f.ReturnTypes, capture.Node.Content(module))
				}
			case 6:
				// @function
				qcParam := tree_sitter.NewQueryCursor()
				qcParam.Exec(queryParameters, capture.Node)
				for {
					match, ok := qcParam.NextMatch()
					if !ok {
						break
					}
					param := Param{}
					for _, queryCapture := range match.Captures {
						switch queryCapture.Index {
						case 0:
							// @parameterName
							param.Name = queryCapture.Node.Content(module)
						case 1:
							// @type
							param.Type = queryCapture.Node.Content(module)
						}
					}
					f.Params = append(f.Params, param)
				}
			}
		}
		if !testFunc {
			functions = append(functions, f)
		}
	}

	return functions, nil
}

func ParseStructs(module []byte) ([]Struct, error) {
	lang := tree_sitter.NewLanguage(tree_sitter_move_on_aptos.Language())
	n, err := tree_sitter.ParseCtx(context.Background(), module, lang)
	if err != nil {
		return nil, fmt.Errorf("parsing AST: %w", err)
	}

	queryStructs, err := tree_sitter.NewQuery([]byte(`
(declaration
  (attributes
    (attribute) @attribute
  )?
  (struct_decl
  	name: (identifier) @name
    (body) @structBody
  )
)
	`), lang)
	if err != nil {
		panic(err)
	}

	queryFields, err := tree_sitter.NewQuery([]byte(`
  (field_annot
	field: (identifier) @fieldName
	(type) @type
  )
	`), lang)
	if err != nil {
		panic(err)
	}

	structsCursor := tree_sitter.NewQueryCursor()
	structsCursor.Exec(queryStructs, n)
	var structs []Struct
	for {
		m, ok := structsCursor.NextMatch()
		if !ok {
			break
		}

		m = structsCursor.FilterPredicates(m, module)
		s := Struct{}
		for _, capture := range m.Captures {
			switch capture.Index {
			case 0:
				// @attribute
				if capture.Node.Content(module) == "event" {
					s.IsEvent = true
				}
			case 1:
				// @name
				s.Name = capture.Node.Content(module)
			case 2:
				// @structBody
				pqFields := tree_sitter.NewQueryCursor()
				pqFields.Exec(queryFields, capture.Node)
				for {
					match, ok := pqFields.NextMatch()
					if !ok {
						break
					}
					f := Param{}
					for _, queryCapture := range match.Captures {
						switch queryCapture.Index {
						case 0:
							// @fieldName
							f.Name = queryCapture.Node.Content(module)
						case 1:
							// @type
							f.Type = queryCapture.Node.Content(module)
						}
					}
					s.Fields = append(s.Fields, f)
				}
			}
		}
		structs = append(structs, s)
	}

	return structs, nil
}

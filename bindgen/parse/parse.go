// nolint
package parse

import (
	"context"
	"fmt"

	tree_sitter_move_on_aptos "github.com/aptos-labs/tree-sitter-move-on-aptos/bindings/go"
	tree_sitter "github.com/smacker/go-tree-sitter"
)

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Func struct {
	IsView  bool `json:"is_view"`
	IsEntry bool `json:"is_entry"`

	Name        string   `json:"name"`
	Params      []Param  `json:"params"`
	ReturnTypes []string `json:"return_types"`
}

type Struct struct {
	IsEvent bool

	Name   string
	Fields []Param
}

type ExternalStruct struct {
	ImportPath string

	Package string
	Module  string
	Name    string
}

func ParseModule(module []byte) (pkg string, mod string, err error) {
	lang := tree_sitter.NewLanguage(tree_sitter_move_on_aptos.Language())
	n, err := tree_sitter.ParseCtx(context.Background(), module, lang)
	if err != nil {
		return "", "", fmt.Errorf("parsing AST: %w", err)
	}

	query, err := tree_sitter.NewQuery([]byte(`
(module
  path: (identifier) @package
  name: (identifier) @module
)
	`), lang)

	queryCursor := tree_sitter.NewQueryCursor()
	queryCursor.Exec(query, n)

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
				if capture.Node.Content(module) == "view" {
					f.IsView = true
				}
			case 2:
				// @modifier
				if capture.Node.Content(module) == "entry" {
					f.IsEntry = true
				}
			case 3:
				// @function_name
				f.Name = capture.Node.Content(module)
			case 4:
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
			case 5:
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

	// query to select all structs
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

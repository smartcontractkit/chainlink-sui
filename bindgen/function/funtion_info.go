package function

import (
	_ "embed"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const modulePrefix = "github.com/smartcontractkit/chainlink-sui/bindings/generated"

//go:embed go.tmpl
var tmpl string

// TODO: this asumes that module names will always be different,
// they can be the same if they are in different packages.
// Refactor to consider package path as well if needed.
type FunctionInfoData struct {
	ImportPath string
	ModuleName string
}

func GenerateGlobalFunctionInfo(baseDir string) error {
	packageNames, err := getPackagesWithFunctionInfo(baseDir)
	if err != nil {
		return err
	}
	t, err := template.New("function_info").Parse(tmpl)
	if err != nil {
		return err
	}

	// Create the output file in the baseDir
	outputPath := filepath.Join(baseDir, "function_info.go")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = t.Execute(outputFile, packageNames)
	if err != nil {
		return err
	}

	return nil
}

func getPackagesWithFunctionInfo(baseDir string) ([]FunctionInfoData, error) {
	baseDir = filepath.Clean(baseDir)

	fiData := []FunctionInfoData{}

	err := filepath.Walk(baseDir, func(
		path string,
		info os.FileInfo,
		err error,
	) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		packageName, err := hasFunctionInfo(path)
		if err != nil {
			return err
		}
		if packageName == nil {
			return nil
		}
		dir := filepath.Dir(path)
		packagePrefix := strings.TrimPrefix(dir, baseDir)
		fiData = append(fiData, FunctionInfoData{
			ImportPath: modulePrefix + packagePrefix,
			ModuleName: strings.TrimPrefix(*packageName, "module_"),
		})

		return nil
	})

	return fiData, err
}

func hasFunctionInfo(filePath string) (*string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	packageName := f.Name.Name

	hasFunctionInfo := false
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.GenDecl:
			if t.Tok == token.CONST {
				for _, spec := range t.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, name := range valueSpec.Names {
						if name.Name == "FunctionInfo" {
							// Found the FunctionInfo constant
							hasFunctionInfo = true
							return false
						}
					}
				}
			}
		}

		return true
	})

	if !hasFunctionInfo {
		return nil, nil
	}

	return &packageName, nil
}

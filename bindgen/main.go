// nolint
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/smartcontractkit/chainlink-sui/bindgen/parse"
	"github.com/smartcontractkit/chainlink-sui/bindgen/template"
)

func main() {
	inputFile := flag.String("input", "", "path to Move Sui contract file to parse")
	outputFolder := flag.String("output", "", "path to output directory")
	moveConfigPath := flag.String("moveConfig", "", "path to Move.toml file")
	uppercase := flag.String("uppercase", "", "list of words to convert to uppercase")
	externalStructs := flag.String("externalStructs", "", "comma-separated list of struct names, usage: --externalStructs ccip::ocr3_base::OCRConfig=github.com/smartcontractkit/chainlink-aptos/bindings/ccip/ocr3_base")

	flag.Parse()

	// Validate the move config path exists before using it
	if *moveConfigPath == "" {
		log.Fatalf("Move config path is required")
	}
	cleanPath := filepath.Clean(*moveConfigPath)
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		log.Fatalf("Move config file does not exist at path: %s", cleanPath)
	}

	log.Printf("Generating bindings for %s", *inputFile)

	if *uppercase != "" {
		for _, w := range strings.Split(*uppercase, ",") {
			template.UppercaseWords = append(template.UppercaseWords, strings.ToUpper(w))
		}
		log.Printf("Capitalizing %v words: %v", len(template.UppercaseWords), strings.Join(template.UppercaseWords, ", "))
	}

	// Parse external structs
	var extStructs []parse.ExternalStruct
	if *externalStructs != "" {
		for _, s := range strings.Split(*externalStructs, ",") {
			// package::module::Struct=github.com/smartcontractkit/chainlink-aptos/bindings/path
			split := strings.Split(s, "=")
			if len(split) != 2 {
				log.Fatalf("Invalid external structure definition: %v", s)
			}
			from := strings.Split(split[0], "::")
			if len(from) != 3 {
				log.Fatalf("Invalid external structure definition: %v", s)
			}
			packageName := from[0]
			moduleName := from[1]
			structName := from[2]
			importPath := split[1]

			log.Printf("Importing struct %v::%v::%v from %v", packageName, moduleName, structName, importPath)
			extStructs = append(extStructs, parse.ExternalStruct{
				ImportPath: importPath,
				Package:    packageName,
				Module:     moduleName,
				Name:       structName,
			})
		}
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	pkg, mod, err := parse.ParseModule(fileBytes)
	if err != nil {
		panic(err)
	}

	funcs, err := parse.ParseFunctions(fileBytes)
	if err != nil {
		panic(err)
	}
	log.Println("Parsed functions:")
	for i, viewFunc := range funcs {
		log.Println(i, viewFunc)
	}
	log.Println("----")
	structs, err := parse.ParseStructs(fileBytes)
	if err != nil {
		panic(err)
	}
	log.Println("Parsed structs:")
	for i, structt := range structs {
		log.Println(i, structt)
	}

	log.Println("----")
	data, err := template.Convert(pkg, mod, structs, funcs, extStructs)
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.Generate(data)
	if err != nil {
		log.Fatal(err)
	}

	outputFile := filepath.Join(*outputFolder, fmt.Sprintf("%s.go", data.Module))

	log.Printf("Writing output to %s", outputFile)
	_ = os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err := os.WriteFile(outputFile, []byte(t), 0600); err != nil {
		panic(err)
	}
}

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

	flag.Parse()

	fmt.Println(fmt.Sprintf(`
	##############################################################
	Generating Go bindings for: %s
	##############################################################
	`, *inputFile))

	// Validate the move config path exists before using it
	if *moveConfigPath == "" {
		log.Fatalf("Move config path is required")
	}
	cleanPath := filepath.Clean(*moveConfigPath)
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		log.Fatalf("Move config file does not exist at path: %s", cleanPath)
	}

	if *uppercase != "" {
		for _, w := range strings.Split(*uppercase, ",") {
			template.UppercaseWords = append(template.UppercaseWords, strings.ToUpper(w))
		}
		log.Printf("Capitalizing %v words: %v", len(template.UppercaseWords), strings.Join(template.UppercaseWords, ", "))
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

	structs, err := parse.ParseStructs(fileBytes)
	if err != nil {
		panic(err)
	}

	data, err := template.Convert(pkg, mod, structs, funcs)
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

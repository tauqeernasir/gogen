package main

import (
	"flag"
	"fmt"
	"gogen/internal/builder"
	"log"
	"strings"
)

func main() {
	var (
		specPath     = flag.String("spec", "", "Path to OpenAPI spec file")
		projectName  = flag.String("name", "", "Project name for the client")
		outputDir    = flag.String("output", "./generated-client", "Output directory")
		language     = flag.String("lang", "typescript", "Target language (typescript, python)")
		templatesDir = flag.String("templates", "", "Custom templates directory")
	)
	flag.Parse()

	if *specPath == "" || *projectName == "" {
		log.Fatal("Both -spec and -name flags are required")
	}

	generator := builder.NewClientGeneratorBuilder().
		WithSpec(*specPath).
		WithProjectName(*projectName).
		WithOutputDir(*outputDir).
		WithLanguage(*language).
		WithTemplatesDir(*templatesDir).
		Build()

	if err := generator.Generate(); err != nil {
		log.Fatal("Failed to generate client:", err)
	}

	fmt.Printf("%s client generated successfully in %s\n", strings.Title(*language), *outputDir)
}

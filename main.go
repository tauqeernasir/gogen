package main

import (
	"flag"
	"fmt"
	"gogen/internal/builder"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var (
		specPath     = flag.String("spec", "", "Path to OpenAPI spec file")
		projectName  = flag.String("name", "", "Project name for the client")
		outputDir    = flag.String("output", "./generated-client", "Output directory")
		language     = flag.String("lang", "typescript", "Target language (typescript, python)")
		templatesDir = flag.String("templates", "", "Custom templates directory")
		prettier     = flag.Bool("prettier", true, "Run prettier after generation")
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

	if *language == "typescript" && *prettier {
		absOutputDir, err := filepath.Abs(*outputDir)
		if err != nil {
			log.Printf("Warning: Could not get absolute path for outputDir '%s': %v", *outputDir, err)
			absOutputDir = *outputDir
		}

		if _, err := exec.LookPath("npx"); err == nil {
			cmd := exec.Command("npx", "prettier", "--write", filepath.Join(absOutputDir, "**", "*.ts"))
			cmd.Dir = absOutputDir
			if err := cmd.Run(); err != nil {
				log.Printf("Warning: Failed to run prettier: %v", err)
			}
		} else {
			log.Printf("Warning: prettier not found, skipping formatting")
		}
	}

	fmt.Printf("%s client generated successfully in %s\n", strings.Title(*language), *outputDir)
}

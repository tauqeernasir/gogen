package builder

import (
	"encoding/json"
	"fmt"
	"gogen/internal/adapters"
	"gogen/internal/models"
	"gogen/internal/openapi"
	"gogen/internal/templates"
	"gogen/internal/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ClientGeneratorBuilder struct {
	spec         *openapi.OpenAPISpec
	projectName  string
	outputDir    string
	language     string
	adapter      adapters.LanguageAdapter
	templateMgr  *templates.Manager
	templatesDir string
}

func NewClientGeneratorBuilder() *ClientGeneratorBuilder {
	return &ClientGeneratorBuilder{
		templateMgr: templates.NewManager(),
	}
}

func (b *ClientGeneratorBuilder) WithSpec(specPath string) *ClientGeneratorBuilder {
	data, err := ioutil.ReadFile(specPath)
	if err != nil {
		log.Fatal("Failed to read spec:", err)
	}

	b.spec = &openapi.OpenAPISpec{}
	if err := json.Unmarshal(data, b.spec); err != nil {
		log.Fatal("Failed to parse spec:", err)
	}

	return b
}

func (b *ClientGeneratorBuilder) WithProjectName(name string) *ClientGeneratorBuilder {
	b.projectName = name
	return b
}

func (b *ClientGeneratorBuilder) WithOutputDir(dir string) *ClientGeneratorBuilder {
	b.outputDir = dir
	return b
}

func (b *ClientGeneratorBuilder) WithLanguage(language string) *ClientGeneratorBuilder {
	b.language = language

	switch strings.ToLower(language) {
	case "typescript", "ts":
		b.adapter = adapters.NewTypeScriptAdapter()
	default:
		log.Fatal("Unsupported language:", language)
	}

	return b
}

func (b *ClientGeneratorBuilder) WithTemplatesDir(dir string) *ClientGeneratorBuilder {
	b.templatesDir = dir
	return b
}

func (b *ClientGeneratorBuilder) Build() *ClientGenerator {
	if b.spec == nil || b.projectName == "" || b.outputDir == "" || b.adapter == nil {
		log.Fatal("Missing required configuration")
	}

	if err := b.templateMgr.LoadTemplates(b.templatesDir); err != nil {
		log.Fatal("Failed to load templates:", err)
	}

	return &ClientGenerator{
		spec:        b.spec,
		projectName: b.projectName,
		outputDir:   b.outputDir,
		language:    b.language,
		adapter:     b.adapter,
		templateMgr: b.templateMgr,
	}
}

// Client generator
type ClientGenerator struct {
	spec        *openapi.OpenAPISpec
	projectName string
	outputDir   string
	language    string
	adapter     adapters.LanguageAdapter
	templateMgr *templates.Manager
}

func (g *ClientGenerator) Generate() error {
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return err
	}

	model := g.buildClientModel()

	files := g.getRequiredFiles()
	for _, file := range files {
		if err := g.generateFile(file, model); err != nil {
			return fmt.Errorf("failed to generate %s: %w", file, err)
		}
	}

	return nil
}

func (g *ClientGenerator) buildClientModel() *models.ClientModel {
	model := &models.ClientModel{
		ProjectName:  g.projectName,
		Description:  g.spec.Info.Description,
		Version:      g.spec.Info.Version,
		Dependencies: g.adapter.GetDependencies(),
		Methods:      g.buildMethods(),
		Types:        g.buildTypes(),
	}

	if len(g.spec.Servers) > 0 {
		model.BaseURL = g.spec.Servers[0].URL
	}

	return model
}

func (g *ClientGenerator) buildMethods() []models.MethodModel {
	var methods []models.MethodModel

	for path, pathItem := range g.spec.Paths {
		operations := map[string]*openapi.Operation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"DELETE": pathItem.Delete,
			"PATCH":  pathItem.Patch,
		}

		for httpMethod, operation := range operations {
			if operation == nil {
				continue
			}

			method := g.buildMethodModel(path, httpMethod, operation)
			methods = append(methods, method)
		}
	}

	return methods
}

func (g *ClientGenerator) buildMethodModel(path, httpMethod string, operation *openapi.Operation) models.MethodModel {
	var parameters []models.ParameterModel
	var requestBody *models.RequestBodyModel

	for _, param := range operation.Parameters {
		parameters = append(parameters, models.ParameterModel{
			Name:        g.adapter.FormatPropertyName(param.Name),
			Type:        g.adapter.ConvertType(param.Schema),
			In:          param.In,
			Required:    param.Required,
			Description: param.Description,
		})
	}

	if operation.RequestBody != nil {
		requestBody = &models.RequestBodyModel{
			Type:     g.getRequestBodyType(operation.RequestBody),
			Required: operation.RequestBody.Required,
		}
	}

	return models.MethodModel{
		//Name:         g.adapter.FormatMethodName(operation.OperationID, httpMethod, operation.Tags),
		Name:         operation.OperationID,
		HTTPMethod:   httpMethod,
		Path:         path,
		Summary:      operation.Summary,
		Description:  operation.Description,
		Parameters:   parameters,
		RequestBody:  requestBody,
		ResponseType: g.getResponseType(operation),
	}
}

func (g *ClientGenerator) buildTypes() []models.TypeModel {
	var types []models.TypeModel

	for name, schema := range g.spec.Components.Schemas {
		var typeModel *models.TypeModel

		switch {
		case schema.Type == "object":
			typeModel = g.processObjectSchema(name, &schema)
			break
		case schema.Type == "array":
			typeModel = g.processArraySchema(name, &schema)
			break
		}

		if typeModel != nil {
			types = append(types, *typeModel)
		}

		//if schema.Type == "array" && schema.Items != nil {
		//	fmt.Println(schema.Items.Type, g.adapter.FormatTypeName(name), schema.Items.Items)
		//	if schema.Items.Ref != "" {
		//		types = append(types, models.TypeModel{
		//			Name: g.adapter.FormatTypeName(name),
		//			Type: g.adapter.ConvertType(&schema),
		//		})
		//	}
		//
		//}
	}

	return types
}

func (g *ClientGenerator) processObjectSchema(name string, schema *openapi.Schema) *models.TypeModel {
	if schema.Properties == nil {
		return nil
	}

	var properties []models.PropertyModel

	for propName, propSchema := range schema.Properties {
		propType := g.adapter.ConvertType(propSchema)

		properties = append(properties, models.PropertyModel{
			Name:     g.adapter.FormatPropertyName(propName),
			Type:     propType,
			Required: utils.Contains(schema.Required, propName),
		})
	}

	return &models.TypeModel{
		Name:       g.adapter.FormatTypeName(name),
		Type:       "interface", // Keep as "interface" for objects
		Properties: properties,
	}
}

func (g *ClientGenerator) processArraySchema(name string, schema *openapi.Schema) *models.TypeModel {
	if schema.Items == nil {
		return &models.TypeModel{
			Name:       g.adapter.FormatTypeName(name),
			Type:       "any[]",
			Properties: nil,
		}
	}

	itemType := g.adapter.ConvertType(schema.Items)
	arrayType := itemType + "[]"

	return &models.TypeModel{
		Name:       g.adapter.FormatTypeName(name),
		Type:       arrayType,
		Properties: nil,
	}
}

func (g *ClientGenerator) getRequestBodyType(requestBody *openapi.RequestBody) string {
	for _, mediaType := range requestBody.Content {
		if mediaType.Schema != nil {
			return g.adapter.ConvertType(mediaType.Schema)
		}
	}
	return "any"
}

func (g *ClientGenerator) getResponseType(operation *openapi.Operation) string {
	for code, response := range operation.Responses {
		if strings.HasPrefix(code, "2") {
			for _, mediaType := range response.Content {
				if mediaType.Schema != nil {
					return g.adapter.ConvertType(mediaType.Schema)
				}
			}
		}
	}
	return "any"
}

func (g *ClientGenerator) getRequiredFiles() []string {
	switch g.language {
	case "typescript", "ts":
		return []string{"package.json", "tsconfig.json", "client", "types", "index", "README.md"}
	case "python", "py":
		return []string{"setup.py", "requirements.txt", "client", "types", "__init__", "README.md"}
	}
	return []string{}
}

func (g *ClientGenerator) generateFile(fileName string, model *models.ClientModel) error {
	templateName := fmt.Sprintf("%s/%s", g.language, fileName)
	tmpl, exists := g.templateMgr.GetTemplate(templateName)
	if !exists {
		return fmt.Errorf("template not found: %s", templateName)
	}

	var outputPath string
	if strings.Contains(fileName, ".") {
		outputPath = filepath.Join(g.outputDir, fileName)
	} else {
		outputPath = filepath.Join(g.outputDir, fileName+"."+g.adapter.GetFileExtension())
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	templateData := g.adapter.GetTemplateData(model)
	return tmpl.Execute(file, templateData)
}

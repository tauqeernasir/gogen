package adapters

import (
	"gogen/internal/models"
	"gogen/internal/openapi"
)

// LanguageAdapter defines the interface for language-specific code generation
type LanguageAdapter interface {
	// GetFileExtension returns the file extension for the target language
	GetFileExtension() string

	// GetDependencies returns the list of dependencies required for the generated client
	GetDependencies() []string

	// ConvertType converts an OpenAPI schema to a language-specific type
	ConvertType(schema *openapi.Schema) string

	// FormatMethodName formats a method name according to language conventions
	FormatMethodName(operationID, httpMethod string, tags []string) string

	// FormatTypeName formats a type name according to language conventions
	FormatTypeName(name string) string

	// FormatPropertyName formats a property name according to language conventions
	FormatPropertyName(name string) string

	// GetTemplateData prepares data for template rendering
	GetTemplateData(model *models.ClientModel) interface{}

	// FormatPath formats a path according to language conventions
	FormatPath(path, httpMethod string) string
}

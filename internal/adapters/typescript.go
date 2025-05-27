package adapters

import (
	"fmt"
	"gogen/internal/models"
	"gogen/internal/openapi"
	"gogen/internal/utils"
	"strings"
)

// TypeScriptAdapter implements LanguageAdapter for TypeScript
type TypeScriptAdapter struct{}

// NewTypeScriptAdapter creates a new TypeScript adapter
func NewTypeScriptAdapter() *TypeScriptAdapter {
	return &TypeScriptAdapter{}
}

// GetFileExtension returns the file extension for TypeScript files
func (ts *TypeScriptAdapter) GetFileExtension() string {
	return "ts"
}

// GetDependencies returns the list of dependencies for TypeScript clients
func (ts *TypeScriptAdapter) GetDependencies() []string {
	return []string{"axios"}
}

// ConvertType converts an OpenAPI schema to a TypeScript type
func (ts *TypeScriptAdapter) ConvertType(schema *openapi.Schema) string {
	if schema == nil {
		return "any"
	}

	if schema.Ref != "" {
		refName := strings.TrimPrefix(schema.Ref, "#/components/schemas/")
		return ts.FormatTypeName(refName)
	}

	switch schema.Type {
	case "string":
		if len(schema.Enum) > 0 {
			var enumValues []string
			for _, e := range schema.Enum {
				enumValues = append(enumValues, fmt.Sprintf("'%v'", e))
			}
			return strings.Join(enumValues, " | ")
		}
		return "string"
	case "integer", "number":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		if schema.Items != nil {
			return ts.ConvertType(schema.Items) + "[]"
		}
		return "any[]"
	case "object":
		if schema.Properties == nil {
			return "Record<string, any>"
		}
		// For complex objects, we'll reference the generated interface
		return "object"
	default:
		return "any"
	}
}

// FormatMethodName formats a method name using camelCase convention
func (ts *TypeScriptAdapter) FormatMethodName(operationID, httpMethod string, tags []string) string {
	if operationID != "" {
		return utils.ToCamelCase(operationID)
	}
	if len(tags) > 0 {
		return utils.ToCamelCase(tags[0] + httpMethod)
	}
	return utils.ToCamelCase(httpMethod + "Request")
}

// FormatTypeName formats a type name using PascalCase convention
func (ts *TypeScriptAdapter) FormatTypeName(name string) string {
	return utils.ToPascalCase(name)
}

// FormatPropertyName formats a property name (no change for TypeScript)
func (ts *TypeScriptAdapter) FormatPropertyName(name string) string {
	return name
}

// GetTemplateData prepares data for TypeScript template rendering
func (ts *TypeScriptAdapter) GetTemplateData(model *models.ClientModel) interface{} {
	return struct {
		*models.ClientModel
		ClientClassName string
	}{
		ClientModel:     model,
		ClientClassName: model.ProjectName + "Client",
	}
}

var builtInTypes = map[string]bool{
	"string":    true,
	"number":    true,
	"boolean":   true,
	"Date":      true,
	"RegExp":    true,
	"Error":     true,
	"Array":     true,
	"Map":       true,
	"Set":       true,
	"Promise":   true,
	"any":       true,
	"unknown":   true,
	"void":      true,
	"null":      true,
	"undefined": true,
	"object":    true,
}

func (ts *TypeScriptAdapter) IsBuiltInType(schema *openapi.Schema) bool {
	baseType := strings.Split(schema.Type, "<")[0]
	baseType = strings.TrimSpace(baseType)

	return builtInTypes[baseType]
}

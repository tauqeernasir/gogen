package adapters

import (
	"fmt"
	"gogen/internal/models"
	"gogen/internal/openapi"
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

	if len(schema.OneOf) > 0 {
		return ts.handleOneOf(schema)
	}

	if len(schema.AllOf) > 0 {
		return ts.handleAllOf(schema)
	}

	if len(schema.AnyOf) > 0 {
		return ts.handleAnyOf(schema)
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
		if schema.Items == nil {
			return "any[]"
		}

		if schema.Items.Ref != "" {
			refName := strings.TrimPrefix(schema.Items.Ref, "#/components/schemas/")
			return ts.FormatTypeName(refName) + "[]"
		}

		itemType := ts.ConvertType(schema.Items)
		return itemType + "[]"
	case "object":
		if schema.Properties == nil {
			return "Record<string, any>"
		}

		var properties []string
		for propName, propSchema := range schema.Properties {
			propType := ts.ConvertType(propSchema)
			properties = append(properties, fmt.Sprintf("%s: %s", propName, propType))
		}

		return "{" + strings.Join(properties, ", ") + "}"
	default:
		return "any"
	}
}

// FormatMethodName formats a method name using camelCase convention
func (ts *TypeScriptAdapter) FormatMethodName(operationID, httpMethod string, tags []string) string {
	if operationID != "" {
		return operationID
	}
	if len(tags) > 0 {
		return tags[0] + httpMethod
	}
	return httpMethod + "Request"
}

// FormatTypeName formats a type name using PascalCase convention
func (ts *TypeScriptAdapter) FormatTypeName(name string) string {
	return name
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

func (ts *TypeScriptAdapter) handleAllOf(schema *openapi.Schema) string {
	var types []string
	for _, subSchema := range schema.AllOf {
		types = append(types, ts.ConvertType(&subSchema))
	}

	return strings.Join(types, " & ")
}

func (ts *TypeScriptAdapter) handleOneOf(schema *openapi.Schema) string {
	var types []string
	for _, subSchema := range schema.OneOf {
		types = append(types, ts.ConvertType(&subSchema))
	}

	return strings.Join(types, " | ")
}

func (ts *TypeScriptAdapter) handleAnyOf(schema *openapi.Schema) string {
	var types []string
	for _, subSchema := range schema.AnyOf {
		types = append(types, ts.ConvertType(&subSchema))
	}

	return strings.Join(types, " | ")
}

func (ts *TypeScriptAdapter) FormatPath(path, httpMethod string) string {
	return strings.ReplaceAll(path, "{", "${")
}

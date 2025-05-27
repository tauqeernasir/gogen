package models

// ClientModel represents the complete client model for code generation
type ClientModel struct {
	ProjectName  string
	Description  string
	Version      string
	BaseURL      string
	Methods      []MethodModel
	Types        []TypeModel
	Dependencies []string
}

// MethodModel represents a single API method
type MethodModel struct {
	Name         string
	HTTPMethod   string
	Path         string
	Summary      string
	Description  string
	Parameters   []ParameterModel
	RequestBody  *RequestBodyModel
	ResponseType string
}

// ParameterModel represents a method parameter
type ParameterModel struct {
	Name        string
	Type        string
	In          string
	Required    bool
	Description string
}

// RequestBodyModel represents a request body
type RequestBodyModel struct {
	Type     string
	Required bool
}

// TypeModel represents a data type/schema
type TypeModel struct {
	Name       string
	Type       string
	Properties []PropertyModel
}

// PropertyModel represents a property of a type
type PropertyModel struct {
	Name     string
	Type     string
	Required bool
}

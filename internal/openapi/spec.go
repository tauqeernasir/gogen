package openapi

// OpenAPISpec represents the root OpenAPI specification
type OpenAPISpec struct {
	Info       Info                `json:"info"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
	Servers    []Server            `json:"servers"`
}

// Info contains metadata about the API
type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// Server represents a server configuration
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// PathItem describes the operations available on a single path
type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
}

// Operation describes a single API operation on a path
type Operation struct {
	OperationID string              `json:"operationId"`
	Summary     string              `json:"summary"`
	Description string              `json:"description"`
	Parameters  []Parameter         `json:"parameters"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses"`
	Tags        []string            `json:"tags"`
}

// Parameter describes a single operation parameter
type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Required    bool    `json:"required"`
	Description string  `json:"description"`
	Schema      *Schema `json:"schema"`
}

// RequestBody describes a single request body
type RequestBody struct {
	Content  map[string]MediaType `json:"content"`
	Required bool                 `json:"required"`
}

// Response describes a single response from an API Operation
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType provides schema and examples for the media type identified by its key
type MediaType struct {
	Schema *Schema `json:"schema"`
}

// Components holds a set of reusable objects for different aspects of the OAS
type Components struct {
	Schemas map[string]Schema `json:"schemas"`
}

// Schema allows the definition of input and output data types
type Schema struct {
	Type                 string             `json:"type"`
	Properties           map[string]*Schema `json:"properties"`
	Items                *Schema            `json:"items"`
	Required             []string           `json:"required"`
	Ref                  string             `json:"$ref"`
	AllOf                []Schema           `json:"allOf"`
	OneOf                []Schema           `json:"oneOf"`
	AnyOf                []Schema           `json:"anyOf"`
	Format               string             `json:"format"`
	Enum                 []interface{}      `json:"enum"`
	AdditionalProperties interface{}        `json:"additionalProperties"`
}

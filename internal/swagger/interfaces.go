package swagger

// OpenAPI is the root document of the OpenAPI document.
type OpenAPI struct {
	OpenAPIVersion string          `json:"openapi" yaml:"openapi"`
	Info           Info            `json:"info"`
	Servers        []Server        `json:"servers"`
	Paths          map[string]Path `json:"paths"`
	Components     Components      `json:"components"`
}

type Components struct {
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes" yaml:"securitySchemes"`
}

type SecurityScheme struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Scheme      string `json:"scheme"`
}

// Info provides metadata about the API.
type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// Server provides connectivity information to a target server. If the servers property is not provided,
// or is an empty array, the default value would be a Server Object with a url value of /.
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// Path contains the available operations for a given path in the API.
type Path struct {
	Summary     string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation  `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation  `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation  `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation  `json:"delete,omitempty" yaml:"delete,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

type Operation struct {
	Tags        []string              `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Security    []map[string][]string `json:"security"`
	Responses   map[int]Response      `json:"responses,omitempty" yaml:"responses,omitempty"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Required    bool                 `json:"required,omitempty" yaml:"required,omitempty"`
}

type MediaType struct {
	Schema Schema `json:"schema"`
}

type Schema struct {
	Type                 string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Format               string                 `json:"format,omitempty" yaml:"format,omitempty"`
	Nullable             bool                   `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	Items                *Schema                `json:"items,omitempty" yaml:"items,omitempty"`
	Properties           map[string]*Schema     `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *Schema                `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Example              map[string]interface{} `json:"example,omitempty" yaml:"example,omitempty"`
	Required             []string               `json:"required,omitempty" yaml:"required,omitempty"`
}

// Response describes a single response from an API Operation, including design-time, static links to operations based on the response.
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// SimpleParam is used to pass in to the Swagger-generating functions.
type SimpleParam struct {
	Type        string
	Description string
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
	Schema      Schema `json:"schema"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

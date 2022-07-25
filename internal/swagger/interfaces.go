package swagger

// OpenAPI is the root document of the OpenAPI document.
type OpenAPI struct {
	OpenAPIVersion string          `json:"openapi"`
	Info           Info            `json:"info"`
	Servers        []Server        `json:"servers"`
	Paths          map[string]Path `json:"paths"`
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
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Get         *Operation `json:"get,omitempty"`
	Put         *Operation `json:"put,omitempty"`
	Post        *Operation `json:"post,omitempty"`
	Delete      *Operation `json:"delete,omitempty"`
}

type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema Schema `json:"schema,omitempty"`
}

type Schema struct {
	Type       string                    `json:"type,omitempty"`
	Properties map[string]SchemaProperty `json:"properties,omitempty"`
	Example    map[string]interface{}    `json:"example,omitempty"`
}

type SchemaProperty struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
}

type Response struct {
	Description string `json:"description"`
}

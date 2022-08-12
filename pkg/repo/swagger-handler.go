package repo

import "github.com/kaphos/webapp/internal/swagger"

type SwaggerHandlerI interface {
	SetSummary(string)
	Summary() string
	SetDescription(string)
	Description() string
	AddParam(string, string, string)
	Params() map[string]swagger.SimpleParam
	AddResponse(int, string, interface{})
	AddResponses(...int)
	Responses() map[int]swagger.Response
}

// swaggerHandler is a helper struct that manages potential responses
// for a given handler. Can be used by both Repo and handlerBase.
type swaggerHandler struct {
	summary     string
	description string
	parameters  map[string]swagger.SimpleParam
	responses   map[int]swagger.Response
}

var _ SwaggerHandlerI = &swaggerHandler{}

func (f *swaggerHandler) Init() {
	if f.responses == nil {
		f.responses = map[int]swagger.Response{}
	}
}

func (f *swaggerHandler) SetSummary(summary string) { f.summary = summary }
func (f *swaggerHandler) Summary() string           { return f.summary }

func (f *swaggerHandler) SetDescription(description string) { f.description = description }
func (f *swaggerHandler) Description() string               { return f.description }

func (f *swaggerHandler) AddParam(name, varType, description string) {
	if f.parameters == nil {
		f.parameters = map[string]swagger.SimpleParam{}
	}

	f.parameters[name] = swagger.SimpleParam{Type: varType, Description: description}
}

func (f *swaggerHandler) Params() map[string]swagger.SimpleParam { return f.parameters }

// Responses returns the list of responses the Handler may return.
func (f *swaggerHandler) Responses() map[int]swagger.Response { return f.responses }

// AddResponse adds a single Swagger response into this Handler. Also supports
// tracking an expected response content, though this is not enforced or checked.
func (f *swaggerHandler) AddResponse(statusCode int, description string, payload interface{}) {
	resp := swagger.Response{Description: description}

	if payload != nil {
		// Associate a response with some content
		content, _ := swagger.GenContent(payload, false)
		resp.Content = *content
	}

	f.Init()
	f.responses[statusCode] = resp
}

var responseDescriptions = map[int]string{
	200: "OK",
	201: "Created",
	400: "Invalid request body", // automatically added for handlers with payloads
	401: "Unauthorised",
	500: "Internal server error", // automatically added for all handlers
}

// AddResponses is a helper function to bulk-add a series of "standard" responses.
// Given the status code, it will automatically include the description (and assume
// that there is no payload) as defined in responseDescriptions.
func (f *swaggerHandler) AddResponses(statusCodes ...int) {
	for _, code := range statusCodes {
		if description, ok := responseDescriptions[code]; ok {
			f.AddResponse(code, description, nil)
		}
	}
}

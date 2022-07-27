package repo

import "github.com/kaphos/webapp/internal/swagger"

type SwaggerHandlerI interface {
	AddResponse(int, string, interface{})
	AddResponses(...int)
	GetResponses() map[int]swagger.Response
}

// swaggerHandler is a helper struct that manages potential responses
// for a given handler. Can be used by both Repo and handlerBase.
type swaggerHandler struct {
	Responses map[int]swagger.Response
}

func (f *swaggerHandler) Init() {
	if f.Responses == nil {
		f.Responses = map[int]swagger.Response{}
	}
}

// GetResponses returns the list of responses the Handler may return.
func (f *swaggerHandler) GetResponses() map[int]swagger.Response { return f.Responses }

// AddResponse adds a single Swagger response into this Handler. Also supports
// tracking an expected response content, though this is not enforced or checked.
func (f *swaggerHandler) AddResponse(statusCode int, description string, payload interface{}) {
	resp := swagger.Response{Description: description}

	if payload != nil {
		// Associate a response with some content
		resp.Content = *swagger.GenContent(payload)
	}

	f.Init()
	f.Responses[statusCode] = resp
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

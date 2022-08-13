package swagger

type HandlerI interface {
	SetSummary(string)
	Summary() string
	SetDescription(string)
	Description() string
	AddParam(string, string, string)
	Params() map[string]SimpleParam
	AddResponse(int, string, interface{})
	AddResponses(...int)
	Responses() map[int]Response
}

// Handler is a helper struct that manages potential responses
// for a given handler. Can be used by both Repo and HandlerBase.
type Handler struct {
	summary     string
	description string
	parameters  map[string]SimpleParam
	responses   map[int]Response
}

var _ HandlerI = &Handler{}

func NewHandler() Handler {
	return Handler{
		responses: map[int]Response{},
	}
}

func (f *Handler) Init() {
	if f.responses == nil {
		f.responses = map[int]Response{}
	}
}

func (f *Handler) SetSummary(summary string) { f.summary = summary }
func (f *Handler) Summary() string           { return f.summary }

func (f *Handler) SetDescription(description string) { f.description = description }
func (f *Handler) Description() string               { return f.description }

func (f *Handler) AddParam(name, varType, description string) {
	if f.parameters == nil {
		f.parameters = map[string]SimpleParam{}
	}

	f.parameters[name] = SimpleParam{Type: varType, Description: description}
}

func (f *Handler) Params() map[string]SimpleParam { return f.parameters }

// Responses returns the list of responses the Handler may return.
func (f *Handler) Responses() map[int]Response { return f.responses }

// SetResponse is an internally-used function to set a Response for a given
// statusCode. Used because responses is set as private, to prevent accidental
// editing.
func (f *Handler) SetResponse(statusCode int, resp Response) {
	f.responses[statusCode] = resp
}

// AddResponse adds a single Swagger response into this Handler. Also supports
// tracking an expected response content, though this is not enforced or checked.
func (f *Handler) AddResponse(statusCode int, description string, payload interface{}) {
	resp := Response{Description: description}

	if payload != nil {
		// Associate a response with some content
		content, _ := GenContent(payload, false)
		resp.Content = *content
	}

	f.responses[statusCode] = resp
}

var ResponseDescriptions = map[int]string{
	200: "OK",
	201: "Created",
	400: "Invalid request body", // automatically added for handlers with payloads
	401: "Unauthorised",
	500: "Internal server error", // automatically added for all handlers
}

// AddResponses is a helper function to bulk-add a series of "standard" responses.
// Given the status code, it will automatically include the description (and assume
// that there is no payload) as defined in ResponseDescriptions.
func (f *Handler) AddResponses(statusCodes ...int) {
	for _, code := range statusCodes {
		if description, ok := ResponseDescriptions[code]; ok {
			f.AddResponse(code, description, nil)
		}
	}
}

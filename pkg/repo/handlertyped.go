package repo

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/swagger"
	"go/types"
	"net/http"
)

// HandlerU represents an untyped/unvalidated HTTP Handler.
// Should create a new instance using NewHandlerU instead of
// instantiating this struct.
type HandlerU struct {
	handlerBase[types.Nil]
	handler handlerFuncU
}

// HandlerP represents a typed and validated HTTP Handler
// that expects a payload (hence, P). Should create a new
// instance using NewHandlerP instead of instantiating this struct.
type HandlerP[T any] struct {
	handlerBase[T]
	handler handlerFuncP[T]
}

var _ HandlerBaseI = &HandlerU{}
var _ HandlerBaseI = &HandlerP[types.Nil]{}

// Handle is an implementation of gin.HandleFunc, and provides automated
// handling of status codes, depending on whether f.handler was successful
// or not. Used by Server internally to attach a Repo to it.
func (f *HandlerU) Handle(c *gin.Context) {
	if ok := f.handler(c); ok {
		c.Status(f.successCode)
	} else if c.Writer.Status() < 300 {
		c.Status(http.StatusTeapot) // catch-all, in case we didn't set
	}
}

// Handle is an implementation of gin.HandleFunc, and provides automated
// handling of status codes, depending on whether f.handler was successful
// or not. Used by Server internally to attach a Repo to it.
func (f *HandlerP[T]) Handle(c *gin.Context) {
	var obj T

	if err := c.ShouldBind(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok := f.handler(c, obj); ok {
		c.Status(f.successCode)
	} else if c.Writer.Status() < 300 {
		c.Status(http.StatusTeapot) // catch-all, in case we didn't set
	}
}

// NewHandlerU creates a new unvalidated handler (i.e., does not expect or parse any
// payload). A method, relativePath and fn must be passed. fn will be called when the
// route matches. Middleware can also optionally be added.
func NewHandlerU(method, relativePath string, fn handlerFuncU, successCode int, successContent interface{}, middleware ...Middleware) HandlerU {
	handler := HandlerU{
		handler: fn,
		handlerBase: handlerBase[types.Nil]{
			method:      method,
			successCode: successCode,
			httpBase: httpBase{
				relativePath: relativePath,
				swaggerHandler: swaggerHandler{
					responses: map[int]swagger.Response{},
				},
			},
		},
	}

	handler.AddResponse(successCode, "Success", successContent)
	handler.AddResponses(500)
	handler.SetMiddleware(middleware...)

	return handler
}

// NewHandlerP creates a new validated handler with an expected payload.
// A method, relativePath and fn must be passed. fn will be called when the
// route matches, and the parsed payload will be passed in.
// Middleware can also optionally be added.
func NewHandlerP[T any](method, relativePath string, fn handlerFuncP[T], successCode int, successContent interface{}, middleware ...Middleware) HandlerP[T] {
	handler := HandlerP[T]{
		handler: fn,
		handlerBase: handlerBase[T]{
			method:      method,
			successCode: successCode,
			httpBase: httpBase{
				relativePath: relativePath,
				swaggerHandler: swaggerHandler{
					responses: map[int]swagger.Response{},
				},
			},
		},
	}

	handler.AddResponse(successCode, "Success", successContent)
	handler.AddResponses(400, 500)
	handler.SetMiddleware(middleware...)
	return handler
}

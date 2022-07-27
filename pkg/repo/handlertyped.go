package repo

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/swagger"
	"go/types"
	"net/http"
)

type HandlerU struct {
	handlerBase[types.Nil]
	Func HandlerFuncU
}

type HandlerP[T any] struct {
	handlerBase[T]
	Func HandlerFuncP[T]
}

var _ HandlerBaseI = &HandlerU{}
var _ HandlerBaseI = &HandlerP[types.Nil]{}

func (f *HandlerU) Handle(c *gin.Context) {
	if ok := f.Func(c); ok {
		c.Status(f.SuccessStatusCode)
	} else if c.Writer.Status() < 300 {
		c.Status(http.StatusTeapot) // catch-all, in case we didn't set
	}
}

func (f *HandlerP[T]) Handle(c *gin.Context) {
	var obj T

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok := f.Func(c, obj); ok {
		c.Status(f.SuccessStatusCode)
	} else if c.Writer.Status() < 300 {
		c.Status(http.StatusTeapot) // catch-all, in case we didn't set
	}
}

// NewHandlerU creates a new unvalidated handler (i.e., does not expect or parse any
// payload). A method, relativePath and fn must be passed. fn will be called when the
// route matches. Middleware can also optionally be added.
func NewHandlerU(method, relativePath string, fn HandlerFuncU, successCode int, successContent interface{}, middleware ...Middleware) HandlerU {
	handler := HandlerU{
		Func: fn,
		handlerBase: handlerBase[types.Nil]{
			Method:            method,
			SuccessStatusCode: successCode,
			httpBase: httpBase{
				RelativePath: relativePath,
				swaggerHandler: swaggerHandler{
					Responses: map[int]swagger.Response{},
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
func NewHandlerP[T any](method, relativePath string, fn HandlerFuncP[T], successCode int, successContent interface{}, middleware ...Middleware) HandlerP[T] {
	handler := HandlerP[T]{
		Func: fn,
		handlerBase: handlerBase[T]{
			Method:            method,
			SuccessStatusCode: successCode,
			httpBase: httpBase{
				RelativePath: relativePath,
				swaggerHandler: swaggerHandler{
					Responses: map[int]swagger.Response{},
				},
			},
		},
	}

	handler.AddResponse(successCode, "Success", successContent)
	handler.AddResponses(400, 500)
	handler.SetMiddleware(middleware...)
	return handler
}

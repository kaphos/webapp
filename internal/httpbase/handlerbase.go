package httpbase

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/swagger"
)

// HandlerBaseI extends HTTPBaseI (which extends SwaggerHandlerI), and adds upon it the ability
// to handle methods and handlers. Extended by U and P, which implement the
// "Handle" function defined here.
type HandlerBaseI interface {
	I
	Method() string
	SuccessCode() int
	Type() interface{}
	Handle(*gin.Context)
}

type HandlerBase[T any] struct {
	HTTPBase
	method      string
	successCode int
}

func (f *HandlerBase[T]) Method() string    { return f.method }
func (f *HandlerBase[T]) SuccessCode() int  { return f.successCode }
func (f *HandlerBase[T]) Type() interface{} { return *new(T) }

func NewHandlerBase[T any](method string, successCode int, relativePath string) HandlerBase[T] {
	return HandlerBase[T]{
		method:      method,
		successCode: successCode,
		HTTPBase: HTTPBase{
			relativePath: relativePath,
			Handler:      swagger.NewHandler(),
		},
	}
}

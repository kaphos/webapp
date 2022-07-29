package repo

import (
	"github.com/gin-gonic/gin"
)

// HandlerBaseI extends HTTPBaseI (which extends SwaggerHandlerI), and adds upon it the ability
// to handle methods and handlers. Extended by HandlerU and HandlerP, which implement the
// "Handle" function defined here.
type HandlerBaseI interface {
	HTTPBaseI
	Method() string
	Type() interface{}
	Handle(*gin.Context)
}

type handlerBase[T any] struct {
	httpBase
	method      string
	successCode int
}

func (f *handlerBase[T]) Method() string    { return f.method }
func (f *handlerBase[T]) Type() interface{} { return *new(T) }

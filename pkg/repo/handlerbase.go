package repo

import (
	"github.com/gin-gonic/gin"
)

type HandlerBaseI interface {
	HTTPBaseI
	GetMethod() string
	GetType() interface{}
	Handle(*gin.Context)
}

// handlerBase extends httpBase (which extends swaggerHandler), and adds upon it the ability
// to handle methods and handlers. Extended by HandlerU and HandlerP, which implement the
// "Handle" function.
type handlerBase[T any] struct {
	httpBase
	Method            string
	SuccessStatusCode int
}

func (f *handlerBase[T]) GetMethod() string    { return f.Method }
func (f *handlerBase[T]) GetType() interface{} { return *new(T) }

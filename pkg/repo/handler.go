package repo

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/swagger"
	"go/types"
	"net/http"
)

type HandlerInterface interface {
	GetMethod() string
	GetRelativePath() string
	GetType() interface{}
	handle(c *gin.Context)
	GetHandlers() []gin.HandlerFunc
	GetResponses() map[string]swagger.Response
}

type HandlerFuncNoPayload[T any] struct {
	Method       string
	RelativePath string
	Func         func(*gin.Context)
	Middleware   []gin.HandlerFunc
	Responses    map[string]swagger.Response
}

type HandlerFuncWithPayload[T any] struct {
	Method       string
	RelativePath string
	Func         func(*gin.Context, T)
	Middleware   []gin.HandlerFunc
	Responses    map[string]swagger.Response
}

var _ HandlerInterface = &HandlerFuncNoPayload[types.Nil]{}
var _ HandlerInterface = &HandlerFuncWithPayload[types.Nil]{}

func NewHandlerU[T any](method, relativePath string, fn func(*gin.Context), middleware ...gin.HandlerFunc) HandlerFuncNoPayload[T] {
	return HandlerFuncNoPayload[T]{
		Method:       method,
		RelativePath: relativePath,
		Func:         fn,
		Middleware:   middleware,
		Responses:    map[string]swagger.Response{},
	}
}

func NewHandlerP[T any](method, relativePath string, fn func(*gin.Context, T), middleware ...gin.HandlerFunc) HandlerFuncWithPayload[T] {
	return HandlerFuncWithPayload[T]{
		Method:       method,
		RelativePath: relativePath,
		Func:         fn,
		Middleware:   middleware,
		Responses:    map[string]swagger.Response{},
	}
}

func (f *HandlerFuncNoPayload[T]) GetMethod() string                { return f.Method }
func (f *HandlerFuncNoPayload[T]) GetMiddleware() []gin.HandlerFunc { return f.Middleware }
func (f *HandlerFuncNoPayload[T]) GetRelativePath() string          { return f.RelativePath }
func (f *HandlerFuncNoPayload[T]) GetType() interface{}             { return nil }
func (f *HandlerFuncNoPayload[T]) AddResponse(statusCode, description string) {
	f.Responses[statusCode] = swagger.Response{Description: description}
}
func (f *HandlerFuncNoPayload[T]) GetResponses() map[string]swagger.Response { return f.Responses }
func (f *HandlerFuncNoPayload[T]) handle(c *gin.Context) {
	f.Func(c)
}
func (f *HandlerFuncNoPayload[T]) GetHandlers() []gin.HandlerFunc {
	handles := make([]gin.HandlerFunc, 0)
	handles = append(handles, f.GetMiddleware()...)
	handles = append(handles, f.handle)
	return handles
}

func (f *HandlerFuncWithPayload[T]) GetMethod() string                { return f.Method }
func (f *HandlerFuncWithPayload[T]) GetMiddleware() []gin.HandlerFunc { return f.Middleware }
func (f *HandlerFuncWithPayload[T]) GetRelativePath() string          { return f.RelativePath }
func (f *HandlerFuncWithPayload[T]) GetType() interface{}             { return *new(T) }
func (f *HandlerFuncWithPayload[T]) AddResponse(statusCode, description string) {
	f.Responses[statusCode] = swagger.Response{Description: description}
}
func (f *HandlerFuncWithPayload[T]) GetResponses() map[string]swagger.Response { return f.Responses }
func (f *HandlerFuncWithPayload[T]) handle(c *gin.Context) {
	var obj T

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errchk": err.Error()})
		return
	}

	f.Func(c, obj)
}
func (f *HandlerFuncWithPayload[T]) GetHandlers() []gin.HandlerFunc {
	handles := make([]gin.HandlerFunc, 0)
	handles = append(handles, f.GetMiddleware()...)
	handles = append(handles, f.handle)
	return handles
}

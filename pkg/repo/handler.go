package repo

import (
	"github.com/gin-gonic/gin"
	"go/types"
	"net/http"
)

type HandlerInterface interface {
	GetMethod() string
	GetRelativePath() string
	handle(c *gin.Context)
	GetHandlers() []gin.HandlerFunc
}

type HandlerFuncNoPayload[T any] struct {
	Method       string
	RelativePath string
	Func         func(*gin.Context)
	Middleware   []gin.HandlerFunc
}

type HandlerFuncWithPayload[T any] struct {
	Method       string
	RelativePath string
	Func         func(*gin.Context, T)
	Middleware   []gin.HandlerFunc
}

var _ HandlerInterface = &HandlerFuncNoPayload[types.Nil]{}
var _ HandlerInterface = &HandlerFuncWithPayload[types.Nil]{}

func (f *HandlerFuncNoPayload[T]) GetMethod() string                { return f.Method }
func (f *HandlerFuncNoPayload[T]) GetMiddleware() []gin.HandlerFunc { return f.Middleware }
func (f *HandlerFuncNoPayload[T]) GetRelativePath() string          { return f.RelativePath }

func (f *HandlerFuncWithPayload[T]) GetMethod() string                { return f.Method }
func (f *HandlerFuncWithPayload[T]) GetMiddleware() []gin.HandlerFunc { return f.Middleware }
func (f *HandlerFuncWithPayload[T]) GetRelativePath() string          { return f.RelativePath }

// Handle an incoming HTTP request.
func (f *HandlerFuncNoPayload[T]) handle(c *gin.Context) {
	f.Func(c)
}

func (f *HandlerFuncNoPayload[T]) GetHandlers() []gin.HandlerFunc {
	handles := make([]gin.HandlerFunc, 0)
	handles = append(handles, f.GetMiddleware()...)
	handles = append(handles, f.handle)
	return handles
}

// Handle an incoming HTTP request.
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

package repo

import (
	"github.com/gin-gonic/gin"
	"go/types"
	"net/http"
)

type HandlerInterface interface {
	GetMethod() string
	GetRelativePath() string
	Handle(c *gin.Context)
}

type HandlerFuncNoPayload[T any] struct {
	Method       string
	RelativePath string
	Func         func(*gin.Context)
}

type HandlerFuncWithPayload[T any] struct {
	Method       string
	RelativePath string
	Func         func(*gin.Context, T)
}

var _ HandlerInterface = &HandlerFuncNoPayload[types.Nil]{}
var _ HandlerInterface = &HandlerFuncWithPayload[types.Nil]{}

func (f *HandlerFuncNoPayload[T]) GetMethod() string         { return f.Method }
func (f *HandlerFuncWithPayload[T]) GetMethod() string       { return f.Method }
func (f *HandlerFuncNoPayload[T]) GetRelativePath() string   { return f.RelativePath }
func (f *HandlerFuncWithPayload[T]) GetRelativePath() string { return f.RelativePath }

// Handle an incoming HTTP request.
func (f *HandlerFuncNoPayload[T]) Handle(c *gin.Context) { f.Func(c) }

// Handle an incoming HTTP request.
func (f *HandlerFuncWithPayload[T]) Handle(c *gin.Context) {
	var obj T

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errchk": err.Error()})
		return
	}

	f.Func(c, obj)
}

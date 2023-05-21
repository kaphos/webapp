package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/httpbase"
	"go/types"
	"net/http"
)

// FuncU is an extension of gin.HandlerFunc, but expects
// a bool response on whether the function was successful or not.
// Used to, for example, return a 401 status code if an auth middleware
// fails, or return 200 status code if a handlers passes.
type FuncU func(*gin.Context) bool

// FuncP is an extension of gin.HandlerFunc, but expects
// a bool response on whether the function was successful or not.
// Used to, for example, return a 401 status code if an auth middleware
// fails, or return 200 status code if a handlers passes. Differs from
// FuncU as it takes in a type T as well.
type FuncP[T any] func(*gin.Context, T) bool

// U represents an untyped/unvalidated HTTP Handler.
// Should create a new instance using NewU instead of
// instantiating this struct.
type U struct {
	httpbase.HandlerBase[types.Nil]
	handler FuncU
}

// P represents a typed and validated HTTP Handler
// that expects a payload (hence, P). Should create a new
// instance using NewP instead of instantiating this struct.
type P[T any] struct {
	httpbase.HandlerBase[T]
	handler FuncP[T]
}

var _ httpbase.HandlerBaseI = &U{}
var _ httpbase.HandlerBaseI = &P[types.Nil]{}

// Handle is an implementation of gin.HandleFunc, and provides automated
// handling of status codes, depending on whether f.handlers was successful
// or not. Used by Server internally to attach a Repo to it.
func (f *U) Handle(c *gin.Context) {
	if ok := f.handler(c); ok {
		c.Status(f.SuccessCode())
	} else if c.Writer.Status() < 300 {
		c.Status(http.StatusTeapot) // catch-all; returned false but no status code was set in the function
	}
}

// Handle is an implementation of gin.HandleFunc, and provides automated
// handling of status codes, depending on whether f.handlers was successful
// or not. Used by Server internally to attach a Repo to it.
func (f *P[T]) Handle(c *gin.Context) {
	var obj T

	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok := f.handler(c, obj); ok {
		c.Status(f.SuccessCode())
	} else if c.Writer.Status() < 300 {
		c.Status(http.StatusTeapot) // catch-all; returned false but no status code was set in the function
	}
}

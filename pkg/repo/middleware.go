package repo

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/swagger"
)

// HandlerFuncU is an extension of gin.HandlerFunc, but expects
// a bool response on whether the function was successful or not.
// Used to, for example, return a 401 status code if an auth middleware
// fails, or return 200 status code if a handler passes.
type HandlerFuncU func(ctx *gin.Context) bool

// HandlerFuncP is an extension of gin.HandlerFunc, but expects
// a bool response on whether the function was successful or not.
// Used to, for example, return a 401 status code if an auth middleware
// fails, or return 200 status code if a handler passes. Differs from
// HandlerFuncU as it takes in a type T as well.
type HandlerFuncP[T any] func(*gin.Context, T) bool

// Middleware is a wrapper around the Gin HandlerFunc object.
type Middleware struct {
	// Function to call to run middleware. The Repo struct will automatically Handle
	// failures by calling Abort, and continue the function call by calling Next.
	Fn             HandlerFuncU
	FailStatusCode int              // Status code to return if middleware fails
	FailResponse   swagger.Response // Swagger response if middleware fails
}

func NewMiddleware(fn HandlerFuncU, failCode int, failDescription string) Middleware {
	return Middleware{
		Fn:             fn,
		FailStatusCode: failCode,
		FailResponse:   swagger.Response{Description: failDescription},
	}
}

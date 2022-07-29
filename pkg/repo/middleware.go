package repo

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/swagger"
)

// handlerFuncU is an extension of gin.HandlerFunc, but expects
// a bool response on whether the function was successful or not.
// Used to, for example, return a 401 status code if an auth middleware
// fails, or return 200 status code if a handler passes.
type handlerFuncU func(ctx *gin.Context) bool

// handlerFuncP is an extension of gin.HandlerFunc, but expects
// a bool response on whether the function was successful or not.
// Used to, for example, return a 401 status code if an auth middleware
// fails, or return 200 status code if a handler passes. Differs from
// handlerFuncU as it takes in a type T as well.
type handlerFuncP[T any] func(*gin.Context, T) bool

// Middleware is a wrapper around the Gin HandlerFunc object.
type Middleware struct {
	// Function to call to run middleware. The Repo struct will automatically Handle
	// failures by calling Abort, and continue the function call by calling Next.
	Fn             handlerFuncU
	FailStatusCode int              // Status code to return if middleware fails
	FailResponse   swagger.Response // Swagger response if middleware fails
}

// NewMiddleware creates a new Middleware object, taking in a function that should
// return a boolean value (true if middleware passes), and a fail code to return
// if the middleware does not pass.
func NewMiddleware(fn handlerFuncU, failCode int, failDescription string) Middleware {
	return Middleware{
		Fn:             fn,
		FailStatusCode: failCode,
		FailResponse:   swagger.Response{Description: failDescription},
	}
}

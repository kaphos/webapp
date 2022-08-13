package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/swagger"
)

// Middleware is a wrapper around the Gin HandlerFunc object.
type Middleware struct {
	// Function to call to run middleware. The Repo struct will automatically Handle
	// failures by calling Abort, and continue the function call by calling Next.
	Fn             func(ctx *gin.Context) bool
	FailStatusCode int              // Status code to return if middleware fails
	FailResponse   swagger.Response // Swagger response if middleware fails
}

// New creates a new Middleware object, taking in a function that should
// return a boolean value (true if middleware passes), and a fail code to return
// if the middleware does not pass.
func New(fn func(ctx *gin.Context) bool, failCode int, failDescription string) Middleware {
	return Middleware{
		Fn:             fn,
		FailStatusCode: failCode,
		FailResponse:   swagger.Response{Description: failDescription},
	}
}

// NewAuth is a shortened function, to automatically set fail status code
// and fail response (401, "Unauthorised") instead of needing to key it in manually
// every time.
func NewAuth(fn func(ctx *gin.Context) bool) Middleware {
	return New(fn, 401, "Unauthorised")
}

package repo

import (
	"github.com/gin-gonic/gin"
)

type HTTPBaseI interface {
	SwaggerHandlerI
	GetRelativePath() string
	GetMiddleware() *[]gin.HandlerFunc
	SetMiddleware(middleware ...Middleware)
}

// httpBase extends swaggerHandler, providing support for tracking relative path and middleware.
// This allows it to be used by both handlerBase and Repo.
type httpBase struct {
	swaggerHandler
	RelativePath string
	Middleware   []gin.HandlerFunc
}

var _ HTTPBaseI = &httpBase{}

func (f *httpBase) GetRelativePath() string { return f.RelativePath }

func (f *httpBase) GetMiddleware() *[]gin.HandlerFunc { return &f.Middleware }

// SetMiddleware takes in a list of Middleware, and both adds it to the chain of middleware
// (which is used when returning handlers), and adds the possible failure status to the
// Swagger documentation.
func (f *httpBase) SetMiddleware(middleware ...Middleware) {
	f.Middleware = make([]gin.HandlerFunc, 0)

	for _, m := range middleware {
		// Process and add the middleware
		f.Middleware = append(f.Middleware, func(c *gin.Context) {
			if ok := m.Fn(c); !ok {
				c.AbortWithStatus(m.FailStatusCode)
			} else {
				c.Next()
			}
		})

		// Add the response value to be tracked by Swagger
		f.Init()
		f.Responses[m.FailStatusCode] = m.FailResponse
	}
}

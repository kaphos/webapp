package repo

import (
	"github.com/gin-gonic/gin"
)

type HTTPBaseI interface {
	SwaggerHandlerI
	RelativePath() string
	Middleware() *[]gin.HandlerFunc
	SetMiddleware(middleware ...Middleware)
}

// httpBase extends swaggerHandler, providing support for tracking relative path and middleware.
// This allows it to be used by both handlerBase and Repo.
type httpBase struct {
	swaggerHandler
	relativePath string
	middleware   []gin.HandlerFunc
}

var _ HTTPBaseI = &httpBase{}

func (f *httpBase) RelativePath() string        { return f.relativePath }
func (f *httpBase) SetRelativePath(path string) { f.relativePath = path }

func (f *httpBase) Middleware() *[]gin.HandlerFunc { return &f.middleware }

// SetMiddleware takes in a list of Middleware, and both adds it to the chain of middleware
// (which is used when returning handlers), and adds the possible failure status to the
// Swagger documentation.
func (f *httpBase) SetMiddleware(middleware ...Middleware) {
	f.middleware = make([]gin.HandlerFunc, 0)

	for _, m := range middleware {
		// Process and add the middleware
		f.middleware = append(f.middleware, func(c *gin.Context) {
			if ok := m.Fn(c); !ok {
				c.AbortWithStatus(m.FailStatusCode)
			} else {
				c.Next()
			}
		})

		// Add the response value to be tracked by Swagger
		f.Init()
		f.responses[m.FailStatusCode] = m.FailResponse
	}
}

// Package httpbase defines a struct that handles some basic routing.
// This struct, HTTPBase, can then be used in both repo.Repo and handler.Base
// to perform basic routing needs.
package httpbase

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/swagger"
	"github.com/kaphos/webapp/pkg/middleware"
)

type I interface {
	swagger.HandlerI
	RelativePath() string
	SetRelativePath(string)
	Middleware() *[]gin.HandlerFunc
	SetMiddleware(middleware ...middleware.Middleware)
}

// HTTPBase extends swagger.Handler, providing support for tracking relative path and middlewares.
// It is used by both HandlerBase and Repo.
type HTTPBase struct {
	swagger.Handler
	relativePath string
	middleware   []gin.HandlerFunc
}

var _ I = &HTTPBase{}

func (f *HTTPBase) RelativePath() string        { return f.relativePath }
func (f *HTTPBase) SetRelativePath(path string) { f.relativePath = path }

func (f *HTTPBase) Middleware() *[]gin.HandlerFunc { return &f.middleware }

// SetMiddleware takes in a list of Middleware, and both adds it to the chain of middleware
// (which is used when returning handlers), and adds the possible failure status to the
// Swagger documentation.
func (f *HTTPBase) SetMiddleware(middleware ...middleware.Middleware) {
	f.Init()
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

		f.AddResponse(m.FailStatusCode, m.FailResponse.Description, m.FailResponse.Content)
	}
}

package handler

import (
	"github.com/kaphos/webapp/internal/httpbase"
	"github.com/kaphos/webapp/pkg/middleware"
	"go/types"
)

// NewU creates a new unvalidated handlers (i.e., does not expect or parse any
// payload). A method, relativePath and fn must be passed. fn will be called when the
// route matches. Middleware can also optionally be added.
func NewU(method, relativePath string, fn FuncU, successCode int, successContent interface{}, middleware ...middleware.Middleware) U {
	h := U{
		handler:     fn,
		HandlerBase: httpbase.NewHandlerBase[types.Nil](method, successCode, relativePath),
	}

	h.AddResponse(successCode, "Success", successContent)
	h.AddResponses(500)
	h.SetMiddleware(middleware...)

	return h
}

// NewP creates a new validated handlers with an expected payload.
// A method, relativePath and fn must be passed. fn will be called when the
// route matches, and the parsed payload will be passed in.
// Middleware can also optionally be added.
func NewP[T any](method, relativePath string, fn FuncP[T], successCode int, successContent interface{}, middleware ...middleware.Middleware) P[T] {
	h := P[T]{
		handler:     fn,
		HandlerBase: httpbase.NewHandlerBase[T](method, successCode, relativePath),
	}

	h.AddResponse(successCode, "Success", successContent)
	h.AddResponses(400, 500)
	h.SetMiddleware(middleware...)
	return h
}

package repo

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/db"
	"go/types"
)

// Interface defines the expected functions that a Server expects any
// attached repositories to have.
type Interface interface {
	Init(database *db.Database)        // initialises any connections/configurations
	GetRelativePath() string           // returns the relative path associated with the repo
	GetMiddleware() *[]gin.HandlerFunc // returns the list of middleware that should be included for the repo
	GetHandlers() *[]HandlerInterface  // retrieve handlers, for attaching to the server and documentation
}

// Repo represents a collection of APIs around one entity.
// Should implement Interface.
type Repo[T any] struct {
	RelativePath string             // path that the repo should be accessible at
	Middleware   []gin.HandlerFunc  // any middleware that should be included
	Handlers     []HandlerInterface // list of handlers
	DB           *db.Database       // database object; initialised by the server
}

var _ Interface = &Repo[types.Nil]{}

// Init is called internally by the server when the Repo is attached to the server,
// to set up the database and tracer instance.
func (r *Repo[T]) Init(database *db.Database) {
	r.DB = database
}

// GetRelativePath returns the path that the repo is configured to listen to.
// Intended to be used by the Server object when attaching a Repo.
func (r *Repo[T]) GetRelativePath() string {
	return r.RelativePath
}

func (r *Repo[T]) AddMiddleware(handlers ...gin.HandlerFunc) {
	r.Middleware = append(r.Middleware, handlers...)
}

func (r *Repo[T]) GetMiddleware() *[]gin.HandlerFunc {
	return &r.Middleware
}

// GetHandlers returns the list of handlers added to the repo.
// Intended to be used by the Server object when attaching a Repo.
func (r *Repo[T]) GetHandlers() *[]HandlerInterface {
	return &r.Handlers
}

func (r *Repo[T]) AddHandler(h HandlerInterface) {
	r.Handlers = append(r.Handlers, h)
}

//// AddHandlerU attaches an unvalidated handler that does not require any validation/parsing of a payload
//// (or any parsing will be done within the handler function).
//func (r *Repo[T]) AddHandlerU(method, relativePath string, responses map[string]swagger.Response, fn func(*gin.Context), middleware ...gin.HandlerFunc) {
//	handler := HandlerFuncNoPayload[T]{
//		Method:       method,
//		RelativePath: relativePath,
//		Func:         fn,
//		Middleware:   middleware,
//		Responses:    responses,
//	}
//	r.Handlers = append(r.Handlers, &handler)
//}
//
//// AddHandlerP attaches a handler that requires parsing of a payload. Will attempt to
//// bind the request body to a new instance of T (generic type which is set when instantiating
//// the repo), and automatically rejects the request if there are any bind errors.
//// Uses Gin's validator, referenced at https://github.com/gin-gonic/gin#model-binding-and-validation.
//func (r *Repo[T]) AddHandlerP(method, relativePath string, responses map[string]swagger.Response, fn func(*gin.Context, T), middleware ...gin.HandlerFunc) {
//	handler := HandlerFuncWithPayload[T]{
//		Method:       method,
//		RelativePath: relativePath,
//		Func:         fn,
//		Middleware:   middleware,
//		Responses:    responses,
//	}
//	r.Handlers = append(r.Handlers, &handler)
//}

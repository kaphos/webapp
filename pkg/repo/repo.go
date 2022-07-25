package repo

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/db"
	"github.com/kaphos/webapp/internal/log"
	"go.opentelemetry.io/otel/trace"
	"go/types"
	"strconv"
)

// Interface defines the expected functions that a Server expects any
// attached repositories to have.
type Interface interface {
	Init(database *db.Database, tracer trace.Tracer) // initialises any connections/configurations
	GetRelativePath() string
	GetHandlers() *[]HandlerInterface // retrieve handlers, for attaching to the server and documentation
}

// Repo represents a collection of APIs around one entity.
type Repo[T any] struct {
	RelativePath string             // path that the repo should be accessible at
	Handlers     []HandlerInterface // list of handlers
	DB           *db.Database       // database object; initialised by the server
	tracer       trace.Tracer       // tracer object; initialised by the server
}

var _ Interface = &Repo[types.Nil]{}

// Init is called internally by the server when the Repo is attached to the server,
// to set up the database and tracer instance.
func (r *Repo[T]) Init(database *db.Database, tracer trace.Tracer) {
	r.DB = database
	r.tracer = tracer
}

// GetRelativePath returns the path that the repo is configured to listen to.
// Intended to be used by the Server object when attaching a Repo.
func (r *Repo[T]) GetRelativePath() string {
	return r.RelativePath
}

// GetHandlers returns the list of handlers added to the repo.
// Intended to be used by the Server object when attaching a Repo.
func (r *Repo[T]) GetHandlers() *[]HandlerInterface {
	return &r.Handlers
}

// AddHandler attaches a handler that does not require any parsing of a payload
// (or any parsing will be done within the handler function).
func (r *Repo[T]) AddHandler(method, relativePath string, fn func(*gin.Context)) {
	log.Get("REPO").Debug("Adding " + method + " handler at " + relativePath + "...")
	handler := HandlerFuncNoPayload[T]{
		Method:       method,
		RelativePath: relativePath,
		Func:         fn,
	}
	r.Handlers = append(r.Handlers, &handler)
	log.Get("REPO").Debug("Now have " + strconv.Itoa(len(r.Handlers)) + " handlers attached.")
}

// AddPayloadHandler attaches a handler that requires parsing of a payload. Will attempt to
// bind the request body to a new instance of T (generic type which is set when instantiating
// the repo), and automatically rejects the request if there are any bind errors.
// Uses Gin's validator, referenced at https://github.com/gin-gonic/gin#model-binding-and-validation.
func (r *Repo[T]) AddPayloadHandler(method, relativePath string, fn func(*gin.Context, T)) {
	log.Get("REPO").Debug("Adding " + method + " handler at " + relativePath + "...")
	handler := HandlerFuncWithPayload[T]{
		Method:       method,
		RelativePath: relativePath,
		Func:         fn,
	}
	r.Handlers = append(r.Handlers, &handler)
	log.Get("REPO").Debug("Now have " + strconv.Itoa(len(r.Handlers)) + " handlers attached.")
}

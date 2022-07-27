package repo

import (
	"github.com/kaphos/webapp/internal/db"
	"go/types"
)

// RepoI defines the expected functions that a Server expects any
// attached repositories to have.
type RepoI interface {
	HTTPBaseI
	Init(database *db.Database)   // initialises any connections/configurations
	GetHandlers() *[]HandlerBaseI // retrieve handlers, for attaching to the server and documentation
}

// Repo represents a collection of APIs around one entity.
// Should implement RepoI.
type Repo[T any] struct {
	httpBase
	DB       *db.Database   // database object; initialised by the server
	Handlers []HandlerBaseI // list of handlers
}

var _ RepoI = &Repo[types.Nil]{}

// Init is called internally by the server when the Repo is attached to the server,
// to set up the database and tracer instance.
func (r *Repo[T]) Init(database *db.Database) {
	r.DB = database
}

// GetHandlers returns the list of handlers added to the repo.
// Intended to be used by the Server object when attaching a Repo.
func (r *Repo[T]) GetHandlers() *[]HandlerBaseI {
	return &r.Handlers
}

func (r *Repo[T]) AddHandler(h HandlerBaseI) {
	r.Handlers = append(r.Handlers, h)
}

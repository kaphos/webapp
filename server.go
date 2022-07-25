package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/db"
	"github.com/kaphos/webapp/internal/log"
	"github.com/kaphos/webapp/internal/telemetry"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/repo"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os"
)

var routerLogger = log.Get("ROUTE")

type Server struct {
	Logger    *zap.Logger
	tracer    trace.Tracer
	db        *db.Database
	router    *gin.Engine
	apiRouter *gin.RouterGroup
}

// NewServer returns a new Server object, while performing
// all initialisation as required (Sentry, tracing, database).
func NewServer(appName string, dbUser, dbPass string, dbConns int32) (Server, error) {
	// Initialise Sentry first, so that any errors that come up can be flagged
	errchk.InitSentry()

	server := Server{
		Logger: log.Get("MAIN"),
		tracer: telemetry.NewTracer(appName, "server"),
	}

	var err error
	server.db, err = db.NewDB(appName, dbUser, dbPass, dbConns)
	if errchk.HaveError(err, "initDB") {
		return Server{}, err
	}

	server.buildRouter()

	return server, nil
}

// Attach a Repo to the server. Initialises the repository by passing in the database connection
// and a tracer object, and adds each of the repository's handlers to the server's Gin engine.
func (s *Server) Attach(repo repo.Interface) {
	s.Logger.Debug("Attaching repo " + repo.GetRelativePath())
	repo.Init(s.db)

	group := s.apiRouter.Group(repo.GetRelativePath(), *repo.GetMiddleware()...)

	for _, handler := range *repo.GetHandlers() {
		s.Logger.Debug(" - Attaching " + handler.GetMethod() + " handler at " + handler.GetRelativePath())
		group.Handle(handler.GetMethod(), handler.GetRelativePath(), handler.GetHandlers()...)
	}
}

// Start the Gin engine/router.
func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	s.Logger.Info("Listening on port " + port)

	return s.router.Run(":" + port)
}

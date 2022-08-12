package webapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/log"
	"github.com/kaphos/webapp/internal/swagger"
	"github.com/kaphos/webapp/internal/telemetry"
	"github.com/kaphos/webapp/pkg/db"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/keycloak"
	"github.com/kaphos/webapp/pkg/repo"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os"
	"regexp"
)

var routerLogger = log.Get("ROUTE")

type Server struct {
	logger    *zap.Logger
	tracer    trace.Tracer
	db        *db.Database
	router    *gin.Engine
	apiRouter *gin.RouterGroup
	apiDocs   *swagger.OpenAPI
}

// NewServer returns a new Server object, while performing
// all initialisation as required (Sentry, tracing, database).
func NewServer(appName, version, dbUser, dbPass string, dbConns int32) (Server, error) {
	// Initialise Sentry first, so that any errors that come up can be flagged
	errchk.InitSentry()

	apiDocs := swagger.Generate(appName, version)

	server := Server{
		logger:  log.Get("MAIN"),
		tracer:  telemetry.NewTracer(appName, "server"),
		apiDocs: &apiDocs,
	}

	var err error
	server.db, err = db.NewDB(appName, dbUser, dbPass, dbConns)
	if errchk.HaveError(err, "initDB") {
		return Server{}, err
	}

	server.buildRouter()

	return server, nil
}

var pathRegexp = regexp.MustCompile("//+")

func buildPath(r repo.HTTPBaseI, h repo.HTTPBaseI) string {
	path := "/" + r.RelativePath() + "/" + h.RelativePath() + "/"
	path = pathRegexp.ReplaceAllString(path, "/")
	return path
}

func (s *Server) addAPIPath(r repo.RepoI, h repo.HandlerBaseI, path, summary, description string,
	params map[string]swagger.SimpleParam) {

	// Build the list of potential responses by both the repo and handler.
	responses := make(map[int]swagger.Response)
	for code, resp := range r.Responses() {
		responses[code] = resp
	}

	for code, resp := range h.Responses() { // Handler after repo, to overwrite anything that may have been declared
		responses[code] = resp
	}

	s.apiDocs.AddPath(h.Type(), r.RelativePath(), h.Method(), path, summary, description, params, responses)
}

// Attach a Repo to the server. Initialises the repository by passing in the database connection
// and a tracer object, and adds each of the repository's handlers to the server's Gin engine.
func (s *Server) Attach(r repo.RepoI) {
	s.logger.Debug(fmt.Sprintf("Attaching repo \"%s\"", r.RelativePath()))
	r.Init(s.db)

	group := s.apiRouter.Group(r.RelativePath(), *r.Middleware()...)

	for _, h := range *r.GetHandlers() {
		path := buildPath(r, h)
		s.logger.Debug(fmt.Sprintf(" - Attaching handler at \"%s\" (%s)", path, h.Method()))

		handlers := make([]gin.HandlerFunc, 0)
		handlers = append(handlers, *h.Middleware()...)
		handlers = append(handlers, h.Handle)
		group.Handle(h.Method(), h.RelativePath(), handlers...)

		// Build Swagger API
		s.addAPIPath(r, h, path, h.Summary(), h.Description(), h.Params())
	}
}

// Start the Gin engine/router.
func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	s.logger.Info("Listening on port " + port)

	return s.router.Run(":" + port)
}

func (s *Server) NewKC(pk string) (keycloak.Keycloak, error) {
	return keycloak.New(pk, s.db)
}

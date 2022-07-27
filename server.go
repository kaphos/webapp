package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/db"
	"github.com/kaphos/webapp/internal/log"
	"github.com/kaphos/webapp/internal/swagger"
	"github.com/kaphos/webapp/internal/telemetry"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/repo"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os"
	"regexp"
)

var routerLogger = log.Get("ROUTE")

type Server struct {
	appName   string
	Logger    *zap.Logger
	tracer    trace.Tracer
	db        *db.Database
	router    *gin.Engine
	apiRouter *gin.RouterGroup
	APIDocs   *swagger.OpenAPI
}

// NewServer returns a new Server object, while performing
// all initialisation as required (Sentry, tracing, database).
func NewServer(appName, dbUser, dbPass string, dbConns int32) (Server, error) {
	// Initialise Sentry first, so that any errors that come up can be flagged
	errchk.InitSentry()

	apiDocs := swagger.Generate(appName, "v1.0.0")

	server := Server{
		appName: appName,
		Logger:  log.Get("MAIN"),
		tracer:  telemetry.NewTracer(appName, "server"),
		APIDocs: &apiDocs,
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
	path := "/" + r.GetRelativePath() + "/" + h.GetRelativePath() + "/"
	path = pathRegexp.ReplaceAllString(path, "/")
	return path
}

func (s *Server) addAPIPath(r repo.RepoI, h repo.HandlerBaseI, path string) {
	// Form the request body, based on the type of the handler.
	reqBody := swagger.BuildRequestBody(h.GetType())

	// Build the list of potential responses by both the repo and handler.
	responses := make(map[int]swagger.Response)
	for code, resp := range r.GetResponses() {
		responses[code] = resp
	}

	for code, resp := range h.GetResponses() { // Handler after repo, to overwrite anything that may have been declared
		responses[code] = resp
	}

	s.APIDocs.AddPath(r.GetRelativePath(), h.GetMethod(), path, reqBody, responses)
}

// Attach a Repo to the server. Initialises the repository by passing in the database connection
// and a tracer object, and adds each of the repository's handlers to the server's Gin engine.
func (s *Server) Attach(r repo.RepoI) {
	s.Logger.Debug("Attaching repo " + r.GetRelativePath())
	r.Init(s.db)

	group := s.apiRouter.Group(r.GetRelativePath(), *r.GetMiddleware()...)

	for _, h := range *r.GetHandlers() {
		path := buildPath(r, h)
		s.Logger.Debug(" - Attaching " + h.GetMethod() + " handler at " + path)

		handlers := make([]gin.HandlerFunc, 0)
		handlers = append(handlers, *h.GetMiddleware()...)
		handlers = append(handlers, h.Handle)
		group.Handle(h.GetMethod(), h.GetRelativePath(), handlers...)

		// Build Swagger API
		s.addAPIPath(r, h, path)
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

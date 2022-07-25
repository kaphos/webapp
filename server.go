package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/db"
	"github.com/kaphos/webapp/internal/errorhandling"
	"github.com/kaphos/webapp/internal/log"
	"github.com/kaphos/webapp/internal/telemetry"
	"github.com/kaphos/webapp/pkg/repo"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Logger    *zap.Logger
	db        *db.Database
	tracer    trace.Tracer
	router    *gin.Engine
	apiRouter *gin.RouterGroup
}

// NewServer returns a new Server object, while performing
// all initialisation as required (Sentry, tracing, database).
func NewServer(appName string, dbUser, dbPass string, dbConns int32) (Server, error) {
	errorhandling.InitSentry() // Initialise Sentry first, so that any errors that come up can be flagged

	server := Server{
		Logger: log.Get("MAIN"),
		tracer: telemetry.NewTracer(appName, "main"),
	}

	var err error
	server.db, err = db.NewDB(appName, dbUser, dbPass, dbConns)
	if errorhandling.HaveError(err, "initDB") {
		return Server{}, err
	}

	server.router, server.apiRouter = buildRouter(appName)

	return server, nil
}

// Attach a Repo to the server. Initialises the repository by passing in the database connection
// and a tracer object, and adds each of the repository's handlers to the server's Gin engine.
func (s *Server) Attach(repo repo.Interface) {
	s.Logger.Debug("Attaching repo " + repo.GetRelativePath())
	repo.Init(s.db, s.tracer)

	for _, handler := range *repo.GetHandlers() {
		path := repo.GetRelativePath() + "/" + handler.GetRelativePath()
		s.Logger.Debug("Attaching " + handler.GetMethod() + " handler at " + path)
		s.apiRouter.Handle(handler.GetMethod(), path, handler.Handle)
	}
}

func buildRouter(appName string) (*gin.Engine, *gin.RouterGroup) {
	routerLogger := log.Get("GIN")
	routerTracer := telemetry.NewTracer(appName, "router")

	router := gin.New()

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	apiGroup := router.Group("/api")
	apiGroup.Use(func(c *gin.Context) {
		_, span := routerTracer.Start(c.Request.Context(), c.Request.Method+" "+c.Request.URL.Path)
		defer span.End()

		start := time.Now()
		c.Next()

		latency := time.Now().Sub(start)

		var sb strings.Builder
		sb.WriteString("(")
		sb.WriteString(latency.String())
		sb.WriteString(") ")
		sb.WriteString(c.Request.Method)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(c.Writer.Status()))
		sb.WriteString(" ")
		sb.WriteString(c.Request.URL.Path)

		routerLogger.Info(sb.String())
	})

	return router, apiGroup
}

// Start the Gin engine/router.
func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	return s.router.Run(":" + port)
}

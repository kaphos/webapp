package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/telemetry"
	"github.com/kaphos/webapp/pkg/errchk"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) loggerMiddleware(c *gin.Context) {
	_, span := s.tracer.Start(c.Request.Context(), c.Request.Method+" "+c.Request.URL.Path)
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
}

func (s *Server) buildRouter() {
	router := gin.New()

	// Auxiliary handlers
	router.GET("/metrics", gin.WrapF(telemetry.PromHandler.ServeHTTP))
	router.GET("/healthcheck", func(c *gin.Context) {
		if err := s.db.Healthcheck(c.Request.Context()); errchk.HaveError(err, "healthcheck") {
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(http.StatusOK)
		}
	})

	apiGroup := router.Group("/api")
	apiGroup.Use(s.loggerMiddleware)

	// TODO: Implement middleware support, and keycloak support

	s.router = router
	s.apiRouter = apiGroup
}

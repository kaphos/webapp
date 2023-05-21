package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/internal/telemetry"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (s *Server) loggerMiddleware(c *gin.Context) {
	_, span := s.tracer.Start(c.Request.Context(), c.Request.Method+" "+c.Request.URL.Path)
	defer span.End()

	start := time.Now()
	c.Next()

	latency := time.Since(start)
	method := c.Request.Method
	status := strconv.Itoa(c.Writer.Status())

	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(latency.String())
	sb.WriteString(") ")
	sb.WriteString(method)
	sb.WriteString(":")
	sb.WriteString(status)
	sb.WriteString(" ")
	sb.WriteString(c.Request.URL.Path)
	routerLogger.Info(sb.String())

	telemetry.PromLogRequest(method, status, latency.Seconds())
}

func (s *Server) buildRouter() {
	if os.Getenv("DEBUG") != "true" {
		// Hide debug messages, unless DEBUG flag is set
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Auxiliary handlers
	router.GET("/metrics", gin.WrapF(telemetry.PromHandler.ServeHTTP))
	router.GET("/healthcheck", func(c *gin.Context) {
		if err := s.DB.Healthcheck(c.Request.Context()); errchk.HaveError(err, "healthcheck") {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.Status(http.StatusOK)
		}
	})

	apiGroup := router.Group("/api")
	apiGroup.Use(s.loggerMiddleware)

	apiGroup.GET("/version", func(c *gin.Context) {
		c.String(200, utils.GetEnv("VERSION", "v0.0.0"))
	})

	s.Router = router
	s.apiRouter = apiGroup
}

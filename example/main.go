package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp"
	"github.com/kaphos/webapp/pkg/middleware"
)

func main() {
	s := setupServer()
	_ = s.GenDocs([]webapp.APIServer{{URL: "http://localhost:5000", Description: "Dev server"}}, "swagger.yml")
	if err := s.Start(); err != nil {
		return
	}
}

func setupServer() *webapp.Server {
	s, err := webapp.NewServer("Test App", "v1", "testuser", "testpass", 1)
	if err != nil {
		return nil
	}

	authMiddleware := setupAuthMiddleware()
	s.Attach(buildPingRepo())
	s.Attach(buildUserRepo())
	s.Attach(buildItemRepo(authMiddleware))
	return &s
}

func setupAuthMiddleware() middleware.Middleware {
	return middleware.NewAuth(func(c *gin.Context) bool {
		return c.GetHeader("auth") == "true"
	})
}

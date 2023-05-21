package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp"
	"github.com/kaphos/webapp/pkg/middleware"
)

const pk = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAr9HwDLMZVt6BgC29QeRNT81V8HY7qZWIKAzLpDJv8G4c/O2jEiw+2TyvW7J4bP2sl6EGyrXsPg2eO5h4vQDJs59MQvqd5cvi/oE+ZV7BZozvZ91FWJGL1GvE4Pn+twDWPIl4mm/cs1UmyqVmqct3Yxb4j7TeQ8Th9K/a9BkIvKmRpJ0Qug7ZinhddYzqzC2tgtBzcPflYw2ZrlzqUYxnr4OaV5nKwBTeuSIow2S6F1nW6nymMrp6CvKbtXZj0XSTcYWazekwRVoDZnWe9Jk3UMapLl9zSy8WeKx4NUcNfq4gqNOcUmZoH3xhIXSjWftJpIgsp5tFzjpL3LtwZiUZQwIDAQAB
-----END PUBLIC KEY-----`

func main() {
	s := setupServer()
	_ = s.GenDocs([]webapp.APIServer{{"http://localhost:5000", "Dev server"}}, "swagger.yml")
	if err := s.Start(); err != nil {
		return
	}
}

func setupServer() *webapp.Server {
	s, err := webapp.NewServer("Test App", "v1", "testuser", "testpass", 1)
	if err != nil {
		return nil
	}

	//authMiddleware := setupAuthMiddleware()
	s.Attach(buildUserRepo())
	//s.Attach(buildItemRepo(authMiddleware))
	return &s
}

func setupAuthMiddleware() middleware.Middleware {
	return middleware.NewAuth(func(ctx *gin.Context) bool {
		return true
	})
}

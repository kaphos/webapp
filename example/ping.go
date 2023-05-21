package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/pkg/handler"
	"github.com/kaphos/webapp/pkg/repo"
	"go/types"
	"net/http"
)

type PingRepo struct{ repo.Repo[types.Nil] }

func (r *PingRepo) ping(c *gin.Context) bool {
	c.String(http.StatusOK, "pong")
	return true
}

func buildPingRepo() repo.RepoI {
	r := PingRepo{}
	r.SetRelativePath("ping")
	h := handler.NewU("GET", "/", r.ping, 200, "pong")
	r.AddHandler(&h)
	return &r
}

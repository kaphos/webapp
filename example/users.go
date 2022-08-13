package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/repo"
	"go/types"
)

type UserRepo struct{ repo.Repo[types.Nil] }

func (r *UserRepo) login(c *gin.Context) bool {
	return true
}

func (r *UserRepo) add(c *gin.Context) bool {
	value, _ := c.Get("kc-sub")
	kcId := value.(string)

	err := r.DB.Exec("addUser", c.Request.Context(), `INSERT INTO users (kc_sub) VALUES ($1)`, kcId)
	errchk.Check(err, "addUser")

	return true
}

func buildUserRepo(authMiddleware repo.Middleware) repo.RepoI {
	r := UserRepo{}
	r.SetRelativePath("user")

	h := repo.NewHandlerU("POST", "/login", r.login, 200, nil, authMiddleware)
	r.AddHandler(&h)

	c := repo.NewHandlerU("PUT", "/add", r.add, 201, nil, authMiddleware)
	r.AddHandler(&c)

	return &r
}

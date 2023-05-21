package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/pkg/handler"
	"github.com/kaphos/webapp/pkg/repo"
	"net/http"
)

type User struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Admin  bool    `json:"admin"`
	Groups int     `json:"groups"`
	Age    float32 `json:"age"`
}

type UserRepo struct{ repo.Repo[User] }

func (r *UserRepo) getAll(c *gin.Context) bool {
	rows, cancel, err := r.DB.Query("getUsers", c.Request.Context(), `SELECT id, name, email, admin, groups, age FROM users`)
	defer cancel()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return false
	}

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Admin, &user.Groups, &user.Age)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return false
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
	return true
}

func buildUserRepo() repo.RepoI {
	r := UserRepo{}
	r.SetRelativePath("users")

	getUsersHandler := handler.NewU("GET", "", r.getAll, 200, make([]User, 0))
	r.AddHandler(&getUsersHandler)

	return &r
}

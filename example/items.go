package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/kaphos/webapp/pkg/handler"
	"github.com/kaphos/webapp/pkg/middleware"
	"github.com/kaphos/webapp/pkg/repo"
	"gopkg.in/guregu/null.v4"
	"net/http"
	"time"
)

type Item struct {
	ID      uuid.UUID   `json:"id" binding:"-"`
	Created time.Time   `json:"created"`
	Edited  null.Time   `json:"edited"`
	Name    string      `json:"name" binding:"required"`
	Owner   null.String `json:"owner"`
	Found   null.Bool   `json:"found"`
	Count   null.Int    `json:"count"`
	Price   null.Float  `json:"price"`
}

type ItemRepo struct{ repo.Repo[Item] }

func (r *ItemRepo) getItems(c *gin.Context) bool {
	c.JSON(http.StatusOK, make([]Item, 0))
	return true
}

func (r *ItemRepo) createItem(c *gin.Context, item Item) bool {
	fmt.Printf("%+v\n", item)
	return true
}

func buildItemRepo(authMiddleware middleware.Middleware) repo.RepoI {
	r := ItemRepo{}
	r.SetRelativePath("items")

	h := handler.NewU("GET", "/", r.getItems, 200, []Item{{}})
	h.SetSummary("Retrieves the list of items stored in the database.")
	h.SetDescription("Simply fetches all items.")
	r.AddHandler(&h)

	c := handler.NewP("POST", "/", r.createItem, 201, nil, authMiddleware)
	c.SetSummary("Creates a new item.")
	c.SetDescription("Only allowed by authenticated users.")
	r.AddHandler(&c)

	return &r
}

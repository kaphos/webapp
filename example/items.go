package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kaphos/webapp/pkg/errchk"
	"github.com/kaphos/webapp/pkg/repo"
	"net/http"
)

type Item struct {
	ID    int    `form:"id" binding:"-"`
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" example:"hello@email.com"`
}

type ItemRepo struct{ repo.Repo[Item] }

func (r *ItemRepo) getItems(c *gin.Context) bool {
	items := make([]Item, 0)

	rows, cancel, err := r.DB.Query("getItems", c.Request.Context(), "SELECT id, name, email FROM items")
	defer cancel()
	if err != nil {
		return false
	}

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Email)
		if errchk.HaveError(err, "getItems1") {
			return false
		}

		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)

	return true
}

func (r *ItemRepo) createItem(c *gin.Context, item Item) bool {
	fmt.Printf("%+v\n", item)
	return true
}

func (r *ItemRepo) updateItem(c *gin.Context, item Item) bool {
	fmt.Printf("%+v\n", item)
	return true
}

func buildItemRepo(authMiddleware repo.Middleware) repo.RepoI {
	r := ItemRepo{}
	r.SetRelativePath("item")

	h := repo.NewHandlerU("GET", "/", r.getItems, 200, []Item{{}})
	h.SetSummary("Retrieves the list of items stored in the database.")
	r.AddHandler(&h)

	c := repo.NewHandlerP("POST", "/", r.createItem, 201, nil)
	c.SetSummary("Creates a new item.")
	c.SetDescription("Only allowed by authenticated users.")
	r.AddHandler(&c)

	p := repo.NewHandlerP("PUT", "/:id", r.updateItem, 200, nil)
	p.AddParam("id", "integer", "ID of item that is to be updated")
	p.SetSummary("Update item.")
	r.AddHandler(&p)

	return &r
}

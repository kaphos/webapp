package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/kaphos/webapp/pkg/handler"
	"github.com/kaphos/webapp/pkg/middleware"
	"github.com/kaphos/webapp/pkg/repo"
	"gopkg.in/guregu/null.v4"
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
	//items := make([]Item, 0)
	//
	//rows, cancel, err := r.DB.Query("getItems", c.Request.Context(), "SELECT id, name FROM items")
	//defer cancel()
	//if err != nil {
	//	return false
	//}
	//
	//for rows.Next() {
	//	var item Item
	//	err := rows.Scan(&item.ID, &item.Name, &item.Email)
	//	if errchk.HaveError(err, "getItems1") {
	//		return false
	//	}
	//
	//	items = append(items, item)
	//}
	//
	//c.JSON(http.StatusOK, items)

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

func buildItemRepo(authMiddleware middleware.Middleware) repo.RepoI {
	r := ItemRepo{}
	r.SetRelativePath("item")

	h := handler.NewU("GET", "/", r.getItems, 200, []Item{{}})
	h.SetSummary("Retrieves the list of items stored in the database.")
	r.AddHandler(&h)

	c := handler.NewP("POST", "/", r.createItem, 201, nil, authMiddleware)
	c.SetSummary("Creates a new item.")
	c.SetDescription("Only allowed by authenticated users.")
	r.AddHandler(&c)

	p := handler.NewP("PUT", "/:id", r.updateItem, 200, nil)
	p.AddParam("id", "integer", "ID of item that is to be updated")
	p.SetSummary("Update item.")
	r.AddHandler(&p)

	return &r
}

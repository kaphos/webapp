package main

import (
	"encoding/json"
	"github.com/kaphos/webapp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	s, w := setup()
	defer teardown(s)
	t.Run("Fetch", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/users", nil)
		s.Router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var resp []User
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.ElementsMatch(t, resp, []User{
			{
				ID:     1,
				Name:   "John",
				Email:  "john@gmail.com",
				Admin:  true,
				Groups: 1,
				Age:    3.5,
			},
			{
				ID:     2,
				Name:   "Tom",
				Email:  "tom@email.com",
				Admin:  false,
				Groups: 2,
				Age:    5.2,
			},
			{
				ID:     3,
				Name:   "Jane",
				Email:  "jane@hotmail.com",
				Admin:  false,
				Groups: 3,
				Age:    8.9,
			},
		})
	})
}

func setup() (*webapp.Server, *httptest.ResponseRecorder) {
	s := setupServer()
	w := httptest.NewRecorder()
	return s, w
}

func teardown(s *webapp.Server) {
	// remove data, etc.
}

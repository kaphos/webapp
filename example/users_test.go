package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetUsers(t *testing.T) {
	s, w := setup()
	req, _ := http.NewRequest("GET", "/api/users/", nil)
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
}

func TestAddUser(t *testing.T) {
	s, w := setup()
	newUser := User{
		Name:   "user-name",
		Email:  "user-email",
		Admin:  true,
		Groups: 1,
		Age:    48.3,
	}
	body, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/api/users/", bytes.NewReader(body))
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var resp int
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.Nil(t, err)
	assert.Greater(t, resp, 0)
}

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

type AddUserTestCase struct {
	name       string
	body       []byte
	statusCode int
}

func TestAddUser(t *testing.T) {
	validBody, _ := json.Marshal(User{
		Name:   "user-name",
		Email:  "hello@email.com",
		Admin:  true,
		Groups: 1,
		Age:    48.3,
	})

	invalidBody, _ := json.Marshal(User{
		Groups: 1,
		Age:    48.3,
	})

	testCases := []AddUserTestCase{
		{name: "ValidRequest", body: validBody, statusCode: http.StatusCreated},
		{name: "IncompletePayload", body: invalidBody, statusCode: http.StatusBadRequest},
		{name: "MissingBody", body: nil, statusCode: http.StatusBadRequest},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			s, w := setup()
			req, _ := http.NewRequest("POST", "/api/users/", bytes.NewReader(testCase.body))
			s.Router.ServeHTTP(w, req)
			assert.Equal(t, testCase.statusCode, w.Code)

			if testCase.statusCode < 300 {
				var resp int
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Greater(t, resp, 0)
			}
		})
	}
}

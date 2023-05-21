package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetItems(t *testing.T) {
	s, w := setup()
	req, _ := http.NewRequest("GET", "/api/items/", nil)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp []Item
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.Nil(t, err)
}

type AddItemTestCase struct {
	authValue  string
	statusCode int
}

func TestAddItem(t *testing.T) {
	testCases := []AddItemTestCase{
		{authValue: "true", statusCode: http.StatusBadRequest},
		{authValue: "false", statusCode: http.StatusUnauthorized},
	}

	for _, testCase := range testCases {
		t.Run("Header-"+testCase.authValue, func(t *testing.T) {
			s, w := setup()
			req, _ := http.NewRequest("POST", "/api/items/", nil)
			req.Header.Add("auth", testCase.authValue)
			s.Router.ServeHTTP(w, req)
			assert.Equal(t, testCase.statusCode, w.Code)
		})
	}
}

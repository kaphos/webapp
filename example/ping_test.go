package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	s, w := setup()
	req, _ := http.NewRequest("GET", "/api/ping/", nil)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp string
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.Nil(t, err)
	assert.Equal(t, resp, "pong")
}

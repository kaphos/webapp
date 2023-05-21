package main

import (
	"github.com/kaphos/webapp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestHealthcheck(t *testing.T) {
	s, w := setup()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func setup() (*webapp.Server, *httptest.ResponseRecorder) {
	s := setupServer()
	w := httptest.NewRecorder()
	return s, w
}

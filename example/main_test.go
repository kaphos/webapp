package main

import (
	"github.com/kaphos/webapp"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func setup() (*webapp.Server, *httptest.ResponseRecorder) {
	s := setupServer()
	w := httptest.NewRecorder()
	return s, w
}

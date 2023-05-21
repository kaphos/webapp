package main

import (
	"github.com/kaphos/webapp"
	"github.com/stretchr/testify/assert"
	"io"
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

type VersionTestCase struct {
	name          string
	setValue      string
	expectedValue string
}

func TestVersion(t *testing.T) {
	for _, testCase := range []VersionTestCase{
		{"present", "v1.4.8", "v1.4.8"},
		{"absent", "", "v0.0.0"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("VERSION", testCase.setValue)
			s, w := setup()
			req, _ := http.NewRequest("GET", "/api/version", nil)
			s.Router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			body, _ := io.ReadAll(w.Body)
			assert.Equal(t, testCase.expectedValue, string(body))
		})
	}
}

func setup() (*webapp.Server, *httptest.ResponseRecorder) {
	s := setupServer()
	w := httptest.NewRecorder()
	return s, w
}

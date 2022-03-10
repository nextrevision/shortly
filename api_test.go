package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/-/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	api := newTestServer()

	w := httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code, "expected 204 status response")
}

func TestShortenUrlHandler(t *testing.T) {
	data := url.Values{}
	data.Set("url", "http://foo")

	req, err := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	api := newTestServer()

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(api.ShortenUrlHandler)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "expected 200 status response")
	assert.NotEmpty(t, w.Body.String(), "expected non-empty response body")
}

func TestShortenUrlHandlerUnsupportedMediaType(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	api := newTestServer()

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(api.ShortenUrlHandler)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code, "expected 415 status response")
}

func TestShortenUrlHandlerEmptyPayload(t *testing.T) {
	data := url.Values{}

	req, err := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	api := newTestServer()

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(api.ShortenUrlHandler)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "expected 400 status response")
}

func TestRedirectUrlHandler(t *testing.T) {
	expectedUrl := "http://foo"

	api := newTestServer()
	id, err := api.service.SaveUrl(expectedUrl)
	assert.Nil(t, err, "error should be nil")
	assert.NotEmpty(t, id, "id should not be empty")

	req, err := http.NewRequest("GET", "/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMovedPermanently, w.Code, "expected 301 status response")
	assert.Equal(t, expectedUrl, w.Header().Get("Location"), "expected location to correctly set")
}

func TestRedirectUrlHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Fatal(err)
	}

	api := newTestServer()

	w := httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "expected 404 status response")
}

func newTestServer() *Api {
	service := NewUrlService(NewMemoryDataStore())
	return NewApi(service)
}

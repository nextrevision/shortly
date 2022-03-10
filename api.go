package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// The Api struct is the primary entrypoint for routing and handling HTTP requests
type Api struct {
	Router  *mux.Router
	service *UrlService
}

func NewApi(service *UrlService) *Api {
	api := &Api{service: service}

	// Create http router and register paths
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	r.HandleFunc("/", api.ShortenUrlHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", api.RedirectUrlHandler).Methods(http.MethodGet)
	r.HandleFunc("/echo/{id}", api.EchoHandler).Methods(http.MethodGet)
	r.Use(httpMetricsMiddleware)
	r.Use(loggingMiddleware)
	api.Router = r

	return api
}

// ShortenUrlHandler receives a form-encoded body in a url="value" format and shortens it and returns a distinct id
func (a *Api) ShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := r.Form.Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := a.service.SaveUrl(url)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"url":   url,
		}).Errorf("error saving url %s", url)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"id":  id,
		"url": url,
	}).Infof("shortened url to %s", id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s/%s", r.Host, id)))
}

// RedirectUrlHandler receives a url with a shortened id in the path and resolves and redirects to it's reference url
func (a *Api) RedirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := a.service.ResolveUrl(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"id":  id,
		"url": url,
	}).Infof("redirected url with id %s", id)

	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusMovedPermanently)
}

// EchoHandler is used when testing the url shortener workflow and simply echos the id variable it's called with
func (a *Api) EchoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("ID: %s", vars["id"])))
}
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	httpRequestDurationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1},
	}, []string{"path", "method", "status"})
	httpRequestStatusMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_status",
		Help: "Duration of HTTP requests.",
	}, []string{"path", "method", "status"})
	httpRequestTotalMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Duration of HTTP requests.",
	})
)

type responseData struct {
	status int
}

// http.ResponseWriter implementation to record status codes
type responseWrapper struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *responseWrapper) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	return size, err
}

func (r *responseWrapper) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Records responses as prom metrics with obfuscated responses and paths to reduce cardinality
func httpMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseData := &responseData{
			status: 200,
		}
		rw := responseWrapper{
			ResponseWriter: w,
			responseData:   responseData,
		}

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		start := time.Now()

		next.ServeHTTP(&rw, r)

		statusGroup := fmt.Sprintf("%dxx", rw.responseData.status/100)
		httpRequestDurationMetric.WithLabelValues(path, r.Method, statusGroup).Observe(time.Since(start).Seconds())
		httpRequestStatusMetric.WithLabelValues(path, r.Method, statusGroup).Inc()
		httpRequestTotalMetric.Inc()
	})
}

// Logs all http requests with some basic context
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.RequestURI,
			"host":   r.Host,
		}).Infof("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

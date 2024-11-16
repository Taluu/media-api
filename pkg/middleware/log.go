package middleware

import (
	"log"
	"net/http"
	"time"
)

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *customResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extendedWriter := &customResponseWriter{ResponseWriter: w}
		defer func(start time.Time) {
			log.Printf("HTTP %s %s status %d, took %s", r.Method, r.URL.Path, extendedWriter.statusCode, time.Since(start))
		}(time.Now())

		log.Printf("Received HTTP %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(extendedWriter, r)
	})
}

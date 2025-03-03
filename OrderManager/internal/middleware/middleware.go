package middleware

import (
	"bytes"
	"io"
	"net/http"
	"orderservice/internal/logger"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		logger.InfoLogger.Printf("Request: %s %s\nHeaders: %v\nBody: %s",
			r.Method, r.URL.Path, r.Header, string(requestBody))

		// Create custom response writer to capture response
		buf := &bytes.Buffer{}
		rw := &responseWriter{ResponseWriter: w, body: buf}

		next.ServeHTTP(rw, r)

		// Log response
		duration := time.Since(start)
		logger.InfoLogger.Printf("Response: %s %s\nStatus: %d\nDuration: %v\nBody: %s",
			r.Method, r.URL.Path, http.StatusOK, duration, rw.body.String())
	})
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

package middlewares

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs information about incoming requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log the request details
		log.Printf(
			"[INFO] remote-ip=%s method=%s endpoint=%s status=%d duration=%v",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration.Round(time.Microsecond),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Recover middleware recovers from panics and returns a 500 status code
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				fmt.Printf("panic: %v\n%s", err, debug.Stack())

				// Return Internal Server Error
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

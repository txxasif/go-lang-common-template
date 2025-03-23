package middleware

import (
	"net/http"
	"time"
)

// RequestLogger logs information about each HTTP request
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log request details
		duration := time.Since(start)
		// In a real application, you would use a proper logger
		// log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, duration)
	})
}

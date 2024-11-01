package middleware

import (
	"context"
	"net/http"
	"time"
)

// SetRequestContextWithTimeout applies a timeout to the request context.
func SetRequestContextWithTimeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a new context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			// Attach the new context with timeout to the request
			r = r.WithContext(ctx)

			// Serve the next handler
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func CORS(logger *zap.Logger, allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			fmt.Println(allowedOrigins)
			// Check if origin is allowed
			if isOriginAllowed(origin, allowedOrigins) {
				// Set CORS headers
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
				w.Header().Set("Access-Control-Expose-Headers", "Link")
				w.Header().Set("Access-Control-Max-Age", "300")
			}

			// Handle preflight OPTIONS request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		// Support wildcard
		if allowed == "*" {
			return true
		}
		// Exact match
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}

	return false
}

package middleware

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func CORS(logger *zap.Logger, allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := isOriginAllowed(origin, allowedOrigins)
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, Last-Event-ID")
				w.Header().Set("Access-Control-Expose-Headers", "Link, Content-Type, Cache-Control")
				w.Header().Set("Access-Control-Max-Age", "300")
				w.Header().Set("Vary", "Origin")
			} else if origin != "" {
				logger.Warn("CORS Origin not allowed",
					zap.String("origin", origin),
					zap.Any("allowed_origins", allowedOrigins),
				)
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
		// Flexible match: if allowed is "loco.prabalesh.com" but browser sends "https://loco.prabalesh.com"
		if !strings.HasPrefix(allowed, "http://") && !strings.HasPrefix(allowed, "https://") {
			if strings.HasSuffix(origin, "://"+allowed) || strings.HasSuffix(origin, "://"+allowed+"/") {
				return true
			}
		}

		// Development convenience: Allow any localhost/127.0.0.1 if any local dev origin is in the list
		if strings.Contains(allowed, "localhost") || strings.Contains(allowed, "127.0.0.1") {
			if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") ||
				origin == "http://localhost" || origin == "http://127.0.0.1" {
				return true
			}
		}
	}

	return false
}

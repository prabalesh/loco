package middleware

import (
	"context"
	"net/http"

	"github.com/prabalesh/loco/backend/internal/infrastructure/auth"
	"go.uber.org/zap"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserRoleKey  contextKey = "user_role"
)

// Auth middleware validates JWT token from cookie and adds user info to context
func Auth(jwtService *auth.JWTService, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get access token from cookie
			cookie, err := r.Cookie("accessToken")
			if err != nil {
				logger.Warn("No access token cookie found", zap.Error(err))
				respondUnauthorized(w, "unauthorized: no access token")
				return
			}

			token := cookie.Value
			if token == "" {
				logger.Warn("Access token cookie is empty")
				respondUnauthorized(w, "unauthorized: empty access token")
				return
			}

			// Validate token
			claims, err := jwtService.ValidateToken(token, false)
			if err != nil {
				logger.Warn("Invalid access token",
					zap.Error(err),
					zap.String("token_prefix", getSafeTokenPrefix(token)),
				)
				respondUnauthorized(w, "unauthorized: invalid token")
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			// Log successful authentication
			logger.Debug("Request authenticated",
				zap.Int("user_id", claims.UserID),
				zap.String("email", claims.Email),
				zap.String("role", claims.Role),
				zap.String("path", r.URL.Path),
			)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RegularOrAdminAuth middleware validates either regular accessToken OR adminAccessToken
func RegularOrAdminAuth(jwtService *auth.JWTService, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			isAdmin := false

			// 1. Try regular access token
			if cookie, err := r.Cookie("accessToken"); err == nil {
				token = cookie.Value
			}

			// 2. If not found, try admin access token
			if token == "" {
				if cookie, err := r.Cookie("adminAccessToken"); err == nil {
					token = cookie.Value
					isAdmin = true
				}
			}

			if token == "" {
				logger.Warn("No regular or admin access token cookie found")
				respondUnauthorized(w, "unauthorized: no access token")
				return
			}

			// Validate token
			claims, err := jwtService.ValidateToken(token, false)
			if err != nil {
				logger.Warn("Invalid access token",
					zap.Error(err),
					zap.Bool("is_admin_cookie", isAdmin),
				)
				respondUnauthorized(w, "unauthorized: invalid token")
				return
			}

			// If it was an admin cookie, ensure it has the admin role
			if isAdmin && claims.Role != "admin" {
				logger.Warn("Admin cookie used but role is not admin", zap.Int("user_id", claims.UserID))
				respondForbidden(w, "forbidden: admin access required")
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth middleware tries to authenticate but doesn't fail if no token
func OptionalAuth(jwtService *auth.JWTService, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get access token from cookie
			cookie, err := r.Cookie("accessToken")
			if err != nil {
				// No token is fine for optional auth
				next.ServeHTTP(w, r)
				return
			}

			token := cookie.Value
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate token
			claims, err := jwtService.ValidateToken(token, false)
			if err != nil {
				// Invalid token is fine for optional auth
				logger.Debug("Optional auth: invalid token", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get role from context (must run after Auth middleware)
			role, ok := r.Context().Value(UserRoleKey).(string)
			if !ok {
				respondForbidden(w, "forbidden: role not found in context")
				return
			}

			// Check if user has one of the allowed roles
			hasRole := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				respondForbidden(w, "forbidden: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions

func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

func respondForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

// getSafeTokenPrefix returns first 10 chars of token for logging (safely)
func getSafeTokenPrefix(token string) string {
	if len(token) > 10 {
		return token[:10] + "..."
	}
	return token
}

// Context helper functions for handlers

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// GetUserEmail extracts user email from context
func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// GetUserRole extracts user role from context
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}

// MustGetUserID extracts user ID or panics (use in protected routes)
func MustGetUserID(ctx context.Context) int {
	userID, ok := GetUserID(ctx)
	if !ok {
		panic("user ID not found in context")
	}
	return userID
}

package middleware

import (
	"context"
	"net/http"

	"github.com/prabalesh/loco/backend/internal/infrastructure/auth"
	"go.uber.org/zap"
)

// RequireAdminAuth validates adminAccessToken and adds user info to context
func RequireAdminAuth(jwtService *auth.JWTService, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get ADMIN access token cookie
			cookie, err := r.Cookie("adminAccessToken")
			if err != nil {
				logger.Warn("No adminAccessToken cookie found", zap.Error(err))
				respondUnauthorized(w, "unauthorized: no admin access token")
				return
			}

			token := cookie.Value
			if token == "" {
				logger.Warn("adminAccessToken cookie is empty")
				respondUnauthorized(w, "unauthorized: empty admin access token")
				return
			}

			// Validate admin access token
			claims, err := jwtService.ValidateToken(token, false)
			if err != nil {
				logger.Warn("Invalid adminAccessToken",
					zap.Error(err),
					zap.String("token_prefix", getSafeTokenPrefix(token)),
				)
				respondUnauthorized(w, "unauthorized: invalid admin token")
				return
			}

			// Verify admin role
			if claims.Role != "admin" {
				logger.Warn("Non-admin attempted admin route",
					zap.Int("user_id", claims.UserID),
					zap.String("role", claims.Role),
				)
				respondForbidden(w, "forbidden: admin access required")
				return
			}

			// Add user info to context (same keys as regular Auth)
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			// Log successful admin authentication
			logger.Debug("Admin request authenticated",
				zap.Int("user_id", claims.UserID),
				zap.String("email", claims.Email),
				zap.String("role", claims.Role),
				zap.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

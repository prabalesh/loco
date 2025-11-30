package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prabalesh/loco/backend/internal/delivery/handler"
	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/infrastructure/auth"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
)

type Dependencies struct {
	Log         *zap.Logger
	Cfg         *config.Config
	Db          *database.Database
	JWTService  *auth.JWTService
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler

	AdminHandler     *handler.AdminHandler
	AdminAuthHandler *handler.AdminAuthHandler
}

func SetupRouter(deps *Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler())

	// auth handler
	mux.HandleFunc("POST /auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /auth/verify-email", deps.AuthHandler.VerifyEmail)
	mux.HandleFunc("POST /auth/resend-verification", deps.AuthHandler.ResendVerificationEmail)
	mux.HandleFunc("POST /auth/refresh", deps.AuthHandler.RefreshToken)
	mux.HandleFunc("POST /auth/logout", deps.AuthHandler.Logout)

	mux.HandleFunc("POST /auth/forgot-password", deps.AuthHandler.ForgotPassword)
	mux.HandleFunc("POST /auth/reset-password", deps.AuthHandler.ResetPassword)

	// protected routes
	authMiddleware := middleware.Auth(deps.JWTService, deps.Log)

	mux.Handle("GET /auth/me", authMiddleware(http.HandlerFunc(deps.AuthHandler.GetMe)))

	// User profile routes
	mux.Handle("GET /users/me", authMiddleware(http.HandlerFunc(deps.UserHandler.GetProfile)))
	mux.Handle("GET /users/{username}", authMiddleware(http.HandlerFunc(deps.UserHandler.GetProfileByUsername)))

	// Other protected routes
	mux.Handle("GET /problems", authMiddleware(http.HandlerFunc(placeholderHandler("Problems"))))
	mux.Handle("GET /submissions", authMiddleware(http.HandlerFunc(placeholderHandler("Submissions"))))

	// Admin routes
	adminAuthMiddleware := middleware.RequireAdminAuth(deps.JWTService, deps.Log)

	mux.HandleFunc("POST /admin/auth/login", deps.AdminAuthHandler.AdminLogin)
	mux.HandleFunc("POST /admin/auth/logout", deps.AdminAuthHandler.AdminLogout)
	mux.HandleFunc("POST /admin/auth/refresh", deps.AdminAuthHandler.AdminRefreshToken)
	mux.Handle("GET /admin/auth/me", adminAuthMiddleware(http.HandlerFunc(deps.AdminAuthHandler.GetAdminProfile)))

	mux.Handle("GET /admin/users", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.ListUsers)))
	mux.Handle("GET /admin/users/{id}", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.GetUser)))
	mux.Handle("DELETE /admin/users/{id}", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.DeleteUser)))
	mux.Handle("PATCH /admin/users/{id}/role", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.UpdateUserRole)))
	mux.Handle("PATCH /admin/users/{id}/status", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.UpdateUserStatus)))
	mux.Handle("GET /admin/analytics", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.GetAnalytics)))

	handler := middleware.Logging(deps.Log)(mux)
	handler = middleware.CORS(deps.Log, deps.Cfg.CORS.AllowedOrigins)(handler)

	return handler
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	}
}

func placeholderHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"` + name + ` endpoint - coming soon"}`))
	}
}

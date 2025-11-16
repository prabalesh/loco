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
}

func SetupRouter(deps *Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler())

	handler := middleware.Logging(deps.Log)(mux)
	handler = middleware.CORS(deps.Log, deps.Cfg.CORS.AllowedOrigins)(handler)

	// auth handler
	mux.HandleFunc("POST /auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /auth/refresh", deps.AuthHandler.RefreshToken)
	mux.HandleFunc("POST /auth/logout", deps.AuthHandler.Logout)

	// protected routes
	authMiddleware := middleware.Auth(deps.JWTService, deps.Log)

	mux.Handle("GET /auth/me", authMiddleware(http.HandlerFunc(deps.AuthHandler.GetMe)))

	// User profile routes
	mux.Handle("GET /users/me", authMiddleware(http.HandlerFunc(deps.UserHandler.GetProfile)))
	mux.Handle("GET /users/{username}", authMiddleware(http.HandlerFunc(deps.UserHandler.GetProfileByUsername)))

	// Other protected routes
	mux.Handle("GET /problems", authMiddleware(http.HandlerFunc(placeholderHandler("Problems"))))
	mux.Handle("GET /submissions", authMiddleware(http.HandlerFunc(placeholderHandler("Submissions"))))

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

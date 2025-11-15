package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prabalesh/loco/backend/internal/delivery/handler"
	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"go.uber.org/zap"
)

type Dependencies struct {
	Log         *zap.Logger
	Cfg         *config.Config
	Db          *database.Database
	AuthHandler *handler.AuthHandler
}

func SetupRouter(deps *Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler())

	handler := middleware.Logging(deps.Log)(mux)

	// auth handler
	mux.HandleFunc("POST /auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /auth/refresh", deps.AuthHandler.RefreshToken)
	mux.HandleFunc("POST /auth/logout", deps.AuthHandler.Logout)

	return handler
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	}
}

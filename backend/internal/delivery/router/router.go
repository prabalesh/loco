package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

func SetupRouter(log *zap.Logger, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler())

	handler := middleware.Logging(log)(mux)
	return handler
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	}
}

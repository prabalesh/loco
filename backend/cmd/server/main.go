package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/prabalesh/loco/backend/internal/delivery/handler"
	"github.com/prabalesh/loco/backend/internal/delivery/router"
	"github.com/prabalesh/loco/backend/internal/infrastructure/auth"
	"github.com/prabalesh/loco/backend/internal/repository/postgres"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"github.com/prabalesh/loco/backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	config.InitConfig()
	cfg := config.GetConfig()

	if err := logger.InitLogger(cfg.Log.Level); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	log := logger.GetLogger()

	log.Info("=== Application Starting ===",
		zap.String("version", "1.0.0"),
		zap.String("port", cfg.Server.Port),
	)

	db, err := database.NewPostgresDB(cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	stats := db.DB.Stats()
	log.Info("Database connection pool initialized",
		zap.Int("open_connections", stats.OpenConnections),
		zap.Int("in_use", stats.InUse),
		zap.Int("idle", stats.Idle),
	)

	if err := database.RunMigrations(db.DB, "./migrations", log); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	version, dirty, err := database.GetMigrationVersion(db.DB, "./migrations")
	if err != nil {
		log.Warn("Could not get migration version", zap.Error(err))
	} else {
		log.Info("Current migration version",
			zap.Uint("version", version),
			zap.Bool("dirty", dirty),
		)
	}

	deps := initializeDependencies(db, log, cfg)
	router := router.SetupRouter(deps)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Server starting",
			zap.String("port", cfg.Server.Port),
			zap.String("address", "http://localhost:"+cfg.Server.Port),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("=== Application Stopped ===")
}

func initializeDependencies(db *database.Database, logger *zap.Logger, cfg *config.Config) *router.Dependencies {
	userRepo := postgres.NewUserRepository(db)
	jwtService := auth.NewJWTService(cfg.JWT.AccessTokenSecret, cfg.JWT.RefreshTokenSecret, cfg.JWT.AccessTokenExpiration, cfg.JWT.RefreshTokenExpiration)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService, logger)
	userUsecase := usecase.NewUserUsecase(userRepo, logger)

	authHanlder := handler.NewAuthHandler(authUsecase, logger, cfg)
	userHandler := handler.NewUserHandler(userUsecase, logger)

	return &router.Dependencies{Log: logger, Cfg: cfg, Db: db, JWTService: jwtService, AuthHandler: authHanlder, UserHandler: userHandler}
}

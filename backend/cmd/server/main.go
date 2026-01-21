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
	"github.com/prabalesh/loco/backend/internal/delivery/router"
	"github.com/prabalesh/loco/backend/internal/di"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"github.com/prabalesh/loco/backend/pkg/logger"
	"github.com/prabalesh/loco/backend/pkg/seeder"
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
		sqlDB, err := db.DB.DB()
		if err != nil {
			log.Error("Failed to get sql db for closing", zap.Error(err))
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Error("Failed to get sql db for stats", zap.Error(err))
	} else {
		stats := sqlDB.Stats()
		log.Info("Database connection pool initialized",
			zap.Int("open_connections", stats.OpenConnections),
			zap.Int("in_use", stats.InUse),
			zap.Int("idle", stats.Idle),
		)
	}

	// Auto Migrate
	// Drop existing constraints to force GORM to recreate them with ON DELETE CASCADE/SET NULL
	db.DB.Exec("ALTER TABLE submissions DROP CONSTRAINT IF EXISTS fk_submissions_user")
	db.DB.Exec("ALTER TABLE submissions DROP CONSTRAINT IF EXISTS fk_submissions_admin")
	db.DB.Exec("ALTER TABLE problems DROP CONSTRAINT IF EXISTS fk_problems_creator")

	if err := db.DB.AutoMigrate(
		&domain.User{},
		&domain.Problem{},
		&domain.Language{},
		&domain.TestCase{},
		&domain.Submission{},
		&domain.ProblemLanguage{},
		&domain.UserProblemStats{},
		&domain.Achievement{},
		&domain.UserAchievement{},
		&domain.Tag{},
		&domain.Category{},
		&domain.ProblemBoilerplate{},
		&domain.ProblemReferenceSolution{},
	); err != nil {
		log.Fatal("Failed to run auto migrations", zap.Error(err))
	}

	// Run Seeder
	if err := seeder.SeedAll(db, log); err != nil {
		log.Fatal("Failed to seed database", zap.Error(err))
	}

	container := di.NewContainer(db, cfg, log)
	router := router.SetupRouter(container.Handlers)

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

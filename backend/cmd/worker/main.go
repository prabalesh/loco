package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/prabalesh/loco/backend/internal/infrastructure/piston"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/prabalesh/loco/backend/internal/infrastructure/worker"
	"github.com/prabalesh/loco/backend/internal/repository/postgres"
	"github.com/prabalesh/loco/backend/internal/services/codegen"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"github.com/prabalesh/loco/backend/pkg/logger"
	"github.com/prabalesh/loco/backend/pkg/redis"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	// 1. Initialize Config
	config.InitConfig()
	cfg := config.GetConfig()

	// 2. Initialize Logger
	if err := logger.InitLogger(cfg.Log.Level); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
	loggers := logger.GetLogger()

	loggers.Info("Starting Worker Service...")

	// 3. Initialize Database
	db, err := database.NewPostgresDB(cfg.Database, loggers)
	if err != nil {
		loggers.Fatal("Failed to connect to database", zap.Error(err))
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		loggers.Fatal("Failed to get sqlDB", zap.Error(err))
	}
	defer sqlDB.Close()

	// 4. Initialize Redis
	redisClient, err := redis.NewRedisClient(cfg.Redis, loggers)
	if err != nil {
		loggers.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// 5. Initialize Repositories
	submissionRepo := postgres.NewSubmissionRepository(db)
	problemRepo := postgres.NewProblemRepository(db)
	testCaseRepo := postgres.NewTestCaseRepository(db)
	languageRepo := postgres.NewLanguageRepository(db)
	problemLanguageRepo := postgres.NewProblemLanguageRepository(db)
	userProblemStatsRepo := postgres.NewUserProblemStatsRepository(db)
	achievementRepo := postgres.NewAchievementRepository(db)
	userRepo := postgres.NewUserRepository(db)
	boilerplateRepo := postgres.NewBoilerplateRepository(db)
	typeImplementationRepo := postgres.NewTypeImplementationRepository(db.DB)

	// 6. Initialize Services
	pistonService := piston.NewPistonService(cfg, loggers)
	jobQueue := queue.NewJobQueue(redisClient, loggers)
	codeGenService := codegen.NewCodeGenService(typeImplementationRepo)
	boilerplateService := codegen.NewBoilerplateService(boilerplateRepo, languageRepo, testCaseRepo, codeGenService)

	achievementUsecase := usecase.NewAchievementUsecase(
		achievementRepo,
		userRepo,
		submissionRepo,
		problemRepo,
		redisClient,
		loggers,
	)

	// 7. Initialize Worker
	submissionWorker := worker.NewWorker(
		jobQueue,
		submissionRepo,
		problemRepo,
		testCaseRepo,
		languageRepo,
		problemLanguageRepo,
		pistonService,
		boilerplateService,
		userProblemStatsRepo,
		loggers,
		redisClient.Client,
		cfg,
	)

	achievementWorker := worker.NewAchievementWorker(
		jobQueue,
		achievementUsecase,
		submissionRepo,
		loggers,
	)

	// 8. Start Worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		submissionWorker.Start(ctx)
	}()

	go func() {
		achievementWorker.Start(ctx)
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	loggers.Info("Shutting down worker...")
	submissionWorker.Stop()
	achievementWorker.Stop()
	cancel()
	loggers.Info("Worker stopped")
}

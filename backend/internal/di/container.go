package di

import (
	"go.uber.org/zap"

	"github.com/prabalesh/loco/backend/internal/delivery/cookies"
	"github.com/prabalesh/loco/backend/internal/delivery/handler"
	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
	"github.com/prabalesh/loco/backend/internal/delivery/router"
	"github.com/prabalesh/loco/backend/internal/infrastructure/auth"
	"github.com/prabalesh/loco/backend/internal/infrastructure/email"
	"github.com/prabalesh/loco/backend/internal/infrastructure/piston"
	"github.com/prabalesh/loco/backend/internal/infrastructure/queue"
	"github.com/prabalesh/loco/backend/internal/infrastructure/worker"
	"github.com/prabalesh/loco/backend/internal/repository/postgres"
	"github.com/prabalesh/loco/backend/internal/usecase"
	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/prabalesh/loco/backend/pkg/database"
	"github.com/prabalesh/loco/backend/pkg/redis"
)

type Container struct {
	Handlers *router.Dependencies
	Worker   *worker.Worker
}

func NewContainer(db *database.Database, cfg *config.Config, logger *zap.Logger) *Container {
	// Repositories
	userRepo := postgres.NewUserRepository(db)
	problemRepo := postgres.NewProblemRepository(db)
	languageRepo := postgres.NewLanguageRepository(db)
	testCaseRepo := postgres.NewTestCaseRepository(db)
	submissionRepo := postgres.NewSubmissionRepository(db)
	problemLanguageRepo := postgres.NewProblemLanguageRepository(db)
	userProblemStatsRepo := postgres.NewUserProblemStatsRepository(db)

	// Redis client
	redisClient, err := redis.NewRedisClient(cfg.Redis, logger)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	// Services
	jwtService := auth.NewJWTService(cfg.JWT.AccessTokenSecret, cfg.JWT.RefreshTokenSecret, cfg.JWT.AccessTokenExpiration, cfg.JWT.RefreshTokenExpiration)
	emailService := email.NewEmailService(cfg, logger)
	pistonService := piston.NewPistonService(cfg, logger)
	jobQueue := queue.NewJobQueue(redisClient, logger)

	cookieManager := cookies.NewCookieManager(cfg)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService, emailService, cfg, logger)
	userUsecase := usecase.NewUserUsecase(userRepo, submissionRepo, logger)
	adminUsecase := usecase.NewAdminUsecase(userRepo, submissionRepo, redisClient.Client, logger)
	problemLanguageUsecase := usecase.NewProblemLanguageUsecase(problemLanguageRepo, problemRepo, languageRepo, logger)
	problemUsecase := usecase.NewProblemUsecase(problemRepo, testCaseRepo, userProblemStatsRepo, cfg, logger)
	languageUsecase := usecase.NewLanguageUsecase(languageRepo, cfg, logger)
	testCaseUsecase := usecase.NewTestCaseUsecase(testCaseRepo, problemRepo, cfg, logger)
	submissionUsecase := usecase.NewSubmissionUsecase(submissionRepo, problemRepo, testCaseRepo, languageRepo, problemLanguageRepo, pistonService, jobQueue, cfg, logger)
	queueStatusUsecase := usecase.NewQueueStatusUsecase(submissionRepo, redisClient.Client, logger)

	// Worker
	submissionWorker := worker.NewWorker(jobQueue, submissionRepo, problemRepo, testCaseRepo, languageRepo, problemLanguageRepo, pistonService, userProblemStatsRepo, logger, redisClient.Client, cfg)

	// Handlers
	authHanlder := handler.NewAuthHandler(authUsecase, logger, cfg, cookieManager)
	userHandler := handler.NewUserHandler(userUsecase, logger)
	adminAuthHandler := handler.NewAdminAuthHandler(authUsecase, logger, cfg, cookieManager)
	adminHandler := handler.NewAdminHandler(adminUsecase, logger)
	problemHandler := handler.NewProblemHandler(problemUsecase, problemLanguageUsecase, languageUsecase, submissionUsecase, logger, cfg)
	languageHandler := handler.NewLanguageHandler(languageUsecase, logger, cfg)
	testCaseHandler := handler.NewTestCaseHandler(testCaseUsecase, logger, cfg)
	submissionHandler := handler.NewSubmissionHandler(submissionUsecase, logger)
	queueHandler := handler.NewQueueHandler(queueStatusUsecase, logger)
	leaderboardUsecase := usecase.NewLeaderboardUsecase(userRepo, logger)
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardUsecase, logger)

	// Middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(redisClient.Client, logger, &cfg.RateLimit)
	submissionRateLimitMiddleware := middleware.NewSubmissionRateLimitMiddleware(redisClient.Client, logger, &cfg.SubmissionRateLimit)
	runCodeRateLimitMiddleware := middleware.NewRunCodeRateLimitMiddleware(redisClient.Client, logger, &cfg.RunCodeRateLimit)

	deps := &router.Dependencies{
		Log:                 logger,
		Cfg:                 cfg,
		Db:                  db,
		JWTService:          jwtService,
		AuthHandler:         authHanlder,
		UserHandler:         userHandler,
		AdminHandler:        adminHandler,
		AdminAuthHandler:    adminAuthHandler,
		ProblemHandler:      problemHandler,
		LanguageHandler:     languageHandler,
		TestCaseHandler:     testCaseHandler,
		SubmissionHandler:   submissionHandler,
		QueueHandler:        queueHandler,
		LeaderboardHandler:  leaderboardHandler,
		RateLimit:           rateLimitMiddleware,
		SubmissionRateLimit: submissionRateLimitMiddleware,
		RunCodeRateLimit:    runCodeRateLimitMiddleware,
	}

	return &Container{
		Handlers: deps,
		Worker:   submissionWorker,
	}
}

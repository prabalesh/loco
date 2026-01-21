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

	v2 "github.com/prabalesh/loco/backend/internal/delivery/handler/v2"
)

type Dependencies struct {
	Log         *zap.Logger
	Cfg         *config.Config
	Db          *database.Database
	JWTService  *auth.JWTService
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler

	AdminHandler        *handler.AdminHandler
	AdminAuthHandler    *handler.AdminAuthHandler
	ProblemHandler      *handler.ProblemHandler
	LanguageHandler     *handler.LanguageHandler
	TestCaseHandler     *handler.TestCaseHandler
	SubmissionHandler   *handler.SubmissionHandler
	RateLimit           *middleware.RateLimitMiddleware
	SubmissionRateLimit *middleware.RateLimitMiddleware
	RunCodeRateLimit    *middleware.RateLimitMiddleware

	LeaderboardHandler   *handler.LeaderboardHandler
	AchievementHandler   *handler.AchievementHandler
	NotificationHandler  *handler.NotificationHandler
	CodeGenHandler       *v2.CodeGenHandler
	CodeExecutionHandler *v2.SubmissionHandler
	V2ProblemHandler     *v2.ProblemHandler
	ValidationHandler    *v2.ValidationHandler
	BulkHandler          *v2.BulkHandler
}

func SetupRouter(deps *Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler())

	// ========== AUTH ROUTES ==========
	mux.HandleFunc("POST /auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /auth/verify-email", deps.AuthHandler.VerifyEmail)
	mux.HandleFunc("POST /auth/resend-verification", deps.AuthHandler.ResendVerificationEmail)
	mux.HandleFunc("POST /auth/refresh", deps.AuthHandler.RefreshToken)
	mux.HandleFunc("POST /auth/logout", deps.AuthHandler.Logout)
	mux.HandleFunc("POST /auth/forgot-password", deps.AuthHandler.ForgotPassword)
	mux.HandleFunc("POST /auth/reset-password", deps.AuthHandler.ResetPassword)

	// Protected auth routes
	authMiddleware := middleware.Auth(deps.JWTService, deps.Log)
	mux.Handle("GET /auth/me", authMiddleware(http.HandlerFunc(deps.AuthHandler.GetMe)))

	// ========== USER ROUTES ==========
	mux.Handle("GET /users/me", authMiddleware(http.HandlerFunc(deps.UserHandler.GetProfile)))
	mux.HandleFunc("GET /users/{username}", deps.UserHandler.GetProfileByUsername)

	// ========== PROBLEM ROUTES (PUBLIC) ==========
	mux.HandleFunc("GET /problems", deps.ProblemHandler.ListProblems)
	mux.HandleFunc("GET /tags", deps.ProblemHandler.ListTags)
	mux.HandleFunc("GET /categories", deps.ProblemHandler.ListCategories)
	mux.HandleFunc("GET /problems/{id}", deps.ProblemHandler.GetProblem)
	mux.HandleFunc("GET /problems/{id}/languages", deps.ProblemHandler.ListProblemLanguages)

	// ========== TEST CASE ROUTES (PUBLIC) ==========
	// Public route for getting sample test cases
	mux.HandleFunc("GET /problems/{problem_id}/test-cases/samples", deps.TestCaseHandler.GetSampleTestCases)

	// // ========== LANGUAGE ROUTES (PUBLIC)
	// mux.HandleFunc("GET /languages", deps.LanguageHandler.ListActiveLanguages)
	// mux.HandleFunc("GET /languages/{identifier}", deps.LanguageHandler.GetLanguage)

	// ========== SUBMISSION ROUTES ==========
	submissionRateLimit := deps.SubmissionRateLimit.RateLimit
	runCodeRateLimit := deps.RunCodeRateLimit.RateLimit
	mux.Handle("POST /problems/{problem_id}/run", authMiddleware(runCodeRateLimit(http.HandlerFunc(deps.SubmissionHandler.RunCode))))
	mux.Handle("POST /problems/{problem_id}/submissions", authMiddleware(submissionRateLimit(http.HandlerFunc(deps.SubmissionHandler.Submit))))
	mux.Handle("GET /problems/{problem_id}/submissions", authMiddleware(http.HandlerFunc(deps.SubmissionHandler.ListUserProblemSubmissions)))
	mux.Handle("GET /submissions", authMiddleware(http.HandlerFunc(deps.SubmissionHandler.ListUserSubmissions)))
	mux.Handle("GET /problems/{problem_id}/submissions/{id}", middleware.RegularOrAdminAuth(deps.JWTService, deps.Log)(http.HandlerFunc(deps.SubmissionHandler.GetSubmission)))

	// ========== LEADERBOARD ROUTES ==========
	mux.HandleFunc("GET /leaderboard", deps.LeaderboardHandler.GetLeaderboard)

	// ========== ACHIEVEMENT ROUTES ==========
	mux.HandleFunc("GET /achievements", deps.AchievementHandler.List)
	mux.HandleFunc("GET /users/{username}/achievements", deps.AchievementHandler.GetUserAchievements)
	mux.Handle("GET /users/me/achievements", authMiddleware(http.HandlerFunc(deps.AchievementHandler.GetMyAchievements)))

	// ========== NOTIFICATION ROUTES ==========
	mux.Handle("GET /notifications/stream", authMiddleware(http.HandlerFunc(deps.NotificationHandler.Stream)))

	// ========== ADMIN ROUTES ==========
	adminAuthMiddleware := middleware.RequireAdminAuth(deps.JWTService, deps.Log)

	// Admin auth
	mux.HandleFunc("POST /admin/auth/login", deps.AdminAuthHandler.AdminLogin)
	mux.HandleFunc("POST /admin/auth/logout", deps.AdminAuthHandler.AdminLogout)
	mux.HandleFunc("POST /admin/auth/refresh", deps.AdminAuthHandler.AdminRefreshToken)
	mux.Handle("GET /admin/auth/me", adminAuthMiddleware(http.HandlerFunc(deps.AdminAuthHandler.GetAdminProfile)))

	// Admin user management
	mux.Handle("GET /admin/users", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.ListUsers)))
	mux.Handle("GET /admin/users/{id}", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.GetUser)))
	mux.Handle("DELETE /admin/users/{id}", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.DeleteUser)))
	mux.Handle("PATCH /admin/users/{id}/role", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.UpdateUserRole)))
	mux.Handle("PATCH /admin/users/{id}/status", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.UpdateUserStatus)))
	mux.Handle("GET /admin/analytics", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.GetAnalytics)))

	// ========== ADMIN PROBLEM ROUTES ==========
	mux.Handle("GET /admin/problems", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.ListAllProblems)))
	mux.Handle("POST /admin/problems", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.CreateProblem)))
	mux.Handle("GET /admin/problems/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.GetProblem)))
	mux.Handle("PUT /admin/problems/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.UpdateProblem)))
	mux.Handle("DELETE /admin/problems/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.DeleteProblem)))
	mux.Handle("POST /admin/problems/{id}/publish", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.PublishProblem)))
	mux.Handle("POST /admin/problems/{id}/archive", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.ArchiveProblem)))
	mux.Handle("GET /admin/problems/stats", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.GetProblemStats)))
	mux.Handle("GET /admin/problems/{id}/languages", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.ListProblemLanguages)))
	mux.Handle("POST /admin/problems/{id}/languages", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.CreateProblemLanguage)))
	mux.Handle("PUT /admin/problems/{id}/languages/{language_id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.UpdateProblemLanguage)))
	mux.Handle("DELETE /admin/problems/{id}/languages/{language_id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.DeleteProblemLanguage)))
	mux.Handle("POST /admin/problems/{id}/languages/{language_id}/validate", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.ValidateProblemLanguage)))
	mux.Handle("GET /admin/problems/{id}/languages/{language_id}/preview", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.PreviewProblemLanguage)))

	// ========== ADMIN TEST CASE ROUTES ==========
	mux.Handle("POST /admin/problems/{problem_id}/test-cases", adminAuthMiddleware(http.HandlerFunc(deps.TestCaseHandler.CreateTestCase)))
	mux.Handle("GET /admin/problems/{problem_id}/test-cases", adminAuthMiddleware(http.HandlerFunc(deps.TestCaseHandler.ListTestCases)))
	mux.Handle("GET /admin/problems/{problem_id}/test-cases/count", adminAuthMiddleware(http.HandlerFunc(deps.TestCaseHandler.CountTestCasesByProblem)))
	mux.Handle("DELETE /admin/problems/{problem_id}/test-cases", adminAuthMiddleware(http.HandlerFunc(deps.TestCaseHandler.DeleteAllTestCases)))
	mux.Handle("POST /admin/problems/{problem_id}/test-cases/reorder", adminAuthMiddleware(http.HandlerFunc(deps.TestCaseHandler.ReorderTestCases)))
	mux.Handle("PUT /admin/test-cases/{id}", adminAuthMiddleware(http.HandlerFunc(deps.TestCaseHandler.UpdateTestCase)))
	mux.Handle("DELETE /admin/test-cases/{id}", adminAuthMiddleware(http.HandlerFunc(deps.TestCaseHandler.DeleteTestCase)))
	mux.Handle("POST /admin/problems/{id}/test-cases/validate", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.ValidateTestCases)))

	// ========== ADMIN LANGUAGE ROUTES ==========
	mux.Handle("POST /admin/languages", adminAuthMiddleware(http.HandlerFunc(deps.LanguageHandler.CreateLanguage)))
	mux.Handle("GET /admin/languages", adminAuthMiddleware(http.HandlerFunc(deps.LanguageHandler.ListLanguages)))
	mux.Handle("GET /admin/languages/active", adminAuthMiddleware(http.HandlerFunc(deps.LanguageHandler.ListActiveLanguages)))
	mux.Handle("GET /admin/languages/{id}", adminAuthMiddleware(http.HandlerFunc(deps.LanguageHandler.GetLanguage)))
	mux.Handle("PUT /admin/languages/{id}", adminAuthMiddleware(http.HandlerFunc(deps.LanguageHandler.UpdateLanguage)))
	mux.Handle("DELETE /admin/languages/{id}", adminAuthMiddleware(http.HandlerFunc(deps.LanguageHandler.DeleteLanguage)))
	mux.Handle("POST /admin/languages/{id}/activate", adminAuthMiddleware(http.HandlerFunc(deps.LanguageHandler.ActivateLanguage)))
	mux.Handle("POST /admin/languages/{id}/deactivate", adminAuthMiddleware(http.HandlerFunc(deps.LanguageHandler.DeactivateLanguage)))

	// ========== ADMIN SUBMISSION ROUTES ==========
	mux.Handle("GET /admin/submissions", adminAuthMiddleware(http.HandlerFunc(deps.SubmissionHandler.ListAdminUserSubmissions)))
	mux.Handle("POST /admin/problems/{id}/submit", adminAuthMiddleware(http.HandlerFunc(deps.SubmissionHandler.AdminSubmit)))
	mux.Handle("GET /admin/problems/{id}/submissions", adminAuthMiddleware(http.HandlerFunc(deps.SubmissionHandler.ListProblemSubmissions)))

	// ========== ADMIN TAG ROUTES ==========
	mux.Handle("POST /admin/tags", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.CreateTag)))
	mux.Handle("PUT /admin/tags/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.UpdateTag)))
	mux.Handle("DELETE /admin/tags/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.DeleteTag)))

	// ========== ADMIN CATEGORY ROUTES ==========
	mux.Handle("POST /admin/categories", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.CreateCategory)))
	mux.Handle("PUT /admin/categories/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.UpdateCategory)))
	mux.Handle("DELETE /admin/categories/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.DeleteCategory)))

	// ========== V2 CODEGEN ROUTES ==========
	mux.Handle("POST /api/v2/codegen/stub", adminAuthMiddleware(http.HandlerFunc(deps.CodeGenHandler.GenerateStub)))
	mux.Handle("GET /api/v2/problems/{problem_id}/stub", http.HandlerFunc(deps.CodeGenHandler.GetProblemStub))
	mux.Handle("GET /api/v2/problems/{problem_id}/boilerplates", adminAuthMiddleware(http.HandlerFunc(deps.CodeGenHandler.GetProblemBoilerplates)))
	mux.Handle("POST /api/v2/problems/{problem_id}/submit", authMiddleware(http.HandlerFunc(deps.CodeExecutionHandler.SubmitCode)))

	// V2 Problems
	mux.Handle("GET /api/v2/admin/problems", adminAuthMiddleware(http.HandlerFunc(deps.V2ProblemHandler.AdminListProblems)))
	mux.Handle("POST /api/v2/admin/problems", adminAuthMiddleware(http.HandlerFunc(deps.V2ProblemHandler.CreateProblem)))
	mux.Handle("GET /api/v2/admin/problems/{id}", adminAuthMiddleware(http.HandlerFunc(deps.V2ProblemHandler.AdminGetProblem)))
	mux.Handle("DELETE /api/v2/admin/problems/{id}", adminAuthMiddleware(http.HandlerFunc(deps.V2ProblemHandler.DeleteProblem)))
	mux.Handle("POST /api/v2/admin/problems/{id}/publish", adminAuthMiddleware(http.HandlerFunc(deps.V2ProblemHandler.PublishProblem)))
	mux.Handle("POST /api/v2/admin/problems/{id}/boilerplates", adminAuthMiddleware(http.HandlerFunc(deps.V2ProblemHandler.RegenerateBoilerplates)))
	mux.HandleFunc("GET /api/v2/problems", deps.V2ProblemHandler.ListProblems)
	mux.HandleFunc("GET /api/v2/problems/{slug}", deps.V2ProblemHandler.GetProblem)

	// V2 Validation
	mux.Handle("POST /api/v2/admin/problems/{id}/validate", adminAuthMiddleware(http.HandlerFunc(deps.ValidationHandler.ValidateReferenceSolution)))
	mux.Handle("GET /api/v2/admin/problems/{id}/validation-status", adminAuthMiddleware(http.HandlerFunc(deps.ValidationHandler.GetValidationStatus)))

	// V2 Custom Types
	mux.Handle("GET /api/v2/admin/custom-types", adminAuthMiddleware(http.HandlerFunc(deps.V2ProblemHandler.GetCustomTypes)))

	// V2 Bulk Import
	bulkRateLimiter := middleware.NewRateLimiter(10, 1*time.Hour)
	mux.Handle("POST /api/v2/admin/problems/bulk", adminAuthMiddleware(bulkRateLimiter.Middleware(deps.BulkHandler.BulkImportProblems)))
	mux.Handle("POST /api/v2/admin/problems/bulk-async", adminAuthMiddleware(bulkRateLimiter.Middleware(deps.BulkHandler.BulkImportProblemsAsync)))

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

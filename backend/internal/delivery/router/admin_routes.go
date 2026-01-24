package router

import (
	"net/http"
	"time"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
)

// SetupAdminRoutes configures all admin-authenticated routes
func SetupAdminRoutes(mux *http.ServeMux, deps *Dependencies, adminAuthMiddleware func(http.Handler) http.Handler) {
	// ========== ADMIN AUTH ==========
	mux.HandleFunc("POST /admin/auth/login", deps.AdminAuthHandler.AdminLogin)
	mux.HandleFunc("POST /admin/auth/logout", deps.AdminAuthHandler.AdminLogout)
	mux.HandleFunc("POST /admin/auth/refresh", deps.AdminAuthHandler.AdminRefreshToken)
	mux.Handle("GET /admin/auth/me", adminAuthMiddleware(http.HandlerFunc(deps.AdminAuthHandler.GetAdminProfile)))

	// ========== ADMIN USER MANAGEMENT ==========
	mux.Handle("GET /admin/users", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.ListUsers)))
	mux.Handle("GET /admin/users/{id}", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.GetUser)))
	mux.Handle("DELETE /admin/users/{id}", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.DeleteUser)))
	mux.Handle("PATCH /admin/users/{id}/role", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.UpdateUserRole)))
	mux.Handle("PATCH /admin/users/{id}/status", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.UpdateUserStatus)))
	mux.Handle("GET /admin/analytics", adminAuthMiddleware(http.HandlerFunc(deps.AdminHandler.GetAnalytics)))

	// ========== ADMIN PROBLEM ROUTES ==========
	mux.Handle("GET /admin/problems", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.ListAllProblems)))
	mux.Handle("POST /admin/problems", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.CreateProblem)))
	mux.Handle("GET /admin/problems/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.AdminGetProblem)))
	mux.Handle("PUT /admin/problems/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.UpdateProblem)))
	mux.Handle("DELETE /admin/problems/{id}", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.DeleteProblem)))
	mux.Handle("POST /admin/problems/{id}/publish", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.PublishProblem)))
	mux.Handle("POST /admin/problems/{id}/archive", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.ArchiveProblem)))
	mux.Handle("GET /admin/problems/stats", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.GetProblemStats)))

	// ========== ADMIN PROBLEM LANGUAGE ROUTES ==========
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

	// ========== ADMIN CODEGEN ROUTES ==========
	mux.Handle("POST /codegen/stub", adminAuthMiddleware(http.HandlerFunc(deps.CodeGenHandler.GenerateStub)))

	// ========== ADMIN VALIDATION ROUTES ==========
	mux.Handle("POST /admin/problems/{id}/validate", adminAuthMiddleware(http.HandlerFunc(deps.ValidationHandler.ValidateReferenceSolution)))
	mux.Handle("GET /admin/problems/{id}/validation-status", adminAuthMiddleware(http.HandlerFunc(deps.ValidationHandler.GetValidationStatus)))

	// ========== ADMIN CUSTOM TYPES ROUTES ==========
	mux.Handle("GET /admin/custom-types", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.GetCustomTypes)))
	mux.Handle("POST /admin/problems/{id}/boilerplates", adminAuthMiddleware(http.HandlerFunc(deps.ProblemHandler.RegenerateBoilerplates)))

	// ========== ADMIN BULK IMPORT ROUTES ==========
	bulkRateLimiter := middleware.NewRateLimiter(10, 1*time.Hour)
	mux.Handle("POST /admin/problems/bulk", adminAuthMiddleware(bulkRateLimiter.Middleware(deps.BulkHandler.BulkImportProblems)))
	mux.Handle("POST /admin/problems/bulk-async", adminAuthMiddleware(bulkRateLimiter.Middleware(deps.BulkHandler.BulkImportProblemsAsync)))
}

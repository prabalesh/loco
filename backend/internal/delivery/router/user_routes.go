package router

import (
	"net/http"

	"github.com/prabalesh/loco/backend/internal/delivery/middleware"
)

// SetupUserRoutes configures all user-authenticated routes
func SetupUserRoutes(mux *http.ServeMux, deps *Dependencies, authMiddleware func(http.Handler) http.Handler) {
	// ========== AUTH ROUTES ==========
	// Protected auth endpoints
	mux.Handle("GET /auth/me", authMiddleware(http.HandlerFunc(deps.AuthHandler.GetMe)))
	mux.Handle("POST /auth/refresh", http.HandlerFunc(deps.AuthHandler.RefreshToken))
	mux.Handle("POST /auth/logout", http.HandlerFunc(deps.AuthHandler.Logout))

	// ========== USER ROUTES ==========
	// User profile
	mux.Handle("GET /users/me", authMiddleware(http.HandlerFunc(deps.UserHandler.GetProfile)))

	// ========== SUBMISSION ROUTES ==========
	// Rate limiters
	submissionRateLimit := deps.SubmissionRateLimit.RateLimit
	runCodeRateLimit := deps.RunCodeRateLimit.RateLimit

	// Submit and run code
	mux.Handle("POST /problems/{problem_id}/run", authMiddleware(runCodeRateLimit(http.HandlerFunc(deps.SubmissionHandler.RunCode))))
	mux.Handle("POST /problems/{problem_id}/submissions", authMiddleware(submissionRateLimit(http.HandlerFunc(deps.SubmissionHandler.Submit))))

	// List submissions
	mux.Handle("GET /problems/{problem_id}/submissions", authMiddleware(http.HandlerFunc(deps.SubmissionHandler.ListUserProblemSubmissions)))
	mux.Handle("GET /submissions", authMiddleware(http.HandlerFunc(deps.SubmissionHandler.ListUserSubmissions)))

	// View submission (user or admin)
	regularOrAdminAuth := middleware.RegularOrAdminAuth(deps.JWTService, deps.Log)
	mux.Handle("GET /problems/{problem_id}/submissions/{id}", regularOrAdminAuth(http.HandlerFunc(deps.SubmissionHandler.GetSubmission)))

	// ========== ACHIEVEMENT ROUTES ==========
	// My achievements
	mux.Handle("GET /users/me/achievements", authMiddleware(http.HandlerFunc(deps.AchievementHandler.GetMyAchievements)))

	// ========== NOTIFICATION ROUTES ==========
	// SSE stream for real-time notifications
	mux.Handle("GET /notifications/stream", authMiddleware(http.HandlerFunc(deps.NotificationHandler.Stream)))
}

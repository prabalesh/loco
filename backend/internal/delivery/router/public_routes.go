package router

import (
	"net/http"
)

// SetupPublicRoutes configures all publicly accessible routes (no authentication required)
func SetupPublicRoutes(mux *http.ServeMux, deps *Dependencies) {
	// ========== HEALTH CHECK ==========
	mux.HandleFunc("GET /health", healthHandler())

	// ========== AUTH ROUTES ==========
	// Registration and login
	mux.HandleFunc("POST /auth/register", deps.AuthHandler.Register)
	mux.HandleFunc("POST /auth/login", deps.AuthHandler.Login)
	mux.HandleFunc("POST /auth/verify-email", deps.AuthHandler.VerifyEmail)
	mux.HandleFunc("POST /auth/resend-verification", deps.AuthHandler.ResendVerificationEmail)
	mux.HandleFunc("POST /auth/forgot-password", deps.AuthHandler.ForgotPassword)
	mux.HandleFunc("POST /auth/reset-password", deps.AuthHandler.ResetPassword)

	// ========== USER ROUTES ==========
	// Public user profiles
	mux.HandleFunc("GET /users/{username}", deps.UserHandler.GetProfileByUsername)

	// ========== PROBLEM ROUTES ==========
	// Browse problems
	mux.HandleFunc("GET /problems", deps.ProblemHandler.ListProblems)
	mux.HandleFunc("GET /problems/{id}", deps.ProblemHandler.GetProblem)
	mux.HandleFunc("GET /problems/{id}/boilerplates", deps.ProblemHandler.ListProblemLanguages)

	// Tags and categories
	mux.HandleFunc("GET /tags", deps.ProblemHandler.ListTags)
	mux.HandleFunc("GET /categories", deps.ProblemHandler.ListCategories)

	// ========== TEST CASE ROUTES ==========
	// Sample test cases (public)
	mux.HandleFunc("GET /problems/{problem_id}/test-cases/samples", deps.TestCaseHandler.GetSampleTestCases)

	// ========== LEADERBOARD ROUTES ==========
	mux.HandleFunc("GET /leaderboard", deps.LeaderboardHandler.GetLeaderboard)

	// ========== ACHIEVEMENT ROUTES ==========
	mux.HandleFunc("GET /achievements", deps.AchievementHandler.List)
	mux.HandleFunc("GET /users/{username}/achievements", deps.AchievementHandler.GetUserAchievements)

	// ========== CODEGEN ROUTES ==========
	// Get problem stub (public for IDE integrations)
	mux.HandleFunc("GET /problems/{problem_id}/stub", deps.CodeGenHandler.GetProblemStub)
}

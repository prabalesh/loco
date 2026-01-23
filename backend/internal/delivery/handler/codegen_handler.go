package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/services/codegen"
)

type CodeGenHandler struct {
	codeGenService     *codegen.CodeGenService
	boilerplateService *codegen.BoilerplateService
	languageRepo       domain.LanguageRepository
	problemRepo        domain.ProblemRepository
	testCaseRepo       domain.TestCaseRepository
}

func NewCodeGenHandler(
	problemRepo domain.ProblemRepository,
	languageRepo domain.LanguageRepository,
	testCaseRepo domain.TestCaseRepository,
	boilerplateService *codegen.BoilerplateService,
	codeGenService *codegen.CodeGenService,
) *CodeGenHandler {
	return &CodeGenHandler{
		codeGenService:     codeGenService,
		problemRepo:        problemRepo,
		languageRepo:       languageRepo,
		testCaseRepo:       testCaseRepo,
		boilerplateService: boilerplateService,
	}
}

type GenerateStubRequest struct {
	FunctionName string                   `json:"function_name"`
	ReturnType   domain.GenericType       `json:"return_type"`
	Parameters   []domain.SchemaParameter `json:"parameters"`
	LanguageSlug string                   `json:"language_slug"`
}

type GenerateStubResponse struct {
	StubCode string `json:"stub_code"`
}

// POST /api/v2/codegen/stub - Generate starter code
func (h *CodeGenHandler) GenerateStub(w http.ResponseWriter, r *http.Request) {
	var req GenerateStubRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	signature := domain.ProblemSchema{
		FunctionName: req.FunctionName,
		ReturnType:   req.ReturnType,
		Parameters:   req.Parameters,
	}

	stubCode, err := h.codeGenService.GenerateStubCode(signature, req.LanguageSlug)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, GenerateStubResponse{
		StubCode: stubCode,
	})
}

// GET /api/v2/problems/{problem_id}/stub?language={language}
func (h *CodeGenHandler) GetProblemStub(w http.ResponseWriter, r *http.Request) {
	problemIDStr := r.PathValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid problem ID")
		return
	}

	languageSlug := r.URL.Query().Get("language")
	if languageSlug == "" {
		languageSlug = "python" // default
	}

	// Get language by slug to get ID
	language, err := h.languageRepo.GetBySlug(languageSlug)
	if err != nil {
		RespondError(w, http.StatusNotFound, "Language not found")
		return
	}

	// Try to get cached stub code
	stubCode, err := h.boilerplateService.GetStubCode(problemID, language.ID)
	if err != nil {
		// FALLBACK: If not found in cache, generate it (but this shouldn't happen with pre-generation)
		// For now, we return error as per requirement 5 (Handle missing boilerplates gracefully - wait, requirement 5 says regenerate on-demand)

		problem, err := h.problemRepo.GetByID(problemID)
		if err != nil {
			RespondError(w, http.StatusNotFound, "Problem not found")
			return
		}

		if problem.FunctionName == nil || problem.ReturnType == nil || problem.Parameters == nil {
			RespondError(w, http.StatusBadRequest, "Problem signature not defined")
			return
		}

		var params []domain.SchemaParameter
		if err := json.Unmarshal([]byte(*problem.Parameters), &params); err != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to parse problem parameters")
			return
		}

		signature := domain.ProblemSchema{
			FunctionName: *problem.FunctionName,
			ReturnType:   domain.GenericType(*problem.ReturnType),
			Parameters:   params,
		}

		testCases, err := h.testCaseRepo.GetByProblemID(problemID)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to fetch test cases")
			return
		}

		// Generate on-demand and store it too
		err = h.boilerplateService.GenerateBoilerplateForLanguage(problemID, language.ID, signature, languageSlug, testCases, problem.ValidationType)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to generate boilerplate on-demand")
			return
		}

		stubCode, _ = h.boilerplateService.GetStubCode(problemID, language.ID)
	}

	RespondJSON(w, http.StatusOK, GenerateStubResponse{
		StubCode: stubCode,
	})
}

// GET /api/v2/problems/{problem_id}/boilerplates
func (h *CodeGenHandler) GetProblemBoilerplates(w http.ResponseWriter, r *http.Request) {
	problemIDStr := r.PathValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid problem ID")
		return
	}

	boilerplates, err := h.boilerplateService.GetBoilerplateStats(problemID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, boilerplates)
}

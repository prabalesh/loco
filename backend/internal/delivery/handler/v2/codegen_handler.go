package v2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prabalesh/loco/backend/internal/delivery/handler"
	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/services/codegen"
)

type CodeGenHandler struct {
	codeGenService *codegen.CodeGenService
	problemRepo    domain.ProblemRepository
}

func NewCodeGenHandler(problemRepo domain.ProblemRepository) *CodeGenHandler {
	return &CodeGenHandler{
		codeGenService: codegen.NewCodeGenService(),
		problemRepo:    problemRepo,
	}
}

type GenerateStubRequest struct {
	FunctionName string              `json:"function_name"`
	ReturnType   string              `json:"return_type"`
	Parameters   []codegen.Parameter `json:"parameters"`
	LanguageSlug string              `json:"language_slug"`
}

type GenerateStubResponse struct {
	StubCode string `json:"stub_code"`
}

// POST /api/v2/codegen/stub - Generate starter code
func (h *CodeGenHandler) GenerateStub(w http.ResponseWriter, r *http.Request) {
	var req GenerateStubRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	signature := codegen.ProblemSignature{
		FunctionName: req.FunctionName,
		ReturnType:   req.ReturnType,
		Parameters:   req.Parameters,
	}

	stubCode, err := h.codeGenService.GenerateStubCode(signature, req.LanguageSlug)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.RespondJSON(w, http.StatusOK, GenerateStubResponse{
		StubCode: stubCode,
	})
}

// GET /api/v2/problems/{problem_id}/stub?language={language}
func (h *CodeGenHandler) GetProblemStub(w http.ResponseWriter, r *http.Request) {
	problemIDStr := r.PathValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, "Invalid problem ID")
		return
	}

	languageSlug := r.URL.Query().Get("language")
	if languageSlug == "" {
		languageSlug = "python" // default
	}

	problem, err := h.problemRepo.GetByID(problemID)
	if err != nil {
		handler.RespondError(w, http.StatusNotFound, "Problem not found")
		return
	}

	if problem.FunctionName == nil || problem.ReturnType == nil || problem.Parameters == nil {
		handler.RespondError(w, http.StatusBadRequest, "Problem signature not defined")
		return
	}

	var params []codegen.Parameter
	if err := json.Unmarshal([]byte(*problem.Parameters), &params); err != nil {
		handler.RespondError(w, http.StatusInternalServerError, "Failed to parse problem parameters")
		return
	}

	signature := codegen.ProblemSignature{
		FunctionName: *problem.FunctionName,
		ReturnType:   *problem.ReturnType,
		Parameters:   params,
	}

	stubCode, err := h.codeGenService.GenerateStubCode(signature, languageSlug)
	if err != nil {
		handler.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.RespondJSON(w, http.StatusOK, GenerateStubResponse{
		StubCode: stubCode,
	})
}

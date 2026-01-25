package piston

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type ExecutionRequest struct {
	Language           string   `json:"language"`
	Version            string   `json:"version"`
	Files              []File   `json:"files"`
	Stdin              string   `json:"stdin"`
	Args               []string `json:"args"`
	RunTimeout         int64    `json:"run_timeout,omitempty"`
	CompileTimeout     int64    `json:"compile_timeout,omitempty"`
	MemoryLimit        int64    `json:"memory_limit,omitempty"`
	CompileMemoryLimit int64    `json:"compile_memory_limit,omitempty"`
}

type File struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content"`
}

type ExecutionResponse struct {
	Language string `json:"language"`
	Version  string `json:"version"`
	Run      Stage  `json:"run"`
	Compile  Stage  `json:"compile"`
}

type Stage struct {
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	Code     int     `json:"code"`
	Signal   string  `json:"signal"`
	Output   string  `json:"output"`
	WallTime float64 `json:"wall_time"`
	CpuTime  float64 `json:"cpu_time"`
	Memory   float64 `json:"memory"`
}

type ExecutionResult struct {
	Output   string
	Error    string
	ExitCode int
	Signal   string
	Runtime  int
	Memory   int
}

type PistonService interface {
	Execute(problemID int, submissionID *int, language, version, code, input string) (*ExecutionResult, error)
}

type pistonService struct {
	client        *http.Client
	baseURL       string
	logger        *zap.Logger
	executionRepo domain.PistonExecutionRepository
}

func NewPistonService(cfg *config.Config, executionRepo domain.PistonExecutionRepository, logger *zap.Logger) PistonService {
	return &pistonService{
		client: &http.Client{
			// client timeout must be slightly higher than the RunTimeout
			Timeout: 45 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		baseURL:       "http://localhost:2000/api/v2",
		logger:        logger,
		executionRepo: executionRepo,
	}
}

func (s *pistonService) Execute(problemID int, submissionID *int, language, version, code, input string) (*ExecutionResult, error) {
	// These values now leverage your new Docker environment settings
	const (
		megabyte           = 1024 * 1024
		runMemoryLimit     = 512 * megabyte  // Increased to 512MB
		compileMemoryLimit = 1024 * megabyte // Increased to 1GB for complex C++ templates
		runTimeoutMs       = 15000           // Increased to 15s (100 test cases safe)
		compileTimeoutMs   = 20000           // Increased to 20s
	)

	reqBody := ExecutionRequest{
		Language:           language,
		Version:            version,
		Files:              []File{{Content: code}},
		Stdin:              input,
		MemoryLimit:        runMemoryLimit,
		CompileMemoryLimit: compileMemoryLimit,
		RunTimeout:         runTimeoutMs,
		CompileTimeout:     compileTimeoutMs,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	resp, err := s.client.Post(s.baseURL+"/execute", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("piston post error: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("piston error %d: %s", resp.StatusCode, string(respBody))
	}

	// Log to database
	if s.executionRepo != nil {
		execution := &domain.PistonExecution{
			ProblemID:    problemID,
			SubmissionID: submissionID,
			Language:     language,
			Version:      version,
			Code:         code,
			Stdin:        input,
			Response:     datatypes.JSON(respBody),
		}
		_ = s.executionRepo.Create(execution)
	}

	var pistonResp ExecutionResponse
	if err := json.Unmarshal(respBody, &pistonResp); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	// Compilation check
	if pistonResp.Compile.Code != 0 || (pistonResp.Compile.Signal != "" && pistonResp.Compile.Signal != "none") {
		return &ExecutionResult{
			Output:   pistonResp.Compile.Stdout,
			Error:    pistonResp.Compile.Output,
			ExitCode: pistonResp.Compile.Code,
		}, nil
	}

	// Runtime check (Capture signal if process was aborted)
	errorMsg := pistonResp.Run.Stderr
	if pistonResp.Run.Signal != "" && pistonResp.Run.Signal != "none" {
		signal := strings.ToUpper(pistonResp.Run.Signal)
		status := "RUNTIME_ERROR"
		switch signal {
		case "SIGKILL":
			status = "TLE"
		case "SIGABRT":
			status = "MLE"
		case "SIGSEGV":
			status = "SIGSEGV"
		default:
			status = fmt.Sprintf("SIGNAL_%s", signal)
		}
		errorMsg = fmt.Sprintf("Process Terminated (Signal: %s, Status: %s)\n%s", signal, status, errorMsg)
	}

	return &ExecutionResult{
		Output:   pistonResp.Run.Stdout,
		Error:    errorMsg,
		ExitCode: pistonResp.Run.Code,
		Signal:   pistonResp.Run.Signal,
		Runtime:  int(pistonResp.Run.CpuTime),
		Memory:   int(pistonResp.Run.Memory) / 1024,
	}, nil
}

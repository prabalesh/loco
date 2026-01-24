package piston

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
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
	Execute(language, version, code, input string) (*ExecutionResult, error)
}

type pistonService struct {
	client  *http.Client
	baseURL string
	logger  *zap.Logger
}

func NewPistonService(cfg *config.Config, logger *zap.Logger) PistonService {
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
		baseURL: "http://localhost:2000/api/v2",
		logger:  logger,
	}
}

func (s *pistonService) Execute(language, version, code, input string) (*ExecutionResult, error) {
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

	fmt.Println("*******************************")
	fmt.Println("JSON Body:", input)
	fmt.Println("*******************************")
	resp, err := s.client.Post(s.baseURL+"/execute", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("piston post error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("piston error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var pistonResp ExecutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&pistonResp); err != nil {
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

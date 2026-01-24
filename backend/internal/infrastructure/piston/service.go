package piston

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

type ExecutionRequest struct {
	Language string   `json:"language"`
	Version  string   `json:"version"`
	Files    []File   `json:"files"`
	Stdin    string   `json:"stdin"`
	Args     []string `json:"args"`
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
	WallTime float64 `json:"wall_time"` // in milliseconds (some versions use time)
	CpuTime  float64 `json:"cpu_time"`  // in milliseconds
	Memory   float64 `json:"memory"`    // in bytes
}

type PistonService interface {
	Execute(language, version, code, input string) (*ExecutionResult, error)
}

type ExecutionResult struct {
	Output   string
	Error    string
	ExitCode int
	Runtime  int // In milliseconds
	Memory   int // In kilobytes
}

type pistonService struct {
	client  *http.Client
	baseURL string
	logger  *zap.Logger
}

func NewPistonService(cfg *config.Config, logger *zap.Logger) PistonService {
	// Default to emkc.org if not configured, or use env var
	baseURL := "http://localhost:2000/api/v2"
	// You might want to add PISTON_URL to config later

	return &pistonService{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		baseURL: baseURL,
		logger:  logger,
	}
}

func (s *pistonService) Execute(language, version, code, input string) (*ExecutionResult, error) {
	reqBody := ExecutionRequest{
		Language: language,
		Version:  version,
		Files: []File{
			{Content: code},
		},
		Stdin: input,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// fmt.Println(string(jsonBody))

	resp, err := s.client.Post(s.baseURL+"/execute", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to execute code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		s.logger.Error("Piston API error",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(bodyBytes)),
		)
		return nil, fmt.Errorf("piston api returned status: %d", resp.StatusCode)
	}

	var pistonResp ExecutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&pistonResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// fmt.Println("*****************************")
	// fmt.Println(pistonResp.Run.CpuTime)
	// fmt.Println(pistonResp.Run.Memory)
	// fmt.Println(int(pistonResp.Run.CpuTime))
	// fmt.Println(int(pistonResp.Run.Memory) / 1024)
	// fmt.Println("*****************************")

	// Check for compile error first
	if pistonResp.Compile.Code != 0 {
		return &ExecutionResult{
			Output:   pistonResp.Compile.Stdout,
			Error:    pistonResp.Compile.Stderr + "\n" + pistonResp.Compile.Output,
			ExitCode: pistonResp.Compile.Code,
			Runtime:  int(pistonResp.Compile.CpuTime),
			Memory:   int(pistonResp.Compile.Memory) / 1024,
		}, nil
	}

	return &ExecutionResult{
		Output:   pistonResp.Run.Stdout,
		Error:    pistonResp.Run.Stderr,
		ExitCode: pistonResp.Run.Code,
		Runtime:  int(pistonResp.Run.CpuTime),
		Memory:   int(pistonResp.Run.Memory) / 1024,
	}, nil
}

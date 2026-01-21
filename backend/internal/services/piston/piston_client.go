package piston

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultPistonURL = "https://emkc.org/api/v2/piston"
	MaxRetries       = 3
	RetryDelay       = 1 * time.Second
	RequestTimeout   = 30 * time.Second
)

type PistonClient struct {
	baseURL    string
	httpClient *http.Client
}

type ExecuteRequest struct {
	Language       string   `json:"language"`
	Version        string   `json:"version"`
	Files          []File   `json:"files"`
	Stdin          string   `json:"stdin"`
	Args           []string `json:"args"`
	CompileTimeout int      `json:"compile_timeout"` // milliseconds
	RunTimeout     int      `json:"run_timeout"`     // milliseconds
}

type File struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content"`
}

type ExecuteResponse struct {
	Language string     `json:"language"`
	Version  string     `json:"version"`
	Run      RunResult  `json:"run"`
	Compile  *RunResult `json:"compile,omitempty"` // For compiled languages
}

type RunResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Output string `json:"output"` // Combined stdout + stderr
	Code   int    `json:"code"`   // Exit code
	Signal string `json:"signal,omitempty"`
}

type Runtime struct {
	Language string   `json:"language"`
	Version  string   `json:"version"`
	Aliases  []string `json:"aliases"`
}

func NewPistonClient(baseURL string) *PistonClient {
	if baseURL == "" {
		baseURL = DefaultPistonURL
	}

	return &PistonClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
	}
}

// Execute sends code to Piston for execution
func (c *PistonClient) Execute(req ExecuteRequest) (*ExecuteResponse, error) {
	var lastErr error

	// Retry logic
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		resp, err := c.executeOnce(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on client errors (4xx)
		if isClientError(err) {
			return nil, err
		}

		// Retry on server errors or network issues
		if attempt < MaxRetries {
			time.Sleep(RetryDelay * time.Duration(attempt))
			continue
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", MaxRetries, lastErr)
}

func (c *PistonClient) executeOnce(req ExecuteRequest) (*ExecuteResponse, error) {
	// Validate request
	if req.Language == "" {
		return nil, errors.New("language is required")
	}
	if len(req.Files) == 0 {
		return nil, errors.New("at least one file is required")
	}

	// Set default timeouts
	if req.CompileTimeout == 0 {
		req.CompileTimeout = 10000 // 10 seconds
	}
	if req.RunTimeout == 0 {
		req.RunTimeout = 3000 // 3 seconds
	}

	// Marshal request
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/execute", c.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if httpResp.StatusCode != http.StatusOK {
		return nil, &PistonError{
			StatusCode: httpResp.StatusCode,
			Message:    string(respBody),
		}
	}

	// Parse response
	var executeResp ExecuteResponse
	if err := json.Unmarshal(respBody, &executeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &executeResp, nil
}

// GetRuntimes fetches available runtimes from Piston
func (c *PistonClient) GetRuntimes() ([]Runtime, error) {
	url := fmt.Sprintf("%s/runtimes", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var runtimes []Runtime
	if err := json.NewDecoder(resp.Body).Decode(&runtimes); err != nil {
		return nil, fmt.Errorf("failed to parse runtimes: %w", err)
	}

	return runtimes, nil
}

// PistonError represents an error from Piston API
type PistonError struct {
	StatusCode int
	Message    string
}

func (e *PistonError) Error() string {
	return fmt.Sprintf("piston error %d: %s", e.StatusCode, e.Message)
}

func isClientError(err error) bool {
	var pistonErr *PistonError
	if errors.As(err, &pistonErr) {
		return pistonErr.StatusCode >= 400 && pistonErr.StatusCode < 500
	}
	return false
}

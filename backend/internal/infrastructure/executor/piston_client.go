package executor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type PistonClient struct {
	baseURL string
	client  *http.Client
}

type ExecuteRequest struct {
	Language           string   `json:"language"`
	Version            string   `json:"version"`
	Files              []File   `json:"files"`
	Stdin              string   `json:"stdin"`
	Args               []string `json:"args,omitempty"`
	CompileTimeout     int      `json:"compile_timeout"`
	RunTimeout         int      `json:"run_timeout"`
	CompileMemoryLimit int      `json:"compile_memory_limit"`
	RunMemoryLimit     int      `json:"run_memory_limit"`
}

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ExecuteResponse struct {
	Run      RunResult  `json:"run"`
	Compile  *RunResult `json:"compile,omitempty"`
	Language string     `json:"language"`
	Version  string     `json:"version"`
}

type RunResult struct {
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	Code     int     `json:"code"`
	Signal   *string `json:"signal"`
	Output   string  `json:"output"`
	Memory   int     `json:"memory"`
	Message  string  `json:"message"`
	Status   string  `json:"status"`
	CpuTime  int     `json:"cpu_time"`
	WallTime int     `json:"wall_time"`
}

func NewPistonClient(baseURL string) *PistonClient {
	return &PistonClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *PistonClient) Execute(req *ExecuteRequest) (*ExecuteResponse, error) {
	body, _ := json.Marshal(req)

	resp, err := p.client.Post(
		p.baseURL+"/api/v2/execute",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ExecuteResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}

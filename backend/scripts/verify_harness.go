package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/piston"
	"github.com/prabalesh/loco/backend/internal/services/codegen"
	"github.com/prabalesh/loco/backend/pkg/config"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Script starting...")
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		return
	}
	cfg := &config.Config{
		Server: config.ServerConfig{
			PistonURL: "http://localhost:2000",
		},
	}
	pistonService := piston.NewPistonService(cfg, logger)
	cg := codegen.NewCodeGenService(nil)

	languages := []struct {
		slug    string
		version string
		code    string
	}{
		{
			slug:    "python",
			version: "3.10.0",
			code: `def solution(n):
    if n == 2:
        import time
        time.sleep(3)
    return n * 2`,
		},
		/*
		   		{
		   			slug:    "javascript",
		   			version: "18.15.0",
		   			code: `function solution(n) {
		       if (n === 2) {
		           while(true);
		       }
		       return n * 2;
		   }`,
		   		},
		*/
	}

	for _, lang := range languages {
		fmt.Printf("\n--- Testing %s ---\n", lang.slug)
		sig := domain.ProblemSchema{
			FunctionName: "solution",
			ReturnType:   domain.TypeInteger,
			Parameters: []domain.SchemaParameter{
				{Name: "n", Type: domain.TypeInteger},
			},
		}

		harness, _ := cg.GenerateTestHarness(sig, lang.code, lang.slug, []domain.TestCase{{Input: "[1]", ExpectedOutput: "2"}, {Input: "[2]", ExpectedOutput: "4"}}, "EXACT")
		testCases := []map[string]interface{}{
			{"input": []interface{}{1}, "expected": 2},
			{"input": []interface{}{2}, "expected": 4},
			{"input": []interface{}{3}, "expected": 6},
		}
		testInput, _ := json.Marshal(testCases)

		fmt.Println("Executing Piston...")
		res, err := pistonService.Execute(lang.slug, lang.version, harness, string(testInput))
		if err != nil {
			log.Fatalf("Execution failed: %v", err)
		}
		fmt.Println("Piston execution finished.")

		fmt.Printf("Output: %s\n", res.Output)
		var results []domain.TestCaseResult
		if err := json.Unmarshal([]byte(res.Output), &results); err != nil {
			fmt.Printf("Failed to parse: %v\n", err)
			continue
		}

		for _, tr := range results {
			fmt.Printf("Test %d: Status=%s, Time=%dms, Memory=%dKB\n", tr.TestID, tr.Status, tr.TimeMS, tr.MemoryKB)
		}
	}
}

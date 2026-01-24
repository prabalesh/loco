package domain

// BoilerplateService defines the interface for boilerplate generation operations
type BoilerplateService interface {
	GenerateAllBoilerplatesForProblem(problem *Problem) error
	RegenerateBoilerplatesForProblem(problem *Problem) error
	GenerateBoilerplateForLanguage(problemID, languageID int, signature ProblemSchema, languageSlug string, testCases []TestCase, validationType string) error
	GetStubCode(problemID, languageID int) (string, error)
	GetTestHarnessTemplate(problemID, languageID int) (string, error)
	InjectUserCodeIntoHarness(template, userCode string) string
	GetBoilerplateStats(problemID int) (map[string]interface{}, error)
	GetBoilerplatesByProblemID(problemID int) ([]ProblemBoilerplate, error)
	DeleteBoilerplatesByProblemID(problemID int) error
}

// ProblemSchema represents the structure needed for code generation
// This might duplicate what's in codegen but needed here if we want full decoupling
// For now, let's assume ProblemSchema is already in domain or we need to move it?
// Checking domain/problem.go or similar might be needed.
// Wait, I saw ProblemSchema used in codegen/boilerplate_service.go line 49.
// Let me double check where ProblemSchema is defined.

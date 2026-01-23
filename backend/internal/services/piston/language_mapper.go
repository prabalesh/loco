package piston

import (
	"errors"
	"fmt"
)

// LanguageMapper maps our language slugs to Piston runtime identifiers
type LanguageMapper struct {
	mappings map[string]PistonRuntime
}

type PistonRuntime struct {
	Language string
	Version  string
	FileExt  string
	FileName string
}

func NewLanguageMapper() *LanguageMapper {
	return &LanguageMapper{
		mappings: map[string]PistonRuntime{
			"python": {
				Language: "python",
				Version:  "3.10.0",
				FileExt:  ".py",
				FileName: "solution.py",
			},
			"javascript": {
				Language: "javascript",
				Version:  "18.15.0",
				FileExt:  ".js",
				FileName: "solution.js",
			},
			"java": {
				Language: "java",
				Version:  "15.0.2",
				FileExt:  ".java",
				FileName: "Solution.java",
			},
			"c++": {
				Language: "cpp",
				Version:  "10.2.0",
				FileExt:  ".cpp",
				FileName: "solution.cpp",
			},
			"c": {
				Language: "c",
				Version:  "10.2.0",
				FileExt:  ".c",
				FileName: "solution.c",
			},
			"go": {
				Language: "go",
				Version:  "1.16.2",
				FileExt:  ".go",
				FileName: "solution.go",
			},
			"rust": {
				Language: "rust",
				Version:  "1.68.2",
				FileExt:  ".rs",
				FileName: "solution.rs",
			},
		},
	}
}

// GetPistonRuntime returns Piston runtime info for our language slug
func (m *LanguageMapper) GetPistonRuntime(languageSlug string) (PistonRuntime, error) {
	runtime, ok := m.mappings[languageSlug]
	if !ok {
		return PistonRuntime{}, fmt.Errorf("unsupported language: %s", languageSlug)
	}
	return runtime, nil
}

// IsSupported checks if a language is supported
func (m *LanguageMapper) IsSupported(languageSlug string) bool {
	_, ok := m.mappings[languageSlug]
	return ok
}

// UpdateMapping allows updating version dynamically (if needed)
func (m *LanguageMapper) UpdateMapping(languageSlug, version string) error {
	runtime, ok := m.mappings[languageSlug]
	if !ok {
		return errors.New("language not found")
	}
	runtime.Version = version
	m.mappings[languageSlug] = runtime
	return nil
}

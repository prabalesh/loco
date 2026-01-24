package validator

import (
	"github.com/prabalesh/loco/backend/internal/domain/dto"
)

func ValidateCreateProblemRequest(req *dto.CreateProblemRequest) map[string]string {
	errors := make(map[string]string)

	if req.Title == "" {
		errors["title"] = "Title is required"
	} else if len(req.Title) < 3 {
		errors["title"] = "Title must be at least 3 characters"
	} else if len(req.Title) > 255 {
		errors["title"] = "Title must not exceed 255 characters"
	}

	if req.Description == "" {
		errors["description"] = "Description is required"
	} else if len(req.Description) < 10 {
		errors["description"] = "Description must be at least 10 characters"
	}

	if req.Difficulty == "" {
		errors["difficulty"] = "Difficulty is required"
	} else if req.Difficulty != "easy" && req.Difficulty != "medium" && req.Difficulty != "hard" {
		errors["difficulty"] = "Difficulty must be easy, medium, or hard"
	}

	return errors
}

func ValidateUpdateProblemRequest(req *dto.UpdateProblemRequest) map[string]string {
	errors := make(map[string]string)

	if req.Title != "" && len(req.Title) < 3 {
		errors["title"] = "Title must be at least 3 characters"
	}

	if req.Description != "" && len(req.Description) < 10 {
		errors["description"] = "Description must be at least 10 characters"
	}

	if req.Difficulty != "" && req.Difficulty != "easy" && req.Difficulty != "medium" && req.Difficulty != "hard" {
		errors["difficulty"] = "Difficulty must be easy, medium, or hard"
	}

	return errors
}

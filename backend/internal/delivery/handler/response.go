package handler

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Data T `json:"data"`
}

func RespondJSON[T any](w http.ResponseWriter, status int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response[T]{Data: data})
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{
		"error": message,
	})
}

func RespondValidationError(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  "validation failed",
		"fields": errors,
	})
}

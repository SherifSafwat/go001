package api

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is a helper to write a consistent error JSON response
func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Error{ // Use the generated Error struct
		Code:    int32(statusCode),
		Message: message,
	}); err != nil {
		// If JSON encoding fails, we can't write to the response writer
		// since we've already written the status code
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Ptr is a helper function to get a pointer to a value (useful for optional parameters)
func Ptr[T any](v T) *T {
	return &v
}

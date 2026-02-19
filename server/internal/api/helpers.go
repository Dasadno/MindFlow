// server/internal/api/helpers.go
package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON сериализует v в JSON и пишет в ResponseWriter.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError пишет APIError в JSON.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, &APIError{Code: code, Message: message})
}

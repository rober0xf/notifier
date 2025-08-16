package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/rober0xf/notifier/internal/adapters/http/dto"
)

// HTTP handler helper functions
func IsJSONRequest(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteErrorResponse(w http.ResponseWriter, status int, message, details string) {
	response := dto.ErrorResponse{Message: message}
	if details != "" {
		log.Printf("Error details: %s", details)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

package httphelpers

import (
	"net/http"
	"strings"
)

// HTTP handler helper functions
func IsJSONRequest(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

package render

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/parthoshuvo/authsvc/uc"
)

// JSONRenderer defines a JSON renderer.
type JSONRenderer struct {
	indent bool
}

// NewJSONRenderer creates a JSON renderer.
func NewJSONRenderer(indent bool) *JSONRenderer {
	return &JSONRenderer{indent}
}

// Render renders an object to JSON.
// Sets the Content-Type header to "application/json".
// Adds an ETag header if needed.
// Sets the HTTP status and renders the content to JSON.
func (jr *JSONRenderer) Render(w http.ResponseWriter, v interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if hb, ok := v.(uc.Hashable); ok {
		w.Header().Set("ETag", fmt.Sprintf("%q", hb.Hash()))
	}
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if jr.indent {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(v)
}

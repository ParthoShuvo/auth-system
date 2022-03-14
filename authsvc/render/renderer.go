package render

import (
	"net/http"
)

// Renderer defines an interface for output renderers.
type Renderer interface {
	Render(http.ResponseWriter, interface{}, int) error
}

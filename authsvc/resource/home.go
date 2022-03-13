package resource

import (
	"fmt"
	"net/http"
)

// HomeHandler defines a resource that renders the authsvc home page
func HomeHandler(homePage string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, homePage)
	}
}

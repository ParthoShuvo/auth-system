package resource

import (
	"net/http"
)

// Action defines an area of functionality used for authorization purposes.
type Action string

func (action Action) String() string {
	return string(action)
}

// Protector defines an action protector.
type Protector interface {
	Protect(Action, http.Handler) http.HandlerFunc
}

type DefaultProtector struct{}

func (dp *DefaultProtector) Protect(action Action, inner http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer ServerError(w, r)
		inner.ServeHTTP(w, r)
	})
}

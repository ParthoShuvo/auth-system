package resource

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/parthoshuvo/authsvc/log4u"
)

// AuthSvcError defines errors that are created by authsvc resources.
type AuthSvcError struct {
	Status int
	Msg    string
}

// ServerError maps errors to internal server errors.
func ServerError(w http.ResponseWriter, rec *http.Request) {
	if r := recover(); r != nil {
		sendISError(w, fmt.Sprintf("%v", r))
	}
}

// sendISError sends an StatusInternalServerError to the client.
func sendISError(w http.ResponseWriter, msg string) {
	sendError(w, errors.New(msg))
}

func sendBadRequestError(w http.ResponseWriter, msg string) {
	sendError(w, &AuthSvcError{Status: http.StatusBadRequest, Msg: msg})
}

// NewError creates a authsvc specific error.
func NewError(status int, msg string) *AuthSvcError {
	return &AuthSvcError{Status: status, Msg: msg}
}

func (e *AuthSvcError) Error() string {
	return e.Msg
}

// sendError sends an Error to the client with the defined status if the error
// is a AuthSvcError or else with a status of StatusInternalServerError.
// If the status is StatusInternalServerError or greater then the error will be logged.
func sendError(w http.ResponseWriter, err error) {
	serr := toAuthSvcError(err)
	if serr.Status >= http.StatusInternalServerError {
		log.Errorln(serr.Error())
	}
	http.Error(w, serr.Error(), serr.Status)
}

func toAuthSvcError(err error) *AuthSvcError {
	if terr, ok := err.(*AuthSvcError); ok {
		return terr
	}
	return NewError(http.StatusInternalServerError, err.Error())
}

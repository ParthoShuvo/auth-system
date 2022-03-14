package resource

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/render"
	"github.com/parthoshuvo/authsvc/uc/user"
)

type AuthResource struct {
	hndlr    *user.Handler
	rndr     render.Renderer
	validate *validator.Validate
}

func NewAuthResource(hndlr *user.Handler, rndr render.Renderer, validate *validator.Validate) *AuthResource {
	return &AuthResource{hndlr, rndr, validate}
}

func (ar *AuthResource) UserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ServerError(w, r)
		rw := requestWrapper(r)
		lusr, err := rw.loginUser()
		if err != nil {
			sendISError(w, fmt.Sprintf("error unmarshalling loginuser [%v]", lusr))
			return
		}

		if err := ar.validate.Struct(lusr); err != nil {
			err = ar.toCustomValidatorError(err)
			log.Errorf("validation error: [%s]", err.Error())
			sendBadRequestError(w, err.Error())
			return
		}

		usr, err := ar.hndlr.ReadUserByLogin(lusr.Email.String())
		if err != nil {
			log.Errorf("user fetching error: [%s]", err.Error())
			sendISError(w, "user fetching error")
			return
		}
		if usr == nil {
			err := fmt.Errorf("user: %s doesn't exists", lusr.Email)
			log.Error(err.Error())
			sendError(w, &AuthSvcError{Status: http.StatusNotFound, Msg: err.Error()})
			return
		}
		if !lusr.isAuthenticated(usr.Password) {
			err := errors.New("login failed, credentials mismatch")
			log.Error(err.Error())
			sendError(w, &AuthSvcError{Status: http.StatusUnauthorized, Msg: err.Error()})
			return
		}
		if !usr.Verified {
			err := fmt.Errorf("login failed, %s is not verified", usr.Email)
			log.Error(err.Error())
			sendError(w, &AuthSvcError{Status: http.StatusForbidden, Msg: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (ar *AuthResource) toCustomValidatorError(err error) error {
	fieldErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		log.Errorf("Failed to convert to ValidationErrors")
		return errors.New("validation error occurred")
	}

	msg := make([]string, 0)
	fieldErr := fieldErrors[0]
	switch fieldErr.Tag() {
	case "required":
		msg = append(msg, "a required field")
	case "validPwd":
		msg = append(msg, "must contain alpha numeric characters, any of special charaters: _!@$%")
	case "email":
		msg = append(msg, "must contain valid email address")
	case "min":
		msg = append(msg, fmt.Sprintf("at least %s characters", fieldErr.Param()))
	case "max":
		msg = append(msg, fmt.Sprintf("at most %s characters", fieldErr.Param()))
	}
	return fmt.Errorf("non-complaint %s: %s", strings.ToLower(fieldErr.Field()), strings.Join(msg, ","))
}

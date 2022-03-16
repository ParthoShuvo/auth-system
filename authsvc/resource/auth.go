package resource

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/render"
	"github.com/parthoshuvo/authsvc/uc/token"
	"github.com/parthoshuvo/authsvc/uc/user"
)

type AuthResource struct {
	usrHndlr  *user.Handler
	toknHndlr *token.Handler
	rndr      render.Renderer
	validate  *validator.Validate
}

func NewAuthResource(usrHandlr *user.Handler, toknHandlr *token.Handler, rndr render.Renderer, validate *validator.Validate) *AuthResource {
	return &AuthResource{usrHandlr, toknHandlr, rndr, validate}
}

func (aurs *AuthResource) UserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer ServerError(w, r)
		rw := requestWrapper(r)
		lusr, err := rw.loginUser()
		if err != nil {
			sendISError(w, fmt.Sprintf("error unmarshalling loginuser [%v]", lusr))
			return
		}

		if err := aurs.validate.Struct(lusr); err != nil {
			err = aurs.toCustomValidatorError(err)
			log.Errorf("validation error: [%s]", err.Error())
			sendError(w, NewError(http.StatusBadRequest, err.Error()))
			return
		}

		usr, err := aurs.usrHndlr.ReadUserByLogin(lusr.Email.String())
		if err != nil {
			log.Errorf("user fetching error: [%s]", err.Error())
			sendISError(w, "user fetching error")
			return
		}
		if usr == nil {
			err := fmt.Errorf("user: %s doesn't exists", lusr.Email)
			log.Error(err.Error())
			sendError(w, NewError(http.StatusNotFound, err.Error()))
			return
		}
		if !lusr.isAuthenticated(usr.Password) {
			err := errors.New("login failed, credentials mismatch")
			log.Error(err.Error())
			sendError(w, NewError(http.StatusUnauthorized, err.Error()))
			return
		}
		if !usr.Verified {
			err := fmt.Errorf("login failed, %s is not verified", usr.Email)
			log.Error(err.Error())
			sendError(w, NewError(http.StatusForbidden, err.Error()))
			return
		}

		toknPair, err := aurs.toknHndlr.NewAuthTokenPair(usr)
		if err != nil {
			sendISError(w, fmt.Sprintf("error occurred while creating tokens: [%v]", err))
			return
		}
		if err := aurs.rndr.Render(w, toknPair, http.StatusOK); err != nil {
			sendISError(w, fmt.Sprintf("error marshalling tokens [%v]", err))
		}
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

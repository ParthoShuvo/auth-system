package resource

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/render"
	"github.com/parthoshuvo/authsvc/uc/adm"
	"github.com/parthoshuvo/authsvc/uc/token"
)

type TokenResource struct {
	toknHndlr *token.Handler
	admHndlr  *adm.Handler
	rndr      render.Renderer
	validate  *validator.Validate
}

func NewTokenResource(toknHandlr *token.Handler, admHndlr *adm.Handler, rndr render.Renderer, validate *validator.Validate) *TokenResource {
	return &TokenResource{toknHandlr, admHndlr, rndr, validate}
}

func (trs *TokenResource) AccessTokenVerifier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := requestWrapper(r)
		accessToken, err := rw.bearerAuth()
		if err != nil {
			log.Error(err.Error())
			sendError(w, NewError(http.StatusBadRequest, err.Error()))
			return
		}
		tokenClaims, err := trs.toknHndlr.VerifyAccessToken(accessToken)
		if err != nil {
			log.Errorf("Invalid token: [%s], error: [%v]", accessToken, err)
			sendError(w, NewError(http.StatusUnauthorized, "Access token has expired or is not yet valid."))
			return
		}
		usrDetails, err := trs.admHndlr.UserDetailsByJWTClaims(tokenClaims)
		if err != nil {
			log.Errorf("error [%v] occurred on user details for user: [%s]", err, tokenClaims.Subject())
			sendISError(w, "error reading user details")
		}
		if err := trs.rndr.Render(w, usrDetails, http.StatusOK); err != nil {
			sendISError(w, fmt.Sprintf("error: [%v] marshalling user details for user [%s]", err, tokenClaims.Subject()))
		}
	}
}

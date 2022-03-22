package resource

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/render"
	"github.com/parthoshuvo/authsvc/uc/adm"
	"github.com/parthoshuvo/authsvc/uc/token"
	"github.com/parthoshuvo/authsvc/uc/user"
)

type TokenResource struct {
	toknHndlr *token.Handler
	admHndlr  *adm.Handler
	usrHndlr  *user.Handler
	rndr      render.Renderer
}

func NewTokenResource(toknHandlr *token.Handler, admHndlr *adm.Handler, usrHndlr *user.Handler, rndr render.Renderer) *TokenResource {
	return &TokenResource{toknHandlr, admHndlr, usrHndlr, rndr}
}

func (trs *TokenResource) AccessTokenVerifier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := requestWrapper(r)
		accessToken, err := unmarshallAccessToken(rw)
		if err != nil {
			sendISError(w, fmt.Sprintf("error unmarshalling access token [%v]", err))
			return
		}
		if accessToken == "" {
			err = errors.New("access token is empty")
			log.Errorf(err.Error())
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
			return
		}
		if err := trs.rndr.Render(w, usrDetails, http.StatusOK); err != nil {
			sendISError(w, fmt.Sprintf("error: [%v] marshalling user details for user [%s]", err, tokenClaims.Subject()))
		}
	}
}

func (trs *TokenResource) TokenPairGenerator() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := requestWrapper(r)
		refreshToken, err := unmarshallRefreshToken(rw)
		if err != nil {
			sendISError(w, fmt.Sprintf("error unmarshalling refresh token [%v]", err))
			return
		}
		if refreshToken == "" {
			err = errors.New("refresh token is empty")
			log.Errorf(err.Error())
			sendError(w, NewError(http.StatusBadRequest, err.Error()))
			return
		}

		tokenClaims, err := trs.toknHndlr.VerifyRefreshToken(refreshToken)
		if err != nil {
			log.Errorf("Invalid token: [%s], error: [%v]", refreshToken, err)
			sendError(w, NewError(http.StatusUnauthorized, "Refresh token has expired or is not yet valid."))
			return
		}
		usr, err := trs.usrHndlr.ReadUserByLogin(tokenClaims.Subject())
		if err != nil {
			log.Errorf("error [%v] occurred on reading user: [%s]", err, tokenClaims.Subject())
			sendISError(w, "error reading user")
			return
		}
		if usr == nil {
			sendError(w, NewError(http.StatusNotFound, "user not found"))
			return
		}

		if err := trs.toknHndlr.RevokeRefreshToken(refreshToken); err != nil {
			log.Errorf("failed to revoke refresh token: [%v]", err)
			sendISError(w, "failed to revoke refresh token")
			return
		}
		toknPair, err := trs.toknHndlr.NewAuthTokenPair(usr)
		if err != nil {
			sendISError(w, fmt.Sprintf("error occurred while creating tokens: [%v]", err))
			return
		}
		if err := trs.rndr.Render(w, toknPair, http.StatusOK); err != nil {
			sendISError(w, fmt.Sprintf("error marshalling tokens [%v]", err))
		}
	}
}

func unmarshallRefreshToken(rw *wrapper) (string, error) {
	data, err := rw.body()
	if err != nil {
		return "", err
	}
	v := struct {
		RefreshToken string `json:"refresh_token"`
	}{}
	err = unmarshall(data, &v)
	if err != nil {
		return "", err
	}
	return v.RefreshToken, nil
}

func unmarshallAccessToken(rw *wrapper) (string, error) {
	data, err := rw.body()
	if err != nil {
		return "", err
	}
	v := struct {
		RefreshToken string `json:"access_token"`
	}{}
	err = unmarshall(data, &v)
	if err != nil {
		return "", err
	}
	return v.RefreshToken, nil
}

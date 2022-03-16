package token

import (
	"github.com/parthoshuvo/authsvc/table/user"
	"github.com/parthoshuvo/authsvc/token"
)

type Handler struct {
	tokenSvc *token.Service
}

func NewHandler(tokenSvc *token.Service) *Handler {
	return &Handler{tokenSvc}
}

func (h *Handler) NewAuthTokenPair(usr *user.User) (*token.AuthTokenPair, error) {
	return h.tokenSvc.NewAuthTokenPair(usr)
}

func (h *Handler) VerifyAccessToken(tokenStr string) (*token.JWTCustomClaims, error) {
	return h.tokenSvc.VerifyAccessToken(tokenStr)
}

func (h *Handler) VerifyRefreshToken(tokenStr string) (*token.JWTCustomClaims, error) {
	return h.tokenSvc.VerifyRefreshToken(tokenStr)
}

func (h *Handler) RevokeRefreshToken(tokenStr string) error {
	return h.tokenSvc.RevokeRefreshToken(tokenStr)
}

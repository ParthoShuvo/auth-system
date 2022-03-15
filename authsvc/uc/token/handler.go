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

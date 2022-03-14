package user

import "github.com/parthoshuvo/authsvc/table/user"

// Handler implements user use-cases.
type Handler struct {
	table *user.Table
}

func NewHandler(t *user.Table) *Handler {
	return &Handler{t}
}

func (h *Handler) ReadUserByLogin(login string) (*user.User, error) {
	return h.table.ReadUserByLogin(login)
}

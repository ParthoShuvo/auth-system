package user

import (
	"github.com/google/uuid"
	"github.com/parthoshuvo/authsvc/table/user"
)

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

func (h *Handler) InsertUser(usr *user.User) (*user.User, error) {
	usr.VerificationCode = uuid.NewString()
	usr, err := h.table.InsertUser(usr)
	if err != nil {
		return usr, err
	}
	usr.Password = ""
	return usr, err
}

func (h *Handler) AssignUserVerification(login string, isVerified bool) error {
	return h.table.AssignUserVerification(login, isVerified)
}

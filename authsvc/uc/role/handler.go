package role

import "github.com/parthoshuvo/authsvc/table/role"

// Handler implements role use-cases.
type Handler struct {
	table *role.Table
}

func NewHandler(t *role.Table) *Handler {
	return &Handler{t}
}

func (hndlr *Handler) ReadUserRoles(login string) ([]*role.Role, error) {
	return hndlr.table.ReadUserRoles(login)
}

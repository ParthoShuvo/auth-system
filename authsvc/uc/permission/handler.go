package permission

import "github.com/parthoshuvo/authsvc/table/permission"

// Handler implements permission use-cases.
type Handler struct {
	table *permission.Table
}

func NewHandler(t *permission.Table) *Handler {
	return &Handler{t}
}

func (hndlr *Handler) ReadUserPermissions(login string) ([]*permission.Permission, error) {
	return hndlr.table.ReadUserPermissions(login)
}

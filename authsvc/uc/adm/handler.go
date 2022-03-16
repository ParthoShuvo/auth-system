package adm

import (
	"reflect"

	permTable "github.com/parthoshuvo/authsvc/table/permission"
	roleTable "github.com/parthoshuvo/authsvc/table/role"
	usrTable "github.com/parthoshuvo/authsvc/table/user"
	"github.com/parthoshuvo/authsvc/token"
	"github.com/parthoshuvo/authsvc/uc/permission"
	"github.com/parthoshuvo/authsvc/uc/role"
	"github.com/parthoshuvo/authsvc/uc/user"
)

type UserDetails struct {
	Firstname   string         `json:"firstname"`
	Lastname    string         `json:"lastname"`
	Email       usrTable.Email `json:"email"`
	Roles       []string       `json:"roles,omitempty"`
	Permissions []string       `json:"permissions,omitempty"`
}

// Handler implements admin use-cases.
type Handler struct {
	usrHndlr  *user.Handler
	roleHndlr *role.Handler
	permHndlr *permission.Handler
}

func NewHandler(usrHndlr *user.Handler, roleHndlr *role.Handler, permHndlr *permission.Handler) *Handler {
	return &Handler{usrHndlr, roleHndlr, permHndlr}
}

func (hndlr *Handler) UserDetailsByJWTClaims(claims *token.JWTCustomClaims) (*UserDetails, error) {
	usr, err := hndlr.usrHndlr.ReadUserByLogin(claims.Subject())
	if err != nil {
		return nil, err
	}
	roles, err := hndlr.roleHndlr.ReadUserRoles(claims.Subject())
	if err != nil {
		return nil, err
	}
	perms, err := hndlr.permHndlr.ReadUserPermissions(claims.Subject())
	if err != nil {
		return nil, err
	}

	return &UserDetails{
		Firstname: usr.Firstname,
		Lastname:  usr.Lastname,
		Email:     usr.Email,
		Roles: hndlr.toStrings(roles, func(v interface{}) string {
			role, _ := v.(*roleTable.Role)
			return role.Name
		}),
		Permissions: hndlr.toStrings(perms, func(v interface{}) string {
			perm, _ := v.(*permTable.Permission)
			return perm.Name
		}),
	}, nil
}

func (hndlr *Handler) toStrings(items interface{}, exec func(interface{}) string) []string {
	v := reflect.ValueOf(items)
	res := make([]string, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		res = append(res, exec(v.Index(i).Interface()))
	}
	return res
}

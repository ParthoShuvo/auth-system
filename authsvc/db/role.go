package db

import (
	"database/sql"

	"github.com/parthoshuvo/authsvc/table/role"
)

// ReadUserRoles fetches all assigned roles for a user.
func (ad *AuthDB) ReadUserRoles(login string) ([]*role.Role, error) {
	return ad.readRoles(func() (*sql.Rows, error) {
		return ad.db.Query("call sp_read_user_role(?)", login)
	})
}

func (ad *AuthDB) readRoles(dbReader func() (*sql.Rows, error)) ([]*role.Role, error) {
	roles := make([]*role.Role, 0, 10)
	rows, err := dbReader()
	if err == sql.ErrNoRows {
		return roles, nil
	}
	if err != nil {
		return roles, err
	}
	defer rows.Close()
	for rows.Next() {
		var role role.Role
		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
		); err != nil {
			return roles, err
		}
		roles = append(roles, &role)
	}
	return roles, nil
}

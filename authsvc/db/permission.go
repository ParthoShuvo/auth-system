package db

import (
	"database/sql"

	"github.com/parthoshuvo/authsvc/table/permission"
)

// ReadUserPermissions fetches all authorized permissions for a user.
func (ad *AuthDB) ReadUserPermissions(login string) ([]*permission.Permission, error) {
	return ad.readPermissions(func() (*sql.Rows, error) {
		return ad.db.Query("call sp_read_user_permission(?)", login)
	})
}

func (ad *AuthDB) readPermissions(dbReader func() (*sql.Rows, error)) ([]*permission.Permission, error) {
	perms := make([]*permission.Permission, 0, 10)
	rows, err := dbReader()
	if err == sql.ErrNoRows {
		return perms, nil
	}
	if err != nil {
		return perms, err
	}
	defer rows.Close()
	for rows.Next() {
		var perm permission.Permission
		if err := rows.Scan(
			&perm.ID,
			&perm.Name,
			&perm.Description,
		); err != nil {
			return perms, err
		}
		perms = append(perms, &perm)
	}
	return perms, nil
}

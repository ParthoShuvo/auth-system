package db

import (
	"database/sql"

	"github.com/parthoshuvo/authsvc/table/user"
)

// ReadUserByLogin reads an user by login.
func (ad *AuthDB) ReadUserByLogin(login string) (*user.User, error) {
	usr := user.User{}
	err := ad.db.QueryRow(
		"call sp_user_get_by_login(?)",
		login).Scan(
		&usr.Firstname,
		&usr.Lastname,
		&usr.Email,
		&usr.Password,
		&usr.RowGUID,
		&usr.Verified,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &usr, err
}

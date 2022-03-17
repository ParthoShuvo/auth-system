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
		&usr.VerificationCode,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &usr, err
}

// InsertUser creates a user.
func (ad *AuthDB) InsertUser(usr *user.User) (*user.User, error) {
	err := ad.db.QueryRow(
		"call sp_insert_user(?,?,?,?,?)",
		usr.Firstname,
		usr.Lastname,
		usr.Email,
		usr.Password,
		usr.VerificationCode).Scan(
		&usr.ID)
	return usr, err
}

// AssignUserVerification assigns verification status to user
func (ad *AuthDB) AssignUserVerification(login string, isVerified bool) error {
	_, err := ad.db.Exec("call sp_user_verification_assignment(?, ?)", login, isVerified)
	return err
}

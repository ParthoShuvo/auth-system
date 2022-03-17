package user

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// User storage structure of table User
type User struct {
	ID               int      `json:"id,omitempty"`
	Firstname        string   `json:"firstname" validate:"required,alphaunicode,max=64"`
	Lastname         string   `json:"lastname" validate:"required,alphaunicode,max=64"`
	Email            Email    `json:"email" validate:"required,email,max=64"`
	Password         Password `json:"password,omitempty" validate:"required,validPwd,min=8,max=64"`
	RowGUID          string   `json:"-"`
	Verified         bool     `json:"-"`
	VerificationCode string   `json:"-"`
}

type Email string

func (e Email) String() string {
	return string(e)
}

func (e Email) IsEmpty() bool {
	return e == ""
}

func (e Email) Equals(other Email) bool {
	return strings.EqualFold(e.String(), other.String())
}

type Password string

func (pw Password) String() string {
	return "•••••"
}

func (pw Password) Equals(other Password) bool {
	return string(pw) == string(other)
}

func (pw Password) Hash() Password {
	data := []byte(pw)
	return Password(fmt.Sprintf("%x", md5.Sum(data)))
}

// Store defines the interface for User storage.
type Store interface {
	ReadUserByLogin(string) (*User, error)
	InsertUser(*User) (*User, error)
	AssignUserVerification(string, bool) error
}

// Table provides implementation of User store
type Table struct {
	store Store
}

func NewTable(s Store) *Table {
	return &Table{s}
}

// ReadUserByLogin fetches an user by login.
func (t *Table) ReadUserByLogin(login string) (*User, error) {
	return t.store.ReadUserByLogin(login)
}

// InsertUser creates a user.
func (t *Table) InsertUser(usr *User) (*User, error) {
	return t.store.InsertUser(usr)
}

// AssignUserVerification assigns verification status to user
func (t *Table) AssignUserVerification(login string, isVerified bool) error {
	return t.store.AssignUserVerification(login, isVerified)
}

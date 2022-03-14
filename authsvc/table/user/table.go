package user

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// User storage structure of table User
type User struct {
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Email     Email    `json:"email"`
	Password  Password `json:"password"`
	RowGUID   string   `json:"-"`
	Verified  bool     `json:"-"`
}

type Email string

func (e Email) String() string {
	return string(e)
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

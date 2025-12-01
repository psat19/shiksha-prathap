package models

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type User struct {
	ID             int
	Email          string
	HashedPassword []byte
	Name           sql.NullString
	Age            sql.NullInt32
	Phone          sql.NullString
}

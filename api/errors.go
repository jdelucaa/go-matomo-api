package api

import (
	"errors"
)

// Exported Errors
var (
	ErrNotFound               = errors.New("entity not found")
	ErrApiUrlCannotBeEmpty    = errors.New("api_url cannot be empty")
	ErrTokenAuthCannotBeEmpty = errors.New("auth_token cannot be empty")
)

type UserError struct {
	User string
	Err  error
}

func (e *UserError) Error() string { return "user: " + e.User }

func (e *UserError) Unwrap() error { return e.Err }

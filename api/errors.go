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

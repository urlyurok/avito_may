package repoerrs

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrUnableToInsert = errors.New("unable to insert")
)

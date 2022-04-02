package postgresql_utilits

import (
	"errors"
	"glide/internal/app"
)

const (
	NoAwards = -1
)

var (
	DefaultErrDB = errors.New("something wrong DB")
	NotFound     = errors.New("user not found")
	Conflict     = errors.New("conflict in db")
)

func NewDBError(externalErr error) *app.GeneralError {
	return &app.GeneralError{
		Err:         DefaultErrDB,
		ExternalErr: externalErr,
	}

}

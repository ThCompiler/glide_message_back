package repository_postgresql

import (
	"errors"
	"github.com/lib/pq"
	postgresql_utilits "glide/internal/pkg/utilits/postgresql"
)

var (
	NicknameAlreadyExist = errors.New("nickname already exist")
	IncorrectCounty      = errors.New("unknown county")
	IncorrectLanguage    = errors.New("unknown language")
)

const (
	codeForeignKeyVal  = "23503"
	countryConstraint  = "users_county_fkey"
	languageConstraint = "user_language_language_fkey"
)

func parsePQError(err *pq.Error) error {
	switch {
	case err.Code == codeForeignKeyVal && err.Constraint == countryConstraint:
		return IncorrectCounty
	case err.Code == codeForeignKeyVal && err.Constraint == languageConstraint:
		return IncorrectLanguage
	default:
		return postgresql_utilits.NewDBError(err)
	}
}

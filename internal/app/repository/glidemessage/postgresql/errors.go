package repository_postgresql

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	postgresql_utilits "glide/internal/pkg/utilits/postgresql"
)

var (
	IncorrectCounty   = errors.New("unknown county")
	IncorrectLanguage = errors.New("unknown language")
)

const (
	codeNullErrorVal   = "22004"
	columnName         = "country"
	codeForeignKeyVal  = "23503"
	countryConstraint  = "glide_message_countries_country_fkey"
	languageConstraint = "glide_message_languages_language_fkey"
)

func parsePQError(err *pq.Error) error {
	switch {
	case err.Code == codeNullErrorVal && err.Column == columnName:
		return IncorrectCounty
	case err.Code == codeForeignKeyVal && err.Column == countryConstraint:
		return IncorrectCounty
	case err.Code == codeForeignKeyVal && err.Column == languageConstraint:
		return IncorrectLanguage
	default:
		return postgresql_utilits.NewDBError(err)
	}
}

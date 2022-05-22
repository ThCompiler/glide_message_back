package repository_postgresql

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"glide/internal/app/repository"
	postgresql_utilits "glide/internal/pkg/utilits/postgresql"
)

var (
	IncorrectCountry  = errors.New("unknown county")
	IncorrectLanguage = errors.New("unknown language")
)

const (
	codeNullErrorVal    = "23502"
	columnName          = "country"
	codeForeignKeyVal   = "23503"
	countryConstraint   = "glide_message_countries_country_fkey"
	languageConstraint  = "glide_message_languages_language_fkey"
	glideMessConstraint = "glide_users_glide_message_fkey"
)

func parsePQError(err *pq.Error) error {
	switch {
	case err.Code == codeNullErrorVal && err.Column == columnName:
		return IncorrectCountry
	case err.Code == codeForeignKeyVal && err.Column == countryConstraint:
		return IncorrectCountry
	case err.Code == codeForeignKeyVal && err.Column == languageConstraint:
		return IncorrectLanguage
	case err.Code == codeForeignKeyVal && err.Column == glideMessConstraint:
		return repository.NotFound
	default:
		return postgresql_utilits.NewDBError(err)
	}
}

package repository_postgresql

import (
	"github.com/jmoiron/sqlx"
	"glide/internal/app/models"
	"glide/internal/app/repository"
)

const (
	queryCountriesGet = `SELECT country_name, picture FROM countries`
	queryLanguagesGet = `SELECT language, picture FROM languages`
)

type InfoRepository struct {
	store *sqlx.DB
}

func NewInfoRepository(st *sqlx.DB) *InfoRepository {
	return &InfoRepository{
		store: st,
	}
}

// GetCountries Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *InfoRepository) GetCountries() ([]models.InfoCountry, error) {
	var res []models.InfoCountry
	if err := repo.store.Select(&res, queryCountriesGet); err != nil {
		return nil, repository.NewDBError(err)
	}
	return res, nil
}

// GetLanguages Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *InfoRepository) GetLanguages() ([]models.InfoLanguage, error) {
	var res []models.InfoLanguage
	if err := repo.store.Select(&res, queryLanguagesGet); err != nil {
		return nil, repository.NewDBError(err)
	}
	return res, nil
}

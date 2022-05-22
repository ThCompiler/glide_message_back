package repository_info

import "glide/internal/app/models"

//go:generate mockgen -destination=mocks/mock_info_repository.go -package=mock_repository -mock_names=Repository=InfoRepository . Repository

type Repository interface {
	// GetCountries Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCountries() ([]models.InfoCountry, error)

	// GetLanguages Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetLanguages() ([]models.InfoLanguage, error)
}

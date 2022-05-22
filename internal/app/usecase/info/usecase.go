package info

import "glide/internal/app/models"

//go:generate mockgen -destination=mocks/mock_statistics_usecase.go -package=mock_usecase -mock_names=Usecase=StatisticsUsecase . Usecase

type Usecase interface {
	// GetCountries Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCountries() ([]models.InfoCountry, error)

	// GetLanguages Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetLanguages() ([]models.InfoLanguage, error)
}

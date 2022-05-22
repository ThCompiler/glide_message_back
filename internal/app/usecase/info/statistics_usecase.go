package info

import (
	"glide/internal/app/models"
	repoInfo "glide/internal/app/repository/info"
)

type InfoUsecase struct {
	repository repoInfo.Repository
}

func NewInfoUsecase(repository repoInfo.Repository) *InfoUsecase {
	return &InfoUsecase{
		repository: repository,
	}
}

// GetCountries Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *InfoUsecase) GetCountries() ([]models.InfoCountry, error) {
	return usecase.repository.GetCountries()
}

// GetLanguages Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *InfoUsecase) GetLanguages() ([]models.InfoLanguage, error) {
	return usecase.repository.GetLanguages()
}

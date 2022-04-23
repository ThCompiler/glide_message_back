package usecase_info

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

// Get Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *InfoUsecase) Get() (*models.Info, error) {
	return usecase.repository.Get()
}

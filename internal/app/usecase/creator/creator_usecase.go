package usecase_creator

import (
	"context"
	"fmt"
	"io"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repoCreator "patreon/internal/app/repository/creator"
	"patreon/internal/microservices/files/delivery/grpc/client"
	repoFiles "patreon/internal/microservices/files/files/repository/files"
	"glide/internal/pkg/utils"

	"github.com/pkg/errors"
)

const NoUser int64 = -2

type CreatorUsecase struct {
	repository     repoCreator.Repository
	repositoryFile client.FileServiceClient
	imageConvector utils.ImageConverter
}

func NewCreatorUsecase(repository repoCreator.Repository, repoClient client.FileServiceClient,
	convector ...utils.ImageConverter) *CreatorUsecase {
	conv := utils.ImageConverter(&utils.ConverterToWebp{})
	if len(convector) != 0 {
		conv = convector[0]
	}
	return &CreatorUsecase{
		repository:     repository,
		repositoryFile: repoClient,
		imageConvector: conv,
	}
}

// Create Errors:
//		CreatorExist
//		models.IncorrectCreatorNickname
//		models.IncorrectCreatorCategory
//		models.IncorrectCreatorDescription
//		repository_postgresql.IncorrectCategory
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *CreatorUsecase) Create(creator *models.Creator) (int64, error) {
	check, err := usecase.repository.ExistsCreator(creator.ID)
	if err != nil && err != repository.NotFound {
		return app.InvalidInt, errors.Wrap(err, fmt.Sprintf("METHOD: usecase_creator.Create; "+
			"ERR: error on get creator with ID = %v", creator.ID))
	}
	if check {
		return app.InvalidInt, CreatorExist
	}

	if err = creator.Validate(); err != nil {
		if errors.Is(err, models.IncorrectCreatorCategory) || errors.Is(err, models.IncorrectCreatorNickname) ||
			errors.Is(err, models.IncorrectCreatorDescription) {
			return -1, err
		}
		return app.InvalidInt, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	return usecase.repository.Create(creator)
}

// GetCreators Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CreatorUsecase) GetCreators() ([]models.Creator, error) {
	return usecase.repository.GetCreators()
}

// SearchCreators Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CreatorUsecase) SearchCreators(pag *models.Pagination,
	searchString string, categories ...string) ([]models.Creator, error) {
	return usecase.repository.SearchCreators(pag, searchString, categories...)
}

// GetCreator Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CreatorUsecase) GetCreator(id int64, userId int64) (*models.CreatorWithAwards, error) {
	cr, err := usecase.repository.GetCreator(id, userId)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("creator with ID = %v not found", id))
	}
	return cr, nil
}

// UpdateCover Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
//			repository_os.ErrorCreate
//   		repository_os.ErrorCopyFile
// 			repository.DefaultErrDB
//			utils.ConvertErr
//  		utils.UnknownExtOfFileName
func (usecase *CreatorUsecase) UpdateCover(data io.Reader, name repoFiles.FileName, id int64) error {
	_, err := usecase.repository.ExistsCreator(id)
	if err != nil {
		return err
	}

	data, name, err = usecase.imageConvector.Convert(context.Background(), data, name)
	if err != nil {
		return errors.Wrap(err, "failed convert to webp of update creator cover")
	}

	path, err := usecase.repositoryFile.SaveFile(context.Background(), data, name, repoFiles.Image)
	if err != nil {
		return err
	}
	err = usecase.repository.UpdateCover(id, app.LoadFileUrl+path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(" err update cover cretor with id %d", id))
	}
	return nil
}

// UpdateAvatar Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
//			repository_os.ErrorCreate
//   		repository_os.ErrorCopyFile
// 			repository.DefaultErrDB
//			utils.ConvertErr
//  		utils.UnknownExtOfFileName
func (usecase *CreatorUsecase) UpdateAvatar(data io.Reader, name repoFiles.FileName, id int64) error {
	_, err := usecase.repository.ExistsCreator(id)
	if err != nil {
		return err
	}

	data, name, err = usecase.imageConvector.Convert(context.Background(), data, name)
	if err != nil {
		return errors.Wrap(err, "failed convert to webp of update creator avatar")
	}

	path, err := usecase.repositoryFile.SaveFile(context.Background(), data, name, repoFiles.Image)
	if err != nil {
		return err
	}

	err = usecase.repository.UpdateAvatar(id, app.LoadFileUrl+path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(" err avatar cover cretor with id %d", id))
	}
	return nil
}

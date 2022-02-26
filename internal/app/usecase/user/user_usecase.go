package usercase_user

import (
	"context"
	"fmt"
	"io"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repoUser "patreon/internal/app/repository/user"
	usePosts "patreon/internal/app/usecase/posts"
	"patreon/internal/microservices/files/delivery/grpc/client"
	repoFiles "patreon/internal/microservices/files/files/repository/files"
	"glide/internal/pkg/utils"

	"github.com/pkg/errors"
)

type UserUsecase struct {
	repository     repoUser.Repository
	repositoryFile client.FileServiceClient
	imageConvector utils.ImageConverter
}

func NewUserUsecase(repository repoUser.Repository, fileClient client.FileServiceClient,
	convector ...utils.ImageConverter) *UserUsecase {
	conv := utils.ImageConverter(&utils.ConverterToWebp{})
	if len(convector) != 0 {
		conv = convector[0]
	}
	return &UserUsecase{
		repository:     repository,
		imageConvector: conv,
		repositoryFile: fileClient,
	}
}

// GetProfile Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *UserUsecase) GetProfile(userID int64) (*models.User, error) {
	u, err := usecase.repository.FindByID(userID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("profile with id %v not found", userID))
	}
	return u, nil
}

// Create Errors:
//		models.EmptyPassword
//		models.IncorrectNickname
// 		models.IncorrectEmailOrPassword
//		repository_postgresql.LoginAlreadyExist
//		repository_postgresql.NicknameAlreadyExist
//		UserExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *UserUsecase) Create(user *models.User) (int64, error) {
	checkUser, err := usecase.repository.FindByLogin(user.Login)
	if err != nil && err != repository.NotFound {
		return -1, errors.Wrap(err, fmt.Sprintf("error on create user with login %v", user.Login))
	}

	if checkUser != nil {
		return -1, UserExist
	}

	if err = user.Validate(); err != nil {
		if errors.Is(err, models.IncorrectEmailOrPassword) || errors.Is(err, models.IncorrectNickname) {
			return -1, err
		}
		return -1, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation user"),
		}
	}

	if err = user.Encrypt(); err != nil {
		if errors.Is(err, models.EmptyPassword) {
			return -1, err
		}

		return -1, app.GeneralError{
			Err:         BadEncrypt,
			ExternalErr: err,
		}
	}

	if err = usecase.repository.Create(user); err != nil {
		return -1, err
	}

	return user.ID, nil
}

// Check Errors:
//		models.IncorrectEmailOrPassword
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *UserUsecase) Check(login string, password string) (int64, error) {
	u, err := usecase.repository.FindByLogin(login)
	if err != nil {
		return -1, err
	}

	if !u.ComparePassword(password) {
		return -1, models.IncorrectEmailOrPassword
	}
	return u.ID, nil
}

// UpdatePassword Errors:
// 		repository.NotFound
//		OldPasswordEqualNew
//		IncorrectEmailOrPassword
//		IncorrectNewPassword
//		models.EmptyPassword
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
//			BadEncrypt
//			app.UnknownError
func (usecase *UserUsecase) UpdatePassword(userId int64, oldPassword, newPassword string) error {
	u, err := usecase.GetProfile(userId)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("profile with id %v not found", userId))
	}
	if !u.ComparePassword(oldPassword) {
		return IncorrectEmailOrPassword
	}
	if u.ComparePassword(newPassword) {
		return OldPasswordEqualNew
	}
	u.MakeEmptyPassword()

	u.Password = newPassword
	if err = u.Encrypt(); err != nil {
		if errors.Is(err, models.EmptyPassword) {
			return err
		}
		return app.GeneralError{
			Err:         BadEncrypt,
			ExternalErr: err,
		}
	}
	if err = u.Validate(); err != nil {
		if errors.Is(err, models.IncorrectEmailOrPassword) {
			return IncorrectNewPassword
		}
		return app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation user"),
		}
	}
	err = usecase.repository.UpdatePassword(userId, u.EncryptedPassword)
	return err
}

// UpdateAvatar Errors:
// 		app.GeneralError with Errors
//			app.UnknownError
//			repository.DefaultErrDB
//			repository_os.ErrorCreate
//   		repository_os.ErrorCopyFile
//			utils.ConvertErr
//  		utils.UnknownExtOfFileName
func (usecase *UserUsecase) UpdateAvatar(data io.Reader, name repoFiles.FileName, userId int64) error {
	var err error
	data, name, err = usecase.imageConvector.Convert(context.Background(), data, name)
	if err != nil {
		return errors.Wrap(err, "failed convert to webp of update user avatar")
	}

	path, err := usecase.repositoryFile.SaveFile(context.Background(), data, name, repoFiles.Image)
	if err != nil {
		return err
	}

	if err := usecase.repository.UpdateAvatar(userId, app.LoadFileUrl+path); err != nil {
		return app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of update avatar"),
		}
	}
	return nil
}

// UpdateNickname Errors:
//		InvalidOldNickname
//		repository.NotFound
//		NicknameExists
// 		app.GeneralError with Errors
//			repository.DefaultErrDB
func (usecase *UserUsecase) UpdateNickname(userID int64, oldNickname string, newNickname string) error {
	u, err := usecase.repository.FindByNickname(oldNickname)
	if err != nil {
		return err
	}
	if u.ID != userID {
		return InvalidOldNickname
	}

	_, err = usecase.repository.FindByNickname(newNickname)
	if err == nil {
		return NicknameExists
	}
	if err != repository.NotFound {
		return err
	}

	if err = usecase.repository.UpdateNickname(oldNickname, newNickname); err != nil {
		return err
	}
	return nil
}

// CheckAccessForAward Errors:
// 		app.GeneralError with Errors
//			repository.DefaultErrDBr
func (usecase *UserUsecase) CheckAccessForAward(userID int64, awardsId int64, creatorId int64) (bool, error) {
	if awardsId == repository.NoAwards || userID == creatorId {
		return true, nil
	}

	if userID == usePosts.EmptyUser {
		return false, nil
	}

	return usecase.repository.IsAllowedAward(userID, awardsId)
}

package usercase_user

import (
	"context"
	"fmt"
	"glide/internal/app"
	"glide/internal/app/models"
	repoFiles "glide/internal/app/repository/files"
	repoUser "glide/internal/app/repository/user"
	"glide/internal/pkg/utilits"
	"io"

	"github.com/pkg/errors"
)

type UserUsecase struct {
	repository     repoUser.Repository
	repositoryFile repoFiles.Repository
	imageConvector utilits.ImageConverter
}

func NewUserUsecase(repository repoUser.Repository, fileRepository repoFiles.Repository,
	convector ...utilits.ImageConverter) *UserUsecase {
	conv := utilits.ImageConverter(&utilits.ConverterToWebp{})
	if len(convector) != 0 {
		conv = convector[0]
	}
	return &UserUsecase{
		repository:     repository,
		imageConvector: conv,
		repositoryFile: fileRepository,
	}
}

// GetProfile Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *UserUsecase) GetProfile(nickname string) (*models.User, error) {
	u, err := usecase.repository.FindByNickname(nickname)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("profile with nickname %v not found", nickname))
	}
	return u, nil
}

// Create Errors:
//		models.EmptyPassword
//		models.IncorrectAge
// 		models.IncorrectNicknameOrPassword
//		repository_postgresql.IncorrectCounty
//		repository_postgresql.IncorrectLanguage
//		repository_postgresql.NicknameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
//			app.UnknownError
//			BadEncrypt
func (usecase *UserUsecase) Create(user *models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		originalError := err
		var generalError *app.GeneralError
		if errors.As(err, &generalError) {
			err = errors.Cause(err).(*app.GeneralError).Err
		}

		if errors.Is(err, models.IncorrectNicknameOrPassword) || errors.Is(err, models.IncorrectAge) {
			return nil, originalError
		}

		return nil, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(originalError, "failed process of validation user"),
		}
	}

	if err := user.Encrypt(); err != nil {
		if errors.Is(err, models.EmptyPassword) {
			return nil, err
		}

		return nil, app.GeneralError{
			Err:         BadEncrypt,
			ExternalErr: err,
		}
	}

	usr, err := usecase.repository.Create(user)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// Check Errors:
//		models.IncorrectNicknameOrPassword
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *UserUsecase) Check(nickname string, password string) (string, error) {
	encPassword, err := usecase.repository.GetPassword(nickname)
	if err != nil {
		return "", err
	}

	u := models.User{
		EncryptedPassword: encPassword,
	}
	if !u.ComparePassword(password) {
		return "", models.IncorrectNicknameOrPassword
	}
	return nickname, nil
}

// Update Errors:
//		models.IncorrectAge
//		repository_postgresql.IncorrectCounty
//		repository_postgresql.IncorrectLanguage
//		repository_postgresql.NicknameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
//			app.UnknownError
func (usecase *UserUsecase) Update(user *models.User) (*models.User, error) {
	if err := user.ValidateUpdate(); err != nil {
		if errors.Is(err, models.IncorrectAge) {
			return nil, err
		}
		return nil, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation user"),
		}
	}

	usr, err := usecase.repository.Update(user)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

/*
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
*/

// UpdateAvatar Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
//			FileSystemError
//			repository.DefaultErrDB
//			utilits.ConvertErr
//  		utilits.UnknownExtOfFileName
func (usecase *UserUsecase) UpdateAvatar(data io.Reader, name repoFiles.FileName, nickname string) error {
	var err error
	data, name, err = usecase.imageConvector.Convert(context.Background(), data, name)
	if err != nil {
		return errors.Wrap(err, "failed convert to webp of update user avatar for user: "+nickname)
	}

	path, err := usecase.repositoryFile.SaveFile(data, name, repoFiles.Image)
	if err != nil {
		return app.GeneralError{
			Err:         FileSystemError,
			ExternalErr: errors.Wrap(err, "error with file: "+string(name)+" for user: "+nickname),
		}
	}

	if err = usecase.repository.UpdateAvatar(nickname, app.LoadFileUrl+path); err != nil {
		return err
	}
	return nil
}

/*
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
*/

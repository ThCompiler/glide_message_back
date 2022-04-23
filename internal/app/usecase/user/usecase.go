package usercase_user

import (
	"glide/internal/app/models"
	repoFiles "glide/internal/app/repository/files"
	"io"
)

//go:generate mockgen -destination=mocks/mock_user_usecase.go -package=mock_usecase -mock_names=Usecase=UserUsecase . Usecase

type Usecase interface {
	// GetProfile Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetProfile(nickname string) (*models.User, error)

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
	Create(user *models.User) (*models.User, error)

	// Check Errors:
	//		models.IncorrectNicknameOrPassword
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Check(login string, password string) (string, error)

	/*// UpdatePassword Errors:
	// 		repository.NotFound
	//		OldPasswordEqualNew
	//		IncorrectEmailOrPassword
	//		IncorrectNewPassword
	//		models.EmptyPassword
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	//			BadEncrypt
	//			app.UnknownError
	UpdatePassword(userId int64, oldPassword, newPassword string) error*/

	// Update Errors:
	//		models.IncorrectAge
	//		repository_postgresql.IncorrectCounty
	//		repository_postgresql.IncorrectLanguage
	//		repository_postgresql.NicknameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	//			app.UnknownError
	Update(user *models.User) (*models.User, error)

	// UpdateAvatar Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors
	//			FileSystemError
	//			repository.DefaultErrDB
	//			utilits.ConvertErr
	//  		utilits.UnknownExtOfFileName
	UpdateAvatar(data io.Reader, name repoFiles.FileName, nickname string) error

	/*// UpdateNickname Errors:
	//		InvalidOldNickname
	//		repository.NotFound
	//		NicknameExists
	// 		app.GeneralError with Errors
	//			app.UnknownError
	UpdateNickname(userID string, oldNickname string, newNickname string) error*/
}

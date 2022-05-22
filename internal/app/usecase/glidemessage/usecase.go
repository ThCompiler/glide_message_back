package glidemessage

import (
	"github.com/sirupsen/logrus"
	"glide/internal/app/models"
	repoFiles "glide/internal/app/repository/files"
	"io"
)

//go:generate mockgen -destination=mocks/mock_posts_usecase.go -package=mock_usecase -mock_names=Usecase=PostsUsecase . Usecase

type Usecase interface {
	// Create Errors:
	//		repository_postgresql.IncorrectCounty
	//		repository_postgresql.IncorrectLanguage
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(log *logrus.Entry, message *models.GlideMessage,
		languages []string, counties []string) (*models.GlideMessage, error)

	// GetGotten Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetGotten(user string) ([]models.GlideMessage, error)

	// GetSent Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetSent(user string) ([]models.GlideMessage, error)

	// UpdatePicture Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors
	//			FileSystemError
	// 			repository.DefaultErrDB
	//			utilits.ConvertErr
	//  		utilits.UnknownExtOfFileName
	UpdatePicture(msgId int64, data io.Reader, name repoFiles.FileName) error

	// Check Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Check(id int64) error

	// Get Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Get(id int64) (*models.GlideMessage, error)

	// ChangeUser Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	ChangeUser(log *logrus.Entry, id int64) error

	// CheckAllowUser Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	CheckAllowUser(id int64, user string) error

	// CheckAllowAuthor Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	CheckAllowAuthor(id int64, user string) error
}

package chats

import (
	"github.com/sirupsen/logrus"
	"glide/internal/app/models"
	repoFiles "glide/internal/app/repository/files"
	"io"
)

//go:generate mockgen -destination=mocks/mock_posts_usecase.go -package=mock_usecase -mock_names=Usecase=PostsUsecase . Usecase

type Usecase interface {
	// Create Errors:
	//		repository_postgresql.ChatAlreadyExists
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(user string, with string) (*models.Chat, error)

	// CheckAllow Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	CheckAllow(user string, chatId int64) error

	// GetChats Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetChats(userId string) ([]models.Chat, error)

	// GetMessages Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetMessages(chatId int64, pag *models.Pagination) ([]models.Message, error)

	// MarkMessages Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	MarkMessages(chatId int64, messageIds []int64) error

	// CreateMessage Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	//			FileSystemError
	// 			repository.DefaultErrDB
	//			utilits.ConvertErr
	//  		utilits.UnknownExtOfFileName
	CreateMessage(log *logrus.Entry, text string, chatId int64, data io.Reader, name repoFiles.FileName, user string) (*models.Message, error)
}

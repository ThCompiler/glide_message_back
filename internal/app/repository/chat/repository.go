package repository_chat

import (
	"glide/internal/app/models"
)

type Repository interface {
	// Create Errors:
	//		ChatAlreadyExists
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(user string, with string) (*models.Chat, error)

	// CheckChat Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	CheckChat(chatId int64) error

	// CheckAllow Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	CheckAllow(user string, chatId int64) error

	// GetChat Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetChat(chatId int64, author string) (*models.Chat, error)

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
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	CreateMessage(text string, chatId int64, image string, user string) (*models.Message, error)
}

package chats

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"glide/internal/app"
	"glide/internal/app/models"
	repoChats "glide/internal/app/repository/chat"
	repoFiles "glide/internal/app/repository/files"
	push_client "glide/internal/microservices/push/delivery/client"
	"glide/internal/pkg/utilits"
	"io"

	"github.com/pkg/errors"
)

type ChatsUsecase struct {
	repository      repoChats.Repository
	filesRepository repoFiles.Repository
	imageConvector  utilits.ImageConverter
	pusher          push_client.Pusher
}

func NewChatsUsecase(repository repoChats.Repository, fileRepository repoFiles.Repository,
	pusher push_client.Pusher, convector ...utilits.ImageConverter) *ChatsUsecase {
	conv := utilits.ImageConverter(&utilits.ConverterToWebp{})
	if len(convector) != 0 {
		conv = convector[0]
	}
	return &ChatsUsecase{
		repository:      repository,
		imageConvector:  conv,
		filesRepository: fileRepository,
		pusher:          pusher,
	}
}

// Create Errors:
//		repository_postgresql.ChatAlreadyExists
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *ChatsUsecase) Create(user string, with string) (*models.Chat, error) {
	res, err := usecase.repository.Create(user, with)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CheckAllow Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *ChatsUsecase) CheckAllow(user string, chatId int64) error {
	return usecase.repository.CheckAllow(user, chatId)
}

// GetChats Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *ChatsUsecase) GetChats(userId string) ([]models.Chat, error) {
	return usecase.repository.GetChats(userId)
}

// GetMessages Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *ChatsUsecase) GetMessages(chatId int64, pag *models.Pagination) ([]models.Message, error) {
	if err := usecase.repository.CheckChat(chatId); err != nil {
		return nil, err
	}

	return usecase.repository.GetMessages(chatId, pag)
}

// MarkMessages Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *ChatsUsecase) MarkMessages(chatId int64, messageIds []int64) error {
	if err := usecase.repository.CheckChat(chatId); err != nil {
		return err
	}

	return usecase.repository.MarkMessages(chatId, messageIds)
}

// CreateMessage Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
//			FileSystemError
// 			repository.DefaultErrDB
//			utilits.ConvertErr
//  		utilits.UnknownExtOfFileName
func (usecase *ChatsUsecase) CreateMessage(log *logrus.Entry, text string, chatId int64, data io.Reader, name repoFiles.FileName,
	user string) (*models.Message, error) {
	chat, err := usecase.repository.GetChat(chatId, user)
	if err != nil {
		return nil, err
	}

	var path = ""
	if data != nil {
		data, name, err = usecase.imageConvector.Convert(context.Background(), data, name)
		if err != nil {
			return nil, errors.Wrap(err, "failed convert to webp of update post cover")
		}

		path, err = usecase.filesRepository.SaveFile(data, name, repoFiles.Image)
		if err != nil {
			return nil, app.GeneralError{
				Err: FileSystemError,
				ExternalErr: errors.Wrap(err, "error with file for message: "+
					string(name)+" for chat: "+fmt.Sprintf("%d", chatId)),
			}
		}
	}

	res, err := usecase.repository.CreateMessage(text, chatId, app.LoadFileUrl+path, user)
	if err != nil {
		return nil, err
	}

	if err = usecase.pusher.NewMessage(res.ID, chat.Companion); err != nil {
		log.Errorf("can't send new message to %s with id %d, gotten err %s", chat.Companion, res.ID, err)
	}

	return res, err
}

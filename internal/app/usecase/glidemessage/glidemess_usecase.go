package glidemessage

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"glide/internal/app"
	"glide/internal/app/models"
	repoChats "glide/internal/app/repository/chat"
	repoFiles "glide/internal/app/repository/files"
	repoGlideMess "glide/internal/app/repository/glidemessage"
	push_client "glide/internal/microservices/push/delivery/client"
	"glide/internal/pkg/utilits"
	"io"

	"github.com/pkg/errors"
)

type GlideMessageUsecase struct {
	repository      repoGlideMess.Repository
	repositoryChats repoChats.Repository
	filesRepository repoFiles.Repository
	imageConvector  utilits.ImageConverter
	pusher          push_client.Pusher
}

func NewGlideMessageUsecase(repository repoGlideMess.Repository, fileRepository repoFiles.Repository,
	pusher push_client.Pusher, convector ...utilits.ImageConverter) *GlideMessageUsecase {
	conv := utilits.ImageConverter(&utilits.ConverterToWebp{})
	if len(convector) != 0 {
		conv = convector[0]
	}
	return &GlideMessageUsecase{
		repository:      repository,
		imageConvector:  conv,
		filesRepository: fileRepository,
		pusher:          pusher,
	}
}

// Create Errors:
//		repository_postgresql.IncorrectCounty
//		repository_postgresql.IncorrectLanguage
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) Create(log *logrus.Entry, message *models.GlideMessage,
	languages []string, counties []string) (*models.GlideMessage, error) {
	res, recipient, err := usecase.repository.Create(message, languages, counties)
	if err != nil {
		return nil, err
	}

	if errPush := usecase.pusher.NewGlideMessage(recipient, res.ID); errPush != nil {
		log.Errorf("Try push new glide message, and got err %s", errPush)
	}

	return res, nil
}

// Apply Errors:
//		repository.NotFound
//		repository_postgresql.ChatAlreadyExists
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) Apply(log *logrus.Entry, user string, msgId int64) (*models.Chat, error) {
	msg, err := usecase.repository.Get(msgId)
	if err != nil {
		return nil, err
	}

	res, err := usecase.repositoryChats.Create(user, msg.Author)
	if err != nil {
		return nil, err
	}

	if err = usecase.repository.Delete(msgId); err != nil {
		log.Errorf("som strange with delete glide message with id %d by user %s, gotten err %s", msgId, user, err)
	}

	return res, nil
}

// GetGotten Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) GetGotten(user string) ([]models.GlideMessage, error) {
	return usecase.repository.GetGotten(user)
}

// GetSent Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) GetSent(user string) ([]models.GlideMessage, error) {
	return usecase.repository.GetSent(user)
}

// UpdatePicture Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
//			FileSystemError
// 			repository.DefaultErrDB
//			utilits.ConvertErr
//  		utilits.UnknownExtOfFileName
func (usecase *GlideMessageUsecase) UpdatePicture(msgId int64, data io.Reader, name repoFiles.FileName) error {
	if err := usecase.repository.Check(msgId); err != nil {
		return err
	}

	var err error
	data, name, err = usecase.imageConvector.Convert(context.Background(), data, name)
	if err != nil {
		return errors.Wrap(err, "failed convert to webp of update post cover")
	}

	path, err := usecase.filesRepository.SaveFile(data, name, repoFiles.Image)
	if err != nil {
		return app.GeneralError{
			Err: FileSystemError,
			ExternalErr: errors.Wrap(err, "error with file for message: "+
				string(name)+" for message: "+fmt.Sprintf("%d", msgId)),
		}
	}

	return usecase.repository.UpdatePicture(msgId, path)
}

// Check Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) Check(id int64) error {
	return usecase.repository.Check(id)
}

// Get Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) Get(id int64) (*models.GlideMessage, error) {
	return usecase.repository.Get(id)
}

// ChangeUser Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) ChangeUser(log *logrus.Entry, id int64) error {
	recipient, err := usecase.repository.ChangeUser(id)
	if err != nil {
		return err
	}

	if errPush := usecase.pusher.NewGlideMessage(recipient, id); errPush != nil {
		log.Errorf("Try push new glide message, and got err %s", errPush)
	}

	return nil
}

// CheckAllowUser Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) CheckAllowUser(id int64, user string) error {
	return usecase.repository.CheckAllowUser(id, user)
}

// CheckAllowAuthor Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *GlideMessageUsecase) CheckAllowAuthor(id int64, user string) error {
	return usecase.repository.CheckAllowAuthor(id, user)
}

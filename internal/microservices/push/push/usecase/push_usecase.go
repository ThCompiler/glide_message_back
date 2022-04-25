package usecase

import (
	"glide/internal/microservices/push"
	"glide/internal/microservices/push/push"
	"glide/internal/microservices/push/push/repository"
)

type PushUsecase struct {
	repository repository.Repository
}

func NewPushUsecase(repository repository.Repository) *PushUsecase {
	return &PushUsecase{
		repository: repository,
	}
}

// PrepareMessagePush with Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) PrepareMessagePush(info *push.MessageInfo) ([]string, *push_models.MessagePush, error) {
	result := &push_models.MessagePush{
		MessageId: info.MessageId,
	}

	var err error

	result.Companion, result.Text, result.ChatId, err = usecase.repository.GetMessageInfo(info.MessageId)
	if err != nil {
		return nil, nil, err
	}

	result.CompanionAvatar, err = usecase.repository.GetUserAvatar(result.Companion)
	return []string{info.Companion}, result, err
}

// PrepareGlidePush with Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PushUsecase) PrepareGlidePush(info *push.GlideInfo) ([]string, *push_models.GlidePush, error) {
	result := &push_models.GlidePush{
		Id: info.GlideId,
	}
	var err error

	result.Author, result.Message, result.Title, result.Country, err = usecase.repository.GetGlideInfo(info.GlideId)
	if err != nil {
		return nil, nil, err
	}

	result.AuthorAvatar, err = usecase.repository.GetUserAvatar(result.Author)
	if err != nil {
		return nil, nil, err
	}

	return []string{info.Companion}, result, err
}

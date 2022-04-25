package usecase

import (
	"glide/internal/microservices/push"
	"glide/internal/microservices/push/push"
)

type Usecase interface {
	// PrepareMessagePush with Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	PrepareMessagePush(info *push.MessageInfo) ([]string, *push_models.MessagePush, error)

	// PrepareGlidePush with Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	PrepareGlidePush(info *push.GlideInfo) ([]string, *push_models.GlidePush, error)
}

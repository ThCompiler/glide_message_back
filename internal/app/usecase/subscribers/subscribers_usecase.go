package usecase_subscribers

import (
	"glide/internal/app/models"
	repository_awards "glide/internal/app/repository/awards"
	repository_subscribers "glide/internal/app/repository/subscribers"

	"github.com/pkg/errors"
)

type SubscribersUsecase struct {
	repoSubscr repository_subscribers.Repository
	repoAwards repository_awards.Repository
}

func NewSubscribersUsecase(repoSubscr repository_subscribers.Repository,
	repoAwards repository_awards.Repository) *SubscribersUsecase {
	return &SubscribersUsecase{
		repoSubscr: repoSubscr,
		repoAwards: repoAwards,
	}
}

// Subscribe Errors:
//		SubscriptionAlreadyExists
//		repository_postgresql.AwardNameNotFound
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (uc *SubscribersUsecase) Subscribe(subscriber *models.Subscriber, token string) error {
	exist, err := uc.repoSubscr.Get(subscriber)
	if err != nil {
		return errors.Wrapf(err, "METHOD: subscribers_usecase.Subscribe; "+
			"ERR: error on checkExists userID = %v creatorID = %v", subscriber.UserID, subscriber.CreatorID)
	}
	if exist {
		return SubscriptionAlreadyExists
	}

	return uc.repoSubscr.Create(subscriber, token)
}

// GetCreators Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (uc *SubscribersUsecase) GetCreators(userID int64) ([]models.CreatorSubscribe, error) {
	return uc.repoSubscr.GetCreators(userID)
}

// GetSubscribers Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (uc *SubscribersUsecase) GetSubscribers(creatorID int64) ([]models.User, error) {
	return uc.repoSubscr.GetSubscribers(creatorID)
}

// UnSubscribe Errors:
//		SubscriptionsNotFound
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (uc *SubscribersUsecase) UnSubscribe(subscriber *models.Subscriber) error {
	exists, err := uc.repoSubscr.Get(subscriber)
	if err != nil {
		return err
	}
	if !exists {
		return SubscriptionsNotFound
	}
	if err = uc.repoSubscr.Delete(subscriber); err != nil {
		return err
	}
	return nil
}

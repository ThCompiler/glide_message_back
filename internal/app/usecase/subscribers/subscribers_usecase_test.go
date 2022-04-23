package usecase_subscribers

import (
	"glide/internal/app"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	"glide/internal/app/usecase"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteSubscribersUsecase struct {
	usecase.SuiteUsecase
	uc Usecase
}

func (s *SuiteSubscribersUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewSubscribersUsecase(s.MockSubscribersRepository, s.MockAwardsRepository)
}

func (s *SuiteSubscribersUsecase) TestSubscribersUsecaseSubscribe_OK() {
	token := "25"
	subscriber := models.TestSubscriber()
	s.MockSubscribersRepository.EXPECT().
		Get(subscriber).
		Return(false, nil).
		Times(1)
	s.MockSubscribersRepository.EXPECT().
		Create(subscriber, token).Times(1).
		Return(nil)
	err := s.uc.Subscribe(subscriber, token)
	assert.NoError(s.T(), err)
}

func (s *SuiteSubscribersUsecase) TestSubscribersUsecaseSubscribe_AlreadyExists() {
	token := "25"
	subscriber := models.TestSubscriber()
	s.MockSubscribersRepository.EXPECT().
		Get(subscriber).
		Return(true, nil).
		Times(1)

	err := s.uc.Subscribe(subscriber, token)
	assert.Equal(s.T(), err, SubscriptionAlreadyExists)
}

func (s *SuiteSubscribersUsecase) TestSubscribersUsecaseSubscribe_CheckExistsError() {
	token := "25"
	subscriber := models.TestSubscriber()
	s.MockSubscribersRepository.EXPECT().
		Get(subscriber).
		Return(false, repository.NewDBError(repository.DefaultErrDB)).
		Times(1)
	err := s.uc.Subscribe(subscriber, token)
	assert.Equal(s.T(), repository.DefaultErrDB, errors.Cause(err).(*app.GeneralError).Err)
}

func (s *SuiteSubscribersUsecase) TestSubscribersUsecaseSubscribe_RepositoryCreateError() {
	token := "25"
	subscriber := models.TestSubscriber()
	s.MockSubscribersRepository.EXPECT().
		Get(subscriber).
		Return(false, nil).
		Times(1)
	s.MockSubscribersRepository.EXPECT().
		Create(subscriber, token).
		Times(1).
		Return(&app.GeneralError{
			Err: repository.DefaultErrDB,
		})
	err := s.uc.Subscribe(subscriber, token)
	assert.Equal(s.T(), repository.DefaultErrDB, errors.Cause(err).(*app.GeneralError).Err)
}

func (s *SuiteSubscribersUsecase) TestSubscriberUsecaseGetCreators_OK() {
	subscriber := models.TestSubscriber()
	creatorSubsc := models.TestCreatorSubscribe()
	expCreators := []models.CreatorSubscribe{*creatorSubsc}
	s.MockSubscribersRepository.EXPECT().
		GetCreators(subscriber.UserID).
		Times(1).
		Return(expCreators, nil)
	res, err := s.uc.GetCreators(subscriber.UserID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expCreators, res)
}

func (s *SuiteSubscribersUsecase) TestSubscriberUsecaseGetCreators_RepositoryError() {
	subscriber := models.TestSubscriber()
	creator := models.TestCreatorSubscribe()
	expCreators := []models.CreatorSubscribe{*creator, *creator}
	s.MockSubscribersRepository.EXPECT().
		GetCreators(subscriber.UserID).
		Times(1).
		Return(expCreators, repository.NewDBError(repository.DefaultErrDB))
	res, err := s.uc.GetCreators(subscriber.UserID)
	assert.Equal(s.T(), expCreators, res)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.DefaultErrDB, errors.Cause(err).(*app.GeneralError).Err)
}

func (s *SuiteSubscribersUsecase) TestSubscriberUsecaseGetSubscribers_OK() {
	subscriber := models.TestSubscriber()
	users := models.TestUsers()
	expUsers := users
	s.MockSubscribersRepository.EXPECT().
		GetSubscribers(subscriber.CreatorID).
		Times(1).
		Return(expUsers, nil)
	res, err := s.uc.GetSubscribers(subscriber.CreatorID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expUsers, res)
}

func (s *SuiteSubscribersUsecase) TestSubscriberUsecaseGetSubscribers_RepositoryError() {
	subscriber := models.TestSubscriber()
	users := models.TestUsers()
	expUsers := users
	s.MockSubscribersRepository.EXPECT().
		GetSubscribers(subscriber.CreatorID).
		Times(1).
		Return(expUsers, repository.NewDBError(repository.DefaultErrDB))

	res, err := s.uc.GetSubscribers(subscriber.CreatorID)
	assert.Equal(s.T(), expUsers, res)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.DefaultErrDB, errors.Cause(err).(*app.GeneralError).Err)
}
func TestSubscribersUsecase(t *testing.T) {
	suite.Run(t, new(SuiteSubscribersUsecase))
}

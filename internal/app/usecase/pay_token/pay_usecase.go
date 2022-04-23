package usecase_pay_token

import (
	uuid "github.com/satori/go.uuid"
	"glide/internal/app/models"
	"glide/internal/app/repository/pay_token"
	"strconv"
	"time"
)

const (
	timeExp = time.Hour * 3
)

type PayTokenUsecase struct {
	repository    repository_pay_token.Repository
	accountNumber string
}

func NewPayTokenUsecase(repository repository_pay_token.Repository, accountNumber string) *PayTokenUsecase {
	return &PayTokenUsecase{
		repository:    repository,
		accountNumber: accountNumber,
	}
}

// GetToken with Errors:
//		app.GeneralError with Errors
//			repository_redis.SetError
func (u *PayTokenUsecase) GetToken(userID int64) (models.PayToken, error) {
	payToken := uuid.NewV4().String()
	userIDtoStr := strconv.Itoa(int(userID))
	err := u.repository.Set(payToken, userIDtoStr, int(timeExp.Seconds()))

	if err != nil {
		return models.PayToken{}, err
	}

	return models.PayToken{Token: payToken}, nil
}

//	CheckToken with Errors:
//	repository_redis.NotFound
//	app.GeneralError with Errors
//		repository_redis.InvalidStorageData
func (u *PayTokenUsecase) CheckToken(token models.PayToken) (bool, error) {
	_, err := u.repository.Get(token.Token)
	if err != nil {
		return false, err
	}
	return true, nil
}

//	CheckTokenByUser with Errors:
//		InvalidUserToken
//		repository_redis.NotFound
//		app.GeneralError with Errors
//			repository_redis.InvalidStorageData
func (u *PayTokenUsecase) CheckTokenByUser(token models.PayToken, userID int64) error {
	userTokenID, err := u.repository.Get(token.Token)
	if err != nil {
		return err
	}
	userTokenIDToInt, err := strconv.Atoi(userTokenID)
	if err != nil {
		return err
	}
	if int64(userTokenIDToInt) != userID {
		return InvalidUserToken
	}

	return nil
}

func (u *PayTokenUsecase) GetAccount() string {
	return u.accountNumber
}

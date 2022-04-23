package payments

import (
	"glide/internal/app/models"
	db_models "glide/internal/app/models"
)

const (
	BaseLimit = 10
)

//go:generate mockgen -destination=mocks/mock_payments_usecase.go -package=mock_usecase -mock_names=Usecase=PaymentsUsecase . Usecase

type Usecase interface {
	// GetUserPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetUserPayments(userID int64, pag *db_models.Pagination) ([]models.UserPayments, error)
	// GetCreatorPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetCreatorPayments(creatorID int64, pag *db_models.Pagination) ([]models.CreatorPayments, error)
	// UpdateStatus Errors:
	//		repository_payments.NotEqualPaymentAmount
	//		repository_payments.CountPaymentsByTokenError
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	UpdateStatus(token string, receiveAmount float64) error
}

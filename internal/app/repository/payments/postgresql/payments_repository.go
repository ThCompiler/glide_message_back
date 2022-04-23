package repository_postgresql

import (
	"fmt"
	"glide/internal/app/models"
	db_models "glide/internal/app/models"
	"glide/internal/app/repository"
	repository_payments "glide/internal/app/repository/payments"
	putilits "glide/internal/app/utilits/postgresql"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

const (
	querySelectUserPayments = "SELECT p.amount, p.date, p.creator_id, u.nickname, cp.category, cp.description, p.status FROM payments p " +
		"JOIN creator_profile cp on p.creator_id = cp.creator_id " +
		"JOIN users u on cp.creator_id = u.users_id where p.users_id = $1 " +
		"ORDER BY p.date DESC "

	querySelectCreatorPayments = "SELECT p.amount, p.date, p.users_id, u.nickname, p.status FROM payments p " +
		"JOIN users u on p.users_id = u.users_id where p.creator_id = $1 and p.status = true " +
		"ORDER BY p.date DESC "
	queryUpdateStatus    = "UPDATE payments SET status = true WHERE pay_token = $1 RETURNING users_id, creator_id, awards_id;"
	queryCountPayments   = "SELECT count(*) from payments where pay_token = $1;"
	queryGetPayment      = "SELECT amount, date, creator_id, users_id, status from payments where pay_token = $1;"
	queryUpdateSubscribe = "UPDATE subscribers set status = true where users_id = $1 and creator_id = $2 and awards_id = $3;"
)

type PaymentsRepository struct {
	store *sqlx.DB
}

func NewPaymentsRepository(store *sqlx.DB) *PaymentsRepository {
	return &PaymentsRepository{
		store: store,
	}
}

// GetUserPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) GetUserPayments(userID int64, pag *db_models.Pagination) ([]models.UserPayments, error) {
	query := querySelectUserPayments

	limit, offset, err := putilits.AddPagination("payments", pag, repo.store)
	query = query + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	if err != nil {
		return nil, err
	}
	if limit == 0 {
		return nil, repository.NotFound
	}

	rows, err := repo.store.Query(query, userID)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	paymentsRes := make([]models.UserPayments, 0, limit)

	for rows.Next() {
		cur := models.UserPayments{}
		if err = rows.Scan(&cur.Amount, &cur.Date, &cur.CreatorID,
			&cur.CreatorNickname, &cur.CreatorCategory, &cur.CreatorDescription, &cur.Status); err != nil {

			_ = rows.Close()
			return nil, repository.NewDBError(errors.Wrapf(err, "method - GetUserPayments"+
				"invalid data in db: table payments"))
		}
		paymentsRes = append(paymentsRes, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return paymentsRes, nil
}

// GetCreatorPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) GetCreatorPayments(creatorID int64, pag *db_models.Pagination) ([]models.CreatorPayments, error) {
	query := querySelectCreatorPayments

	limit, offset, err := putilits.AddPagination("payments", pag, repo.store)
	query = query + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	if err != nil {
		return nil, err
	}
	if limit == 0 {
		return nil, repository.NotFound
	}

	rows, err := repo.store.Query(query, creatorID)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	paymentsRes := make([]models.CreatorPayments, 0, limit)

	for rows.Next() {
		cur := models.CreatorPayments{}
		if err = rows.Scan(&cur.Amount, &cur.Date, &cur.UserID, &cur.UserNickname, &cur.Status); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(errors.Wrapf(err, "method - GetUserPayments"+
				"invalid data in db: table payments"))
		}
		paymentsRes = append(paymentsRes, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return paymentsRes, nil
}

// UpdateStatus Errors:
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) UpdateStatus(token string) error {
	begin, err := repo.store.Begin()
	if err != nil {
		return repository.NewDBError(err)
	}
	awardsID, usersID, creatorID := 0, 0, 0
	err = repo.store.QueryRow(queryUpdateStatus, token).Scan(&usersID, &creatorID, &awardsID)
	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}
	_, err = repo.store.Exec(queryUpdateSubscribe, usersID, creatorID, awardsID)
	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = begin.Commit(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}

// CheckCountPaymentsByToken Errors:
//		repository_payments.CountPaymentsByTokenError
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) CheckCountPaymentsByToken(token string) error {
	count := 0
	err := repo.store.QueryRow(queryCountPayments, token).Scan(&count)
	if err != nil {
		return repository.NewDBError(err)
	}
	if count != 1 {
		return repository_payments.CountPaymentsByTokenError
	}
	return nil
}

// GetPaymentByToken Errors:
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) GetPaymentByToken(token string) (models.Payments, error) {
	res := models.Payments{}
	err := repo.store.QueryRow(queryGetPayment, token).Scan(&res.Amount, &res.Date, &res.CreatorID, &res.UserID,
		&res.Status)
	if err != nil {
		return res, repository.NewDBError(err)
	}
	return res, nil
}

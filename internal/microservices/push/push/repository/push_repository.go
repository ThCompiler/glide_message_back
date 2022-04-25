package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"glide/internal/app"
	repository "glide/internal/pkg/utilits/postgresql"
)

const (
	GetUserAvatarQuery = `SELECT nickname, avatar FROM users WHERE nickname = $1`

	GetMessageInfoQuery = `SELECT u.nickname, cp.avatar FROM creator_profile as cp 
									JOIN users u on cp.creator_id = u.users_id
									WHERE cp.creator_id = $1`

	GetGlideInfoQuery = `SELECT name, price FROM awards where awards_id = $1`
)

type PushRepository struct {
	store *sqlx.DB
}

func NewPushRepository(st *sqlx.DB) *PushRepository {
	return &PushRepository{
		store: st,
	}
}

// GetUserAvatar Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetUserAvatar(username string) (avatar string, err error) {
	if err = repo.store.QueryRow(GetUserAvatarQuery, username).Scan(&avatar); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repository.NotFound
		}
		return "", repository.NewDBError(err)
	}

	return avatar, nil
}

// GetMessageInfo Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetMessageInfo(messageId int64) (author string, text string, chatId int64, err error) {
	if err = repo.store.QueryRow(GetMessageInfoQuery, messageId).Scan(&author, &text, &chatId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", app.InvalidInt, repository.NotFound
		}
		return "", "", app.InvalidInt, repository.NewDBError(err)
	}

	return author, text, chatId, nil
}

// GetGlideInfo Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PushRepository) GetGlideInfo(glideId int64) (author string, message string, title string, country string, err error) {
	if err = repo.store.QueryRow(GetGlideInfoQuery, glideId).Scan(&author, &message, title, country); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", "", "", repository.NotFound
		}
		return "", "", "", "", repository.NewDBError(err)
	}

	return author, message, title, country, nil
}

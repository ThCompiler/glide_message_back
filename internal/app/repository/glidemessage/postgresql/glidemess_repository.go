package repository_postgresql

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	repository_glidemess "glide/internal/app/repository/glidemessage"
	putilits "glide/internal/pkg/utilits/postgresql"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

const (
	createQuery = `
		WITH cnt AS (
		    SELECT country_name as cnt_name FROM countries WHERE lower(country_name) = lower($4)
		)
		INSERT INTO glide_message (title, message, author, country) VALUES ($1, $2, $3, (SELECT cnt_name FROM cnt LIMIT 1)) 
		RETURNING id, title, message, author, created, country`

	createQueryLanguagesStart = `
						WITH lng AS (
    						SELECT language as lng_name FROM languages WHERE lower(language) in (
							`

	createQueryLanguagesEnd = `?)
						)
						INSERT INTO glide_message_countries (language, glide_message) SELECT lng_name, ? FROM lng`

	createQueryCountriesStart = `
						WITH cnt AS (
    						SELECT country_name as cnt_name FROM countries WHERE lower(country_name) in (
							`

	createQueryCountriesEnd = `?)
						)
						INSERT INTO glide_message_countries (country, glide_message) SELECT cnt_name, ? FROM cnt`

	searchUserQuery = `
				WITH usr AS (
				    SELECT usr.nickname FROM users as usr 
				    JOIN user_language ul ON usr.nickname = ul.nickname 
				                        	AND ul.language in (SELECT gml.language FROM glide_message_languages as gml WHERE gml.glide_message = $1)
					WHERE usr.country in (SELECT gmc.country FROM glide_message_countries as gmc WHERE gmc.glide_message = $1) 
					  and usr.nickname not in (SELECT visited_user FROM glide_users WHERE glide_users.glide_message = $1)
				
				)
			`

	checkQuery = `SELECT id FROM chat WHERE id = $1`

	checkAllowQuery = `SELECT id FROM chat WHERE id = $1 and (companion = $2 or author = $2)`

	createMessageQuery = `INSERT INTO messages (message, chat, picture, author) VALUES ($1, $2, $3, $4) 
		RETURNING id, message, picture, author, is_viewed, created`

	getChatsQuery = `
			WITH latest_messages as (
				SELECT min(created) as data FROM messages GROUP BY chat
			)
				SELECT chat.id, chat.author, u.avatar, m.id, m.message, m.picture, m.author, m.is_viewed, m.created FROM chat 
					JOIN messages m on chat.id = m.chat and m.created in(latest_messages)
					JOIN users u on chat.author = u.nickname
				WHERE chat.companion = $1
			UNION 
				SELECT chat.id, chat.companion, u.avatar, m.id, m.message, m.picture, m.author, m.is_viewed, m.created FROM chat
					JOIN messages m on chat.id = m.chat and m.created in(latest_messages)
				    JOIN users u on chat.companion = u.nickname
				WHERE chat.author = $1
	`

	getMessagesQuery = `
				SELECT id, message, picture, author, is_viewed, created FROM messages WHERE chat = $1`

	markMessages = ` UPDATE messages SET is_viewed=true WHERE chat = $1 and id in (?)`
)

type GlideMessageRepository struct {
	store *sqlx.DB
}

var _ = repository_glidemess.Repository(&GlideMessageRepository{})

func NewGlideMessageRepository(st *sqlx.DB) *GlideMessageRepository {
	return &GlideMessageRepository{
		store: st,
	}
}

func (repo *GlideMessageRepository) addToInsert(queryStart string, queryEnd string, arr []string, id int64) (string, []interface{}) {
	var argsString []string
	var args []interface{}
	for _, str := range arr {
		argsString = append(argsString, "lower(?)")

		args = append(args, str)
	}

	query := fmt.Sprintf("%s %s %s", queryStart,
		strings.Join(argsString, ", "), queryEnd)
	query = repo.store.Rebind(query)
	return query, args
}

// Create Errors:
//		IncorrectCounty
//		IncorrectLanguage
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) Create(message *models.GlideMessage, languages []string, counties []string) (*models.GlideMessage, error) {
	tx, err := repo.store.Beginx()
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	if err = tx.QueryRowx(createQuery,
		message.Title,
		message.Message,
		message.Author,
		message.Country).
		Scan(
			&message.ID,
			&message.Title,
			&message.Message,
			&message.Author,
			&message.Created,
			&message.Country); err != nil {
		_ = tx.Rollback()
		return nil, parsePQError(err.(*pq.Error))
	}

	query, args := repo.addToInsert(createQueryLanguagesStart, createQueryLanguagesEnd, languages, message.ID)
	if _, err = tx.Exec(createQueryLanguagesStart, args); err != nil {
		_ = tx.Rollback()
		return nil, parsePQError(err.(*pq.Error))
	}

	query, args = repo.addToInsert(createQueryCountriesStart, createQueryCountriesEnd, counties, message.ID)
	if _, err = tx.Exec(query, args); err != nil {
		_ = tx.Rollback()
		return nil, parsePQError(err.(*pq.Error))
	}

	if _, err = tx.Exec(searchUserQuery, message.ID); err != nil {
		_ = tx.Rollback()
		return nil, parsePQError(err.(*pq.Error))
	}

	return message, nil
}

// GetGotten Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) GetGotten(user string) ([]models.GlideMessage, error) {

}

// GetSent Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) GetSent(user string) ([]models.GlideMessage, error) {

}

// UpdatePicture Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) UpdatePicture(msgId int64, picture string) error {

}

// Check Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) Check(id int64) error {

}

// Get Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) Get(id int64) (*models.GlideMessage, error) {

}

// ChangeUser Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) ChangeUser(id int64) error {

}

// CheckAllow Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) CheckAllow(user string, chatId int64) error {
	if err := repo.store.QueryRowx(checkAllowQuery, chatId, user).Scan(&chatId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// CheckChat Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) CheckChat(chatId int64) error {
	if err := repo.store.QueryRowx(checkQuery, chatId).Scan(&chatId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// GetChats Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) GetChats(userId string) ([]models.Chat, error) {
	rows, err := repo.store.Queryx(getChatsQuery, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	var chats []models.Chat

	for rows.Next() {
		var chat models.Chat
		err = rows.Scan(
			&chat.ID,
			&chat.Companion,
			&chat.CompanionAvatar,
			&chat.LastMessage.ID,
			&chat.LastMessage.Text,
			&chat.LastMessage.Picture,
			&chat.LastMessage.Author,
			&chat.LastMessage.IsViewed,
			&chat.LastMessage.Created)

		if err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return chats, nil
}

// GetMessages Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) GetMessages(chatId int64, pag *models.Pagination) ([]models.Message, error) {
	rows, err := repo.store.Queryx(getMessagesQuery, chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	var chats []models.Message

	for rows.Next() {
		var chat models.Message
		err = rows.Scan(
			&chat.ID,
			&chat.Text,
			&chat.Picture,
			&chat.Author,
			&chat.IsViewed,
			&chat.Created)

		if err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return chats, nil
}

// MarkMessages Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) MarkMessages(chatId int64, messageIds []int64) error {
	query, args, err := sqlx.In(markMessages, messageIds)
	if err != nil {
		return repository.NewDBError(err)
	}

	query = putilits.CustomRebind(2, query)

	res, err := repo.store.Exec(query, chatId, args)
	if err != nil {
		return repository.NewDBError(err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return repository.NewDBError(err)
	}

	if n == 0 {
		return repository.NotFound
	}

	return nil
}

// CreateMessage Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) CreateMessage(text string, chatId int64, image string, user string) (*models.Message, error) {
	ms := &models.Message{}
	if err := repo.store.QueryRowx(createMessageQuery, text, chatId, image, user).
		Scan(&ms.ID, &ms.Text, &ms.Picture, &ms.Author, &ms.IsViewed, &ms.Created); err != nil {
		return nil, repository.NewDBError(err)
	}

	return ms, nil
}

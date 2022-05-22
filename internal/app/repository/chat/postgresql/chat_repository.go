package repository_postgresql

import (
	"database/sql"
	"github.com/lib/pq"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	repository_chat "glide/internal/app/repository/chat"
	putilits "glide/internal/pkg/utilits/postgresql"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

const (
	createQuery = `INSERT INTO chat (author, companion) VALUES ($1, $2) 
		RETURNING id, author, companion`

	checkQuery = `SELECT id FROM chat WHERE id = $1`

	checkAllowQuery = `SELECT id FROM chat WHERE id = $1 and (companion = $2 or author = $2)`

	createMessageQuery = `INSERT INTO messages (message, chat, picture, author) VALUES ($1, $2, $3, $4) 
		RETURNING id, message, picture, author, is_viewed, created`

	getChatQuery = `
			SELECT id, companion FROM chat WHERE id = $1 and author = $2
			UNION 
			SELECT id, author FROM chat WHERE id = $1 and companion = $2
			`

	getChatsQuery = `
			WITH latest_messages as (
				SELECT min(created) as data FROM messages GROUP BY chat
			), countView as (
			    SELECT count(*) as cnt, chat FROM messages WHERE messages.author = $1 GROUP BY chat
			)
				SELECT chat.id, chat.author, u.avatar, cv.cnt, m.id, m.message, m.picture, m.author, m.is_viewed, m.created FROM chat
					 LEFT JOIN messages m on chat.id = m.chat and m.created in(SELECT data FROM latest_messages)
				     LEFT JOIN countView cv on cv.chat = m.chat
					 JOIN users u on chat.author = u.nickname
				WHERE chat.companion = $1
				UNION
				SELECT chat.id, chat.companion, u.avatar, cv.cnt, m.id, m.message, m.picture, m.author, m.is_viewed, m.created FROM chat
					LEFT JOIN messages m on chat.id = m.chat and m.created in(SELECT data FROM latest_messages)
				    LEFT JOIN countView cv on cv.chat = m.chat
					JOIN users u on chat.companion = u.nickname
				WHERE chat.author = $1
	`

	getMessagesQuery = `
				SELECT id, message, picture, author, is_viewed, created FROM messages WHERE chat = $1`

	markMessages = ` UPDATE messages SET is_viewed=true WHERE chat = $1 and id in (?)`
)

type ChatRepository struct {
	store *sqlx.DB
}

var _ = repository_chat.Repository(&ChatRepository{})

func NewChatRepository(st *sqlx.DB) *ChatRepository {
	return &ChatRepository{
		store: st,
	}
}

// Create Errors:
//		ChatAlreadyExists
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *ChatRepository) Create(user string, with string) (*models.Chat, error) {
	id := int64(0)
	if err := repo.store.QueryRowx(createQuery, user, with).
		Scan(&id, &user, &with); err != nil {
		return nil, parsePQError(err.(*pq.Error))
	}

	return &models.Chat{ID: id, Companion: with}, nil
}

// GetChat Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *ChatRepository) GetChat(chatId int64, author string) (*models.Chat, error) {
	chat := &models.Chat{}
	if err := repo.store.QueryRowx(getChatQuery, chatId, author).
		Scan(
			&chat.ID,
			&chat.Companion,
			&chat.CompanionAvatar,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return chat, nil
}

// CheckAllow Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *ChatRepository) CheckAllow(user string, chatId int64) error {
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
func (repo *ChatRepository) CheckChat(chatId int64) error {
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
func (repo *ChatRepository) GetChats(userId string) ([]models.Chat, error) {
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

		msgId := sql.NullInt64{}
		chatCountNotViewed := sql.NullInt64{}
		msgText := sql.NullString{}
		msgPicture := sql.NullString{}
		msgAuthor := sql.NullString{}
		msgViwed := sql.NullBool{}
		msgCreated := sql.NullTime{}
		err = rows.Scan(
			&chat.ID,
			&chat.Companion,
			&chat.CompanionAvatar,
			&chatCountNotViewed,
			&msgId,
			&msgText,
			&msgPicture,
			&msgAuthor,
			&msgViwed,
			&msgCreated,
		)

		if chatCountNotViewed.Valid {
			chat.CountNotViewed = chatCountNotViewed.Int64
		} else {
			chat.CountNotViewed = 0
		}

		if msgId.Valid {
			chat.LastMessage = &models.Message{
				ID:       msgId.Int64,
				Text:     msgText.String,
				Picture:  msgPicture.String,
				Author:   msgAuthor.String,
				IsViewed: msgViwed.Bool,
				Created:  msgCreated.Time,
			}
		} else {
			chat.LastMessage = nil
		}

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
func (repo *ChatRepository) GetMessages(chatId int64, pag *models.Pagination) ([]models.Message, error) {
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
func (repo *ChatRepository) MarkMessages(chatId int64, messageIds []int64) error {
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
func (repo *ChatRepository) CreateMessage(text string, chatId int64, image string, user string) (*models.Message, error) {
	ms := &models.Message{}
	if err := repo.store.QueryRowx(createMessageQuery, text, chatId, image, user).
		Scan(&ms.ID, &ms.Text, &ms.Picture, &ms.Author, &ms.IsViewed, &ms.Created); err != nil {
		return nil, repository.NewDBError(err)
	}

	return ms, nil
}

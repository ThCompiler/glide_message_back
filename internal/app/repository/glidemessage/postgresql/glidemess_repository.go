package repository_postgresql

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	repository_glidemess "glide/internal/app/repository/glidemessage"
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

	createQueryLanguagesEnd = `)
						)
						INSERT INTO glide_message_languages (language, glide_message) SELECT lng_name, ? FROM lng`

	createQueryCountriesStart = `
						WITH cnt AS (
    						SELECT country_name as cnt_name FROM countries WHERE lower(country_name) in (
							`

	createQueryCountriesEnd = `)
						)
						INSERT INTO glide_message_countries (country, glide_message) SELECT cnt_name, ? FROM cnt`

	updateVisitedUserQuery = `
				UPDATE glide_users SET is_actual = false WHERE glide_message = $1 AND is_actual
			`

	updatePictureQuery = `
				UPDATE glide_message SET picture = $2 WHERE id = $1
			`

	deleteGlideMessage = `
				DELETE FROM glide_message WHERE id = $1
			`

	searchUserQuery = `
				WITH usr_f AS (
					SELECT usr.nickname as nick FROM users as usr
					 JOIN user_language ul ON usr.nickname = ul.nickname
						AND ul.language in (SELECT gml.language FROM glide_message_languages as gml WHERE gml.glide_message = $1)
					WHERE usr.country in (SELECT gmc.country FROM glide_message_countries as gmc WHERE gmc.glide_message = $1)
					  and usr.nickname not in (SELECT visited_user FROM glide_users WHERE glide_users.glide_message = $1)
				), usr_s AS (
					SELECT usr.nickname as nick FROM users as usr
					WHERE usr.nickname not in (
						SELECT visited_user FROM glide_users WHERE glide_users.glide_message = $1
						UNION
						SELECT nick FROM usr_f
					)
				)
				INSERT INTO glide_users (visited_user, glide_message) SELECT COALESCE(
								(SELECT usr_f.nick FROM usr_f OFFSET random() * (SELECT count(*) from usr_f) LIMIT 1),
								(SELECT usr_s.nick FROM usr_s OFFSET random() * (SELECT count(*) from usr_s) LIMIT 1)
            	), $1 RETURNING visited_user
 			`

	getMessagesQuery = `
				SELECT id, title, message, picture, author, country, created FROM glide_message where id = $1 `

	checkMessagesQuery = `
				SELECT id FROM glide_message where id = $1 `

	checkAllowUserQuery = `
				SELECT glide_message FROM glide_users where glide_message = $1 and visited_user = $2 and is_actual`

	addVisitedUserQuery = `
				INSERT INTO glide_users (visited_user, glide_message, is_actual) VALUES ($2, $1, false)`

	checkAllowAuthorQuery = `
				SELECT id FROM glide_message where id = $1 and author = $2`

	getGottenMessagesQuery = `
				SELECT glide_message.id, title, message, picture, author, country, created FROM glide_message
						JOIN glide_users as gu on gu.glide_message = glide_message.id and gu.visited_user = $1 and gu.is_actual`

	getSentMessagesQuery = `
				SELECT id, title, message, picture, author, country, created FROM glide_message WHERE author = $1`
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
func (repo *GlideMessageRepository) Create(message *models.GlideMessage,
	languages []string, counties []string) (*models.GlideMessage, string, error) {
	tx, err := repo.store.Beginx()
	if err != nil {
		return nil, "", repository.NewDBError(err)
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
		if _, can := err.(*pq.Error); can {
			return nil, "", parsePQError(err.(*pq.Error))
		}
		return nil, "", repository.NewDBError(err)
	}

	if len(languages) != 0 {
		query, args := repo.addToInsert(createQueryLanguagesStart, createQueryLanguagesEnd, languages, message.ID)
		if _, err = tx.Exec(query, args...); err != nil {
			_ = tx.Rollback()
			if _, can := err.(*pq.Error); can {
				return nil, "", parsePQError(err.(*pq.Error))
			}
			return nil, "", repository.NewDBError(err)
		}
	}

	if len(counties) != 0 {
		query, args := repo.addToInsert(createQueryCountriesStart, createQueryCountriesEnd, counties, message.ID)
		if _, err = tx.Exec(query, args...); err != nil {
			_ = tx.Rollback()
			if _, can := err.(*pq.Error); can {
				return nil, "", parsePQError(err.(*pq.Error))
			}
			return nil, "", repository.NewDBError(err)
		}
	}

	if _, err = tx.Exec(addVisitedUserQuery, message.ID, message.Author); err != nil {
		_ = tx.Rollback()
		return nil, "", repository.NewDBError(err)
	}

	if _, err = tx.Exec(updateVisitedUserQuery, message.ID); err != nil {
		_ = tx.Rollback()
		return nil, "", repository.NewDBError(err)
	}

	nickname := ""

	if err = tx.QueryRowx(searchUserQuery, message.ID).Scan(&nickname); err != nil {
		_ = tx.Rollback()
		return nil, "", repository.NewDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, "", repository.NewDBError(err)
	}

	return message, nickname, nil
}

// GetGotten Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) GetGotten(user string) ([]models.GlideMessage, error) {
	rows, err := repo.store.Queryx(getGottenMessagesQuery, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.GlideMessage{}, nil
		}
		return nil, repository.NewDBError(err)
	}

	var msgs []models.GlideMessage

	for rows.Next() {
		var msg models.GlideMessage
		err = rows.Scan(
			&msg.ID,
			&msg.Title,
			&msg.Message,
			&msg.Picture,
			&msg.Author,
			&msg.Country,
			&msg.Created)

		if err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		msgs = append(msgs, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return msgs, nil
}

// GetSent Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) GetSent(user string) ([]models.GlideMessage, error) {
	rows, err := repo.store.Queryx(getSentMessagesQuery, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.GlideMessage{}, nil
		}
		return nil, repository.NewDBError(err)
	}

	var msgs []models.GlideMessage

	for rows.Next() {
		var msg models.GlideMessage
		err = rows.Scan(
			&msg.ID,
			&msg.Title,
			&msg.Message,
			&msg.Picture,
			&msg.Author,
			&msg.Country,
			&msg.Created)

		if err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		msgs = append(msgs, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return msgs, nil
}

// UpdatePicture Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) UpdatePicture(msgId int64, picture string) error {
	res, err := repo.store.Exec(updatePictureQuery, msgId, picture)
	if err != nil {
		return repository.NewDBError(err)
	}

	rw, err := res.RowsAffected()
	if err != nil {
		return repository.NewDBError(err)
	}

	if rw == 0 {
		return repository.NotFound
	}

	return nil
}

// Delete Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) Delete(msgId int64) error {
	res, err := repo.store.Exec(deleteGlideMessage, msgId)
	if err != nil {
		return repository.NewDBError(err)
	}

	rw, err := res.RowsAffected()
	if err != nil {
		return repository.NewDBError(err)
	}

	if rw == 0 {
		return repository.NotFound
	}

	return nil
}

// Check Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) Check(id int64) error {
	if err := repo.store.QueryRowx(checkMessagesQuery, id).
		Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// Get Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) Get(id int64) (*models.GlideMessage, error) {
	msg := &models.GlideMessage{}

	if err := repo.store.QueryRowx(getMessagesQuery, id).
		Scan(
			&msg.ID,
			&msg.Title,
			&msg.Message,
			&msg.Picture,
			&msg.Author,
			&msg.Country,
			&msg.Created,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return msg, nil
}

// ChangeUser Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) ChangeUser(id int64) (string, error) {
	if err := repo.Check(id); err != nil {
		return "", err
	}

	tx, err := repo.store.Beginx()
	if err != nil {
		return "", repository.NewDBError(err)
	}

	if _, err = tx.Exec(updateVisitedUserQuery, id); err != nil {
		_ = tx.Rollback()
		return "", repository.NewDBError(err)
	}

	nickname := ""

	if err = tx.QueryRowx(searchUserQuery, id).Scan(&nickname); err != nil {
		_ = tx.Rollback()
		return "", repository.NewDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return "", repository.NewDBError(err)
	}

	return nickname, nil
}

// CheckAllowUser Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) CheckAllowUser(id int64, user string) error {
	if err := repo.store.QueryRowx(checkAllowUserQuery, id, user).
		Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// CheckAllowAuthor Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *GlideMessageRepository) CheckAllowAuthor(id int64, user string) error {
	if err := repo.store.QueryRowx(checkAllowAuthorQuery, id, user).
		Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

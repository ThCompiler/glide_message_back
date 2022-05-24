package repository_postgresql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"glide/internal/app/models"
	"glide/internal/app/repository"

	"github.com/pkg/errors"
)

const (
	updateNicknameQuery = `UPDATE users SET nickname = $1 WHERE nickname = $2`

	updateAvatarQuery = `UPDATE users SET avatar = $1 WHERE nickname = $2`

	getPasswordQuery = `SELECT password FROM users WHERE nickname=$1`

	findByNicknameQuery = `SELECT nickname, fullname, about, age, avatar, lower(country) FROM users WHERE nickname=$1`

	getBlacklistQuery = `SELECT users.nickname, users.fullname, users.about, users.avatar FROM black_list
							  JOIN users on black_list.fobbiged_user = users.nickname WHERE author=$1`

	addToBlackListQuery = `
					INSERT INTO black_list (fobbiged_user, author) VALUES ($1, $2)
		`

	deleteFromBlackListQuery = `
					DELETE FROM black_list WHERE fobbiged_user = $1 and author = $2
		`

	findByNicknameGetLanguagesQuery = `SELECT lower(language) FROM user_language WHERE nickname=$1`

	createQuery = `    
						WITH cnt AS (
						    SELECT country_name as country_nm FROM countries WHERE lower(country_name) = lower($6)
						), sel AS (
						    SELECT nickname, fullname, about, password, age, country
							FROM users
							WHERE nickname = $1 LIMIT 1
						), ins as (
							INSERT INTO users (nickname, fullname, about, password, age, country)
								SELECT $1, $2, $3, $4, $5, (select country_nm from cnt limit 1)
								WHERE not exists (select 1 from sel)
							RETURNING nickname, fullname, about, age, country
						)
						SELECT nickname, fullname, about, age, country, 0
						FROM ins
						UNION ALL
						SELECT nickname, fullname, about, age, country, 1
						FROM sel
					`
	addLanguagesToUsersQuery = `
						WITH lng AS (
						    SELECT language as lng_name FROM languages WHERE lower(language) = lower($1)
						)
						INSERT INTO user_language (language, nickname) VALUES ((select lng_name from lng limit 1), $2)`

	deleteLanguagesForUsersQuery = `DELETE FROM user_language WHERE nickname = $1`

	updateUserQuery = `
					WITH cnt AS (
						SELECT country_name as country_nm FROM countries WHERE lower(country_name) = lower(NULLIF(TRIM($4), ''))
					)
					UPDATE users SET 
					    fullname = COALESCE(NULLIF(TRIM($1), ''), fullname),
					    about = COALESCE(NULLIF(TRIM($2), ''), about),
					    age = COALESCE(NULLIF($3, 0), age),
						country = COALESCE((select country_nm from cnt limit 1), country)
					WHERE nickname = $5
					RETURNING nickname, fullname, about, age, country`
)

type UserRepository struct {
	store *sqlx.DB
}

func NewUserRepository(st *sqlx.DB) *UserRepository {
	return &UserRepository{
		store: st,
	}
}

// Create Errors:
// 		NicknameAlreadyExist
//		IncorrectCounty
//		IncorrectLanguage
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) Create(u *models.User) (*models.User, error) {
	tx, err := repo.store.Beginx()
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	var have_conflict = 0

	country := sql.NullString{u.Country,
		u.Country != "",
	}

	if err = tx.QueryRow(
		createQuery,
		u.Nickname,
		u.Fullname,
		u.About,
		u.EncryptedPassword,
		u.Age,
		country).
		Scan(&u.Nickname,
			&u.Fullname,
			&u.About,
			&u.Age,
			&country,
			&have_conflict); err != nil {
		_ = tx.Rollback()
		return nil, parsePQError(err.(*pq.Error))
	}

	if country.Valid {
		u.Country = country.String
	} else {
		u.Country = ""
	}

	if have_conflict == 1 {
		return u, NicknameAlreadyExist
	}

	for _, lang := range u.Languages {
		if _, err = tx.Exec(
			addLanguagesToUsersQuery,
			lang,
			u.Nickname); err != nil {
			_ = tx.Rollback()
			return nil, parsePQError(err.(*pq.Error))
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return u, nil
}

// FindByNickname Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) FindByNickname(nickname string) (*models.User, error) {
	user := models.User{}

	country := sql.NullString{}

	if err := repo.store.QueryRow(findByNicknameQuery, nickname).
		Scan(&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Age,
			&user.Avatar,
			&country); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	if country.Valid {
		user.Country = country.String
	} else {
		user.Country = ""
	}

	if err := repo.store.Select(&user.Languages, findByNicknameGetLanguagesQuery, nickname); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return &user, nil
}

// GetBlacklist Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) GetBlacklist(nickname string) ([]models.User, error) {
	rows, err := repo.store.Queryx(getBlacklistQuery, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Avatar,
		)

		if err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return users, nil
}

// GetPassword Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) GetPassword(nickname string) (string, error) {
	password := ""

	if err := repo.store.QueryRow(getPasswordQuery, nickname).
		Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			return "", repository.NotFound
		}
		return "", repository.NewDBError(err)
	}

	return password, nil
}

// UpdatePassword Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) UpdatePassword(id int64, newEncryptedPassword string) error {
	query := `UPDATE users SET encrypted_password = $1 WHERE users_id = $2`

	row, err := repo.store.Query(query,
		newEncryptedPassword, id)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}

// UpdateAvatar Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) UpdateAvatar(nickname string, newAvatar string) error {
	res, err := repo.store.Exec(updateAvatarQuery, newAvatar, nickname)
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

// Update Errors:
//		IncorrectCounty
//		IncorrectLanguage
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) Update(u *models.User) (*models.User, error) {
	tx, err := repo.store.Beginx()
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	country := sql.NullString{u.Country,
		u.Country != "",
	}

	if err = tx.QueryRow(
		updateUserQuery,
		u.Fullname,
		u.About,
		u.Age,
		country,
		u.Nickname).
		Scan(&u.Nickname,
			&u.Fullname,
			&u.About,
			&u.Age,
			&country); err != nil {
		_ = tx.Rollback()
		return nil, parsePQError(err.(*pq.Error))
	}

	if country.Valid {
		u.Country = country.String
	} else {
		u.Country = ""
	}

	if len(u.Languages) != 0 {
		if _, err = tx.Exec(deleteLanguagesForUsersQuery, u.Nickname); err != nil {
			_ = tx.Rollback()
			return nil, repository.NewDBError(err)
		}

		for _, lang := range u.Languages {
			if _, err = tx.Exec(
				addLanguagesToUsersQuery,
				lang,
				u.Nickname); err != nil {
				_ = tx.Rollback()
				return nil, parsePQError(err.(*pq.Error))
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return u, nil
}

// AddToBlacklist Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) AddToBlacklist(author string, nickname string) error {
	_, err := repo.store.Exec(addToBlackListQuery, nickname, author)
	if err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// DeleteFromBlacklist Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) DeleteFromBlacklist(author string, nickname string) error {
	_, err := repo.store.Exec(deleteFromBlackListQuery, nickname, author)
	if err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// UpdateNickname Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) UpdateNickname(oldNickname string, newNickname string) error {
	row, err := repo.store.Exec(updateNicknameQuery, newNickname, oldNickname)
	if err != nil {
		return repository.NewDBError(err)
	}
	if cntChangesRows, err := row.RowsAffected(); err != nil || cntChangesRows != 1 {
		if err != nil {
			return repository.NewDBError(err)
		}
		return repository.NewDBError(
			errors.Wrapf(err,
				"UPDATE_NICKNAME_REPO: expected changes only one row in db, change %v", cntChangesRows))
	}
	return nil
}

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

	findByNicknameQuery             = `SELECT nickname, fullname, about, age, country FROM users WHERE nickname=$1`
	findByNicknameGetLanguagesQuery = `SELECT language FROM user_language WHERE nickname=$1`

	createQuery = `    
						WITH sel AS (
						    SELECT nickname, fullname, about, password, age, country
							FROM users
							WHERE nickname = $1 LIMIT 1
						), ins as (
							INSERT INTO users (nickname, fullname, about, password, age, country)
								SELECT $1, $2, $3, $4, $5, $6
								WHERE not exists (select 1 from sel)
							RETURNING nickname, fullname, about, age, country
						)
						SELECT nickname, fullname, about, age, country, 0
						FROM ins
						UNION ALL
						SELECT nickname, fullname, about, age, country, 1
						FROM sel
					`
	addLanguagesToUsersQuery = `INSERT INTO user_language (language, nickname) VALUES ($1, $2)`

	deleteLanguagesForUsersQuery = `DELETE FROM user_language WHERE nickname = $1`

	updateUserQuery = `
					UPDATE users SET 
					    fullname = COALESCE(NULLIF(TRIM($1), ''), fullname),
					    about = COALESCE(NULLIF(TRIM($2), ''), about),
					    age = COALESCE(NULLIF($3, 0), age),
						country = COALESCE(NULLIF(TRIM($4), ''), country),
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

	for lang := range u.Languages {
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

		for lang := range u.Languages {
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

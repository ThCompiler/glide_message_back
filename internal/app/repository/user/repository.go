package repository_user

import (
	"glide/internal/app/models"
)

type Repository interface {
	// Create Errors:
	// 		repository_postgresql.NicknameAlreadyExist
	//		repository_postgresql.IncorrectCounty
	//		repository_postgresql.IncorrectLanguage
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(*models.User) (*models.User, error)

	// FindByNickname Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	FindByNickname(nickname string) (*models.User, error)

	// Update Errors:
	//		repository_postgresql.IncorrectCounty
	//		repository_postgresql.IncorrectLanguage
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Update(*models.User) (*models.User, error)

	// UpdateAvatar Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	UpdateAvatar(id string, newAvatar string) error
}

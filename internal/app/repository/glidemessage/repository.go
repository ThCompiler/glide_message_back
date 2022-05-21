package repository_glidemess

import (
	"glide/internal/app/models"
)

type Repository interface {
	// Create Errors:
	//		IncorrectCounty
	//		IncorrectLanguage
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(message *models.GlideMessage, languages []string, counties []string) (*models.GlideMessage, error)

	// GetGotten Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetGotten(user string) ([]models.GlideMessage, error)

	// GetSent Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetSent(user string) ([]models.GlideMessage, error)

	// UpdatePicture Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	UpdatePicture(msgId int64, picture string) error

	// Check Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Check(id int64) error

	// Get Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Get(id int64) (*models.GlideMessage, error)

	// ChangeUser Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	ChangeUser(id int64) error
}

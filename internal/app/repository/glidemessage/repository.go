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
	Create(message *models.GlideMessage,
		languages []string, counties []string, age int64) (*models.GlideMessage, string, error)

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

	// Delete Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Delete(msgId int64) error

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
	ChangeUser(id int64) (string, error)

	// CheckAllowUser Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	CheckAllowUser(id int64, user string) error

	// CheckAllowAuthor Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	CheckAllowAuthor(id int64, user string) error
}

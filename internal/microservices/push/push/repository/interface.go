package repository

type Repository interface {
	// GetUserAvatar Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetUserAvatar(username string) (avatar string, err error)

	// GetMessageInfo Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetMessageInfo(messageId int64) (author string, text string, chatId int64, err error)

	// GetGlideInfo Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetGlideInfo(glideId int64) (author string, message string, title string, country string, err error)
}

package usecase_access

//go:generate mockgen -destination=mocks/mock_access_usecase.go -package=mock_usecase -mock_names=Usecase=AccessUsecase . Usecase

type Usecase interface {
	// CheckAccess Errors:
	//		NoAccess
	//		FirstQuery
	// 		app.GeneralError with Errors
	// 			repository_access.InvalidStorageData
	CheckAccess(userIp string) (bool, error)

	// Create Errors:
	// 		app.GeneralError with Errors
	// 			repository_access.SetError
	Create(userIp string) (bool, error)

	// Update Errors:
	//		NoAccess
	// 		app.GeneralError with Errors
	// 			repository_access.InvalidStorageData
	Update(userIp string) (int64, error)

	// AddToBlackList Errors:
	// 		app.GeneralError with Errors
	// 			repository_access.InvalidStorageData
	AddToBlackList(userIp string) error

	// CheckBlackList Errors:
	// 		app.GeneralError with Errors
	// 			repository_access.InvalidStorageData
	CheckBlackList(userIp string) (bool, error)
}

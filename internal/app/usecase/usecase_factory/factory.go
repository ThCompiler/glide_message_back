package usecase_factory

import (
	"glide/internal/app"
	useUser "glide/internal/app/usecase/user"
)

type UsecaseFactory struct {
	paymentsConfig    app.Payments
	repositoryFactory RepositoryFactory
	userUsecase       useUser.Usecase
}

func NewUsecaseFactory(repositoryFactory RepositoryFactory, paymentsConf app.Payments) *UsecaseFactory {
	return &UsecaseFactory{
		repositoryFactory: repositoryFactory,
		paymentsConfig:    paymentsConf,
	}
}

func (f *UsecaseFactory) GetUserUsecase() useUser.Usecase {
	if f.userUsecase == nil {
		f.userUsecase = useUser.NewUserUsecase(f.repositoryFactory.GetUserRepository(),
			f.repositoryFactory.GetFileRepository())
	}
	return f.userUsecase
}

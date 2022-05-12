package usecase_factory

import (
	useUser "glide/internal/app/usecase/user"
)

type UsecaseFactory struct {
	repositoryFactory RepositoryFactory
	userUsecase       useUser.Usecase
}

func NewUsecaseFactory(repositoryFactory RepositoryFactory) *UsecaseFactory {
	return &UsecaseFactory{
		repositoryFactory: repositoryFactory,
	}
}

func (f *UsecaseFactory) GetUserUsecase() useUser.Usecase {
	if f.userUsecase == nil {
		f.userUsecase = useUser.NewUserUsecase(f.repositoryFactory.GetUserRepository(),
			f.repositoryFactory.GetFileRepository())
	}
	return f.userUsecase
}

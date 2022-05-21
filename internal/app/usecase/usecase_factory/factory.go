package usecase_factory

import (
	useChat "glide/internal/app/usecase/chats"
	useUser "glide/internal/app/usecase/user"
)

type UsecaseFactory struct {
	repositoryFactory RepositoryFactory
	userUsecase       useUser.Usecase
	chatUsecase       useChat.Usecase
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

func (f *UsecaseFactory) GetChatUsecase() useChat.Usecase {
	if f.chatUsecase == nil {
		f.chatUsecase = useChat.NewChatsUsecase(f.repositoryFactory.GetChatRepository(),
			f.repositoryFactory.GetFileRepository(), f.repositoryFactory.GetPusher())
	}
	return f.chatUsecase
}

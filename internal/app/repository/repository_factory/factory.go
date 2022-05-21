package repository_factory

import (
	"github.com/sirupsen/logrus"
	"glide/internal/app"
	repoChats "glide/internal/app/repository/chat"
	repChatsPsql "glide/internal/app/repository/chat/postgresql"
	repoFiles "glide/internal/app/repository/files"
	repository_os "glide/internal/app/repository/files/os"
	repUser "glide/internal/app/repository/user"
	repUserPsql "glide/internal/app/repository/user/postgresql"
	push_client "glide/internal/microservices/push/delivery/client"
)

type RepositoryFactory struct {
	expectedConnections app.ExpectedConnections
	logger              *logrus.Logger
	userRepository      repUser.Repository
	fileRepository      repoFiles.Repository
	chatRepository      repoChats.Repository
	pusher              push_client.Pusher
}

func NewRepositoryFactory(logger *logrus.Logger, expectedConnections app.ExpectedConnections) *RepositoryFactory {
	return &RepositoryFactory{
		expectedConnections: expectedConnections,
		logger:              logger,
	}
}

func (f *RepositoryFactory) GetUserRepository() repUser.Repository {
	if f.userRepository == nil {
		f.userRepository = repUserPsql.NewUserRepository(f.expectedConnections.SqlConnection)
	}
	return f.userRepository
}

func (f *RepositoryFactory) GetFileRepository() repoFiles.Repository {
	if f.fileRepository == nil {
		f.fileRepository = repository_os.NewFileRepository(f.expectedConnections.PathFiles)
	}
	return f.fileRepository
}

func (f *RepositoryFactory) GetChatRepository() repoChats.Repository {
	if f.chatRepository == nil {
		f.chatRepository = repChatsPsql.NewChatRepository(f.expectedConnections.SqlConnection)
	}
	return f.chatRepository
}

func (f *RepositoryFactory) GetPusher() push_client.Pusher {
	if f.pusher == nil {
		f.pusher = push_client.NewPushSender(f.expectedConnections.RabbitSession)
	}
	return f.pusher
}

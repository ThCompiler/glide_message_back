package usecase_factory

import (
	repoChats "glide/internal/app/repository/chat"
	repFiles "glide/internal/app/repository/files"
	repoGlideMess "glide/internal/app/repository/glidemessage"
	repInfo "glide/internal/app/repository/info"
	repUser "glide/internal/app/repository/user"
	push_client "glide/internal/microservices/push/delivery/client"
)

//go:generate mockgen -destination=mocks/mock_repository_factory.go -package=mock_repository_factory . RepositoryFactory

type RepositoryFactory interface {
	GetUserRepository() repUser.Repository
	GetFileRepository() repFiles.Repository
	GetChatRepository() repoChats.Repository
	GetGlideMessageRepository() repoGlideMess.Repository
	GetInfoRepository() repInfo.Repository
	GetPusher() push_client.Pusher
}

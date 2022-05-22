package handler_factory

import (
	useChat "glide/internal/app/usecase/chats"
	useGlideMess "glide/internal/app/usecase/glidemessage"
	"glide/internal/app/usecase/info"
	useUser "glide/internal/app/usecase/user"
)

//go:generate mockgen -destination=mocks/mock_usecase_factory.go -package=mock_usecase_factory . UsecaseFactory

type UsecaseFactory interface {
	GetUserUsecase() useUser.Usecase
	GetChatUsecase() useChat.Usecase
	GetGlideMessageUsecase() useGlideMess.Usecase
	GetInfoUsecase() info.Usecase
}

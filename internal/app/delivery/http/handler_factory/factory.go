package handler_factory

import (
	"glide/internal/app"
	"glide/internal/app/delivery/http/handlers"
	"glide/internal/app/delivery/http/handlers/chat_handler"
	"glide/internal/app/delivery/http/handlers/chat_id_message_handler"
	"glide/internal/app/delivery/http/handlers/chat_message_handler"
	"glide/internal/app/delivery/http/handlers/login_handler"
	"glide/internal/app/delivery/http/handlers/logout_handler"
	"glide/internal/app/delivery/http/handlers/user_handler"
	upd_user_avatar_handler "glide/internal/app/delivery/http/handlers/user_handler/upd_avatar_handler"
	user_nickname_profile_handler "glide/internal/app/delivery/http/handlers/user_nickname_handler/profile_handler"
	"glide/internal/microservices/auth/delivery/grpc/client"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
)

const (
	ROOT = iota
	LOGIN
	LOGOUT
	USER
	USER_NICKNAME_PROFILE
	USER_BLACKLIST
	USER_NICKNAME_BLACKLIST
	USER_UPDATE_AVATAR
	CHAT
	CHAT_ID_MESSAGE
	CHAT_ID_MESSAGE_STATUS
)

type HandlerFactory struct {
	usecaseFactory    UsecaseFactory
	sessionClientConn *grpc.ClientConn
	logger            *logrus.Logger
	urlHandler        *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, usecaseFactory UsecaseFactory, sClientConn *grpc.ClientConn) *HandlerFactory {
	return &HandlerFactory{
		usecaseFactory:    usecaseFactory,
		logger:            logger,
		sessionClientConn: sClientConn,
	}
}

func (f *HandlerFactory) initAllHandlers() map[int]app.Handler {
	ucUser := f.usecaseFactory.GetUserUsecase()
	ucChat := f.usecaseFactory.GetChatUsecase()
	sManager := client.NewSessionClient(f.sessionClientConn)

	return map[int]app.Handler{
		LOGIN:                  login_handler.NewLoginHandler(f.logger, sManager, ucUser),
		LOGOUT:                 logout_handler.NewLogoutHandler(f.logger, sManager),
		USER:                   user_handler.NewProfileHandler(f.logger, sManager, ucUser),
		USER_NICKNAME_PROFILE:  user_nickname_profile_handler.NewGetProfileHandler(f.logger, sManager, ucUser),
		USER_UPDATE_AVATAR:     upd_user_avatar_handler.NewUpdateUserAvatarHandler(f.logger, sManager, ucUser),
		CHAT:                   chat_handler.NewChatHandler(f.logger, sManager, ucChat),
		CHAT_ID_MESSAGE:        chat_id_message_handler.NewChatIdMessageHandler(f.logger, sManager, ucChat),
		CHAT_ID_MESSAGE_STATUS: chat_message_handler.NewChatMessageHandler(f.logger, sManager, ucChat),
	}
}

func (f *HandlerFactory) GetHandleUrls() *map[string]app.Handler {
	if f.urlHandler != nil {
		return f.urlHandler
	}

	hs := f.initAllHandlers()
	f.urlHandler = &map[string]app.Handler{
		//"/":                     "I am a joke?",
		"/login":  hs[LOGIN],
		"/logout": hs[LOGOUT],
		// /user     ---------------------------------------------------------////
		"/user": hs[USER],
		"/user/{" + handlers.UserNickname + "}/profile": hs[USER_NICKNAME_PROFILE],
		"/user/update/avatar":                           hs[USER_UPDATE_AVATAR],
		// /chat     ---------------------------------------------------------////
		"/chat": hs[CHAT],
		"/chat/{" + handlers.ChatId + "}/message":        hs[CHAT_ID_MESSAGE],
		"/chat/{" + handlers.ChatId + "}/message/status": hs[CHAT_ID_MESSAGE_STATUS],
	}
	return f.urlHandler
}

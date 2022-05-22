package handler_factory

import (
	"glide/internal/app"
	"glide/internal/app/delivery/http/handlers"
	chat_handler "glide/internal/app/delivery/http/handlers/chat"
	"glide/internal/app/delivery/http/handlers/chat_id/message"
	"glide/internal/app/delivery/http/handlers/chat_id/message_status"
	glide_create_handler "glide/internal/app/delivery/http/handlers/glide/create"
	glide_gotten_handler "glide/internal/app/delivery/http/handlers/glide/gotten"
	glide_id_apply_handler "glide/internal/app/delivery/http/handlers/glide/id/apply"
	glide_id_info_handler "glide/internal/app/delivery/http/handlers/glide/id/info"
	glide_id_picture_handler "glide/internal/app/delivery/http/handlers/glide/id/picture"
	glide_id_redirect_handler "glide/internal/app/delivery/http/handlers/glide/id/redirect"
	glide_sent_handler "glide/internal/app/delivery/http/handlers/glide/sent"
	"glide/internal/app/delivery/http/handlers/login"
	"glide/internal/app/delivery/http/handlers/logout"
	system_countries_handler "glide/internal/app/delivery/http/handlers/system/countries"
	system_languages_handler "glide/internal/app/delivery/http/handlers/system/languages"
	"glide/internal/app/delivery/http/handlers/user"
	upd_user_avatar_handler "glide/internal/app/delivery/http/handlers/user/upd_avatar"
	user_nickname_profile_handler "glide/internal/app/delivery/http/handlers/user_nickname/profile"
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
	GLIDE_MESSAGE_SENT
	GLIDE_MESSAGE_GOTTEN
	GLIDE_MESSAGE_CREATE
	GLIDE_MESSAGE_ID_PICTURE
	GLIDE_MESSAGE_ID_APPLY
	GLIDE_MESSAGE_ID_REDIRECT
	GLIDE_MESSAGE_ID_INFO
	SYSTEM_LANGUAGES
	SYSTEM_COUNTRIES
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
	ucInfo := f.usecaseFactory.GetInfoUsecase()
	ucGlideMessage := f.usecaseFactory.GetGlideMessageUsecase()
	sManager := client.NewSessionClient(f.sessionClientConn)

	return map[int]app.Handler{
		LOGIN:                     login_handler.NewLoginHandler(f.logger, sManager, ucUser),
		LOGOUT:                    logout_handler.NewLogoutHandler(f.logger, sManager),
		USER:                      user_handler.NewProfileHandler(f.logger, sManager, ucUser),
		USER_NICKNAME_PROFILE:     user_nickname_profile_handler.NewGetProfileHandler(f.logger, sManager, ucUser),
		USER_UPDATE_AVATAR:        upd_user_avatar_handler.NewUpdateUserAvatarHandler(f.logger, sManager, ucUser),
		CHAT:                      chat_handler.NewChatHandler(f.logger, sManager, ucChat),
		CHAT_ID_MESSAGE:           chat_id_message_handler.NewChatIdMessageHandler(f.logger, sManager, ucChat),
		CHAT_ID_MESSAGE_STATUS:    chat_message_handler.NewChatMessageHandler(f.logger, sManager, ucChat),
		GLIDE_MESSAGE_SENT:        glide_sent_handler.NewGlideSentHandler(f.logger, sManager, ucGlideMessage),
		GLIDE_MESSAGE_GOTTEN:      glide_gotten_handler.NewGlideGottenHandler(f.logger, sManager, ucGlideMessage),
		GLIDE_MESSAGE_CREATE:      glide_create_handler.NewGlideIdApplyHandler(f.logger, sManager, ucGlideMessage),
		GLIDE_MESSAGE_ID_PICTURE:  glide_id_picture_handler.NewGlideIdPictureHandler(f.logger, sManager, ucGlideMessage),
		GLIDE_MESSAGE_ID_APPLY:    glide_id_apply_handler.NewGlideIdApplyHandler(f.logger, sManager, ucGlideMessage),
		GLIDE_MESSAGE_ID_REDIRECT: glide_id_redirect_handler.NewGlideIdRedirectHandler(f.logger, sManager, ucGlideMessage),
		GLIDE_MESSAGE_ID_INFO:     glide_id_info_handler.NewGlideIdInfoHandler(f.logger, sManager, ucGlideMessage),
		SYSTEM_LANGUAGES:          system_languages_handler.NewSystemLanguagesHandler(f.logger, ucInfo),
		SYSTEM_COUNTRIES:          system_countries_handler.NewSystemCountriesHandler(f.logger, ucInfo),
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
		// /glide     ---------------------------------------------------------////
		"/glide/sent":   hs[GLIDE_MESSAGE_SENT],
		"/glide/gotten": hs[GLIDE_MESSAGE_GOTTEN],
		"/glide/create": hs[GLIDE_MESSAGE_CREATE],
		"/glide/{" + handlers.GlideMessageId + "}/picture":  hs[GLIDE_MESSAGE_ID_PICTURE],
		"/glide/{" + handlers.GlideMessageId + "}/info":     hs[GLIDE_MESSAGE_ID_INFO],
		"/glide/{" + handlers.GlideMessageId + "}/apply":    hs[GLIDE_MESSAGE_ID_APPLY],
		"/glide/{" + handlers.GlideMessageId + "}/redirect": hs[GLIDE_MESSAGE_ID_REDIRECT],
		// /system     ---------------------------------------------------------////
		"/system/countries": hs[SYSTEM_COUNTRIES],
		"/system/languages": hs[SYSTEM_LANGUAGES],
	}
	return f.urlHandler
}

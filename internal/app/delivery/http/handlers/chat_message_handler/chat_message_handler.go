package chat_message_handler

import (
	"github.com/microcosm-cc/bluemonday"
	"glide/internal/app/delivery/http/handlers"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/delivery/http/models"
	"glide/internal/app/middleware"
	usecase_chats "glide/internal/app/usecase/chats"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ChatMessageHandler struct {
	sessionClient session_client.AuthCheckerClient
	chatsUsecase  usecase_chats.Usecase
	bh.BaseHandler
}

func NewChatMessageHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucChats usecase_chats.Usecase) *ChatMessageHandler {
	h := &ChatMessageHandler{
		sessionClient: sManager,
		chatsUsecase:  ucChats,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodPut, h.PUT, session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
		middleware.NewChatsMiddleware(log, ucChats).CheckCorrectChatIdFunc,
	)

	return h
}

func (h *ChatMessageHandler) PUT(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestMessageIds{}

	err := h.GetRequestBody(r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	chatId, code, err := h.GetInt64FromParam(w, r, handlers.ChatId)

	if err != nil {
		h.Error(w, r, code, err)
		return
	}

	err = h.chatsUsecase.MarkMessages(chatId, req.ToArray())
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	h.Log(r).Debugf("sussess mark messages %v", req)
	w.WriteHeader(http.StatusOK)
}

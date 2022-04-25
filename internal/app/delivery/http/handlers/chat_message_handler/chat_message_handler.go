package chat_message_handler

import (
	"glide/internal/app/delivery/http/handlers/handler_errors"
	models_http "glide/internal/app/delivery/http/models"
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
	h.AddMethod(http.MethodGet, h.PUT,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *ChatMessageHandler) PUT(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	u, err := h.userUsecase.GetProfile(userID.(string))
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Log(r).Debugf("get user %s", u)
	h.Respond(w, r, http.StatusOK, models_http.ToProfileResponse(*u))
}

package chat_id_message_handler

import (
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"glide/internal/app"
	"glide/internal/app/delivery/http/handlers"
	models_http "glide/internal/app/delivery/http/models"
	"glide/internal/app/middleware"
	"glide/internal/app/models"
	usecase_chats "glide/internal/app/usecase/chats"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"glide/internal/pkg/handler/handler_errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ChatIdMessageHandler struct {
	sessionClient session_client.AuthCheckerClient
	chatsUsecase  usecase_chats.Usecase
	bh.BaseHandler
}

func NewChatIdMessageHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucChats usecase_chats.Usecase) *ChatIdMessageHandler {
	h := &ChatIdMessageHandler{
		sessionClient: sManager,
		chatsUsecase:  ucChats,
		BaseHandler:   *bh.NewBaseHandler(log),
	}

	h.AddMiddleware(session_middleware.NewSessionMiddleware(h.sessionClient, log).Check,
		middleware.NewChatsMiddleware(log, ucChats).CheckCorrectChatId)

	h.AddMethod(http.MethodGet, h.GET)

	h.AddMethod(http.MethodPost, h.POST)

	return h
}

func (h *ChatIdMessageHandler) GET(w http.ResponseWriter, r *http.Request) {
	chatId, code, err := h.GetInt64FromParam(w, r, handlers.ChatId)

	if err != nil {
		h.Error(w, r, code, err)
		return
	}

	u, err := h.chatsUsecase.GetMessages(chatId, &models.Pagination{
		Offset: 0,
		Limit:  200,
	})

	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Log(r).Debugf("get messages from chat %d", chatId)
	h.Respond(w, r, http.StatusOK, models_http.ToResponseMessages(u))
}

func (h *ChatIdMessageHandler) POST(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	chatId, code, err := h.GetInt64FromParam(w, r, handlers.ChatId)

	if err != nil {
		h.Error(w, r, code, err)
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	file, filename, code, err := h.GetFilesFromRequest(w, r, handlers.MAX_UPLOAD_SIZE,
		"image", []string{"image/png", "image/jpeg", "image/jpg"})

	if err != nil {
		if _, can := err.(*app.GeneralError); can {
			if !errors.Is(err.(*app.GeneralError).Err, handler_errors.InvalidFormFieldName) {
				h.HandlerError(w, r, code, err)
				return
			}
		} else {
			h.HandlerError(w, r, code, err)
			return
		}
	}

	text := r.FormValue("text")

	if text == "" {
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidBody)
		return
	}

	text = bluemonday.UGCPolicy().Sanitize(text)

	msg, err := h.chatsUsecase.CreateMessage(h.Log(r), text, chatId, file, filename, userID.(string))
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPOST)
		return
	}

	h.Respond(w, r, http.StatusCreated, models_http.ToResponseMessage(*msg))
}

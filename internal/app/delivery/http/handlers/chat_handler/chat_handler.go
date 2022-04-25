package chat_handler

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

type ChatHandler struct {
	sessionClient session_client.AuthCheckerClient
	chatsUsecase  usecase_chats.Usecase
	bh.BaseHandler
}

func NewChatHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucChats usecase_chats.Usecase) *ChatHandler {
	h := &ChatHandler{
		sessionClient: sManager,
		chatsUsecase:  ucChats,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	h.AddMethod(http.MethodGet, h.POST,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *ChatHandler) GET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	u, err := h.chatsUsecase.GetChats(userID.(string))
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Log(r).Debugf("get user %s", u)
	h.Respond(w, r, http.StatusOK, models_http.ToResponseChats(u))
}

func (h *ChatHandler) POST(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	compan, status := h.GetStringFromQueries(w, r, "with")
	if status == bh.EmptyQuery {
		h.Log(r).Warnf("not found with in query")
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidQueries)
		return
	}

	us, err := h.chatsUsecase.Create(userID.(string), compan)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPOST)
		return
	}

	h.Respond(w, r, http.StatusCreated, models_http.ToResponseChat(*us))
}

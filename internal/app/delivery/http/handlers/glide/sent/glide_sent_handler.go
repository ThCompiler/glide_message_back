package glide_sent_handler

import (
	models_http "glide/internal/app/delivery/http/models"
	ucGlideMessage "glide/internal/app/usecase/glidemessage"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"glide/internal/pkg/handler/handler_errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GlideSentHandler struct {
	sessionClient       session_client.AuthCheckerClient
	glideMessageUsecase ucGlideMessage.Usecase
	bh.BaseHandler
}

func NewGlideSentHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucGlideMessage ucGlideMessage.Usecase) *GlideSentHandler {
	h := &GlideSentHandler{
		sessionClient:       sManager,
		glideMessageUsecase: ucGlideMessage,
		BaseHandler:         *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *GlideSentHandler) GET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	u, err := h.glideMessageUsecase.GetSent(userID.(string))
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Respond(w, r, http.StatusOK, models_http.ToResponseGlideMessages(u))
}

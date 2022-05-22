package glide_id_info_handler

import (
	"github.com/gorilla/mux"
	"glide/internal/app/delivery/http/handlers"
	models_http "glide/internal/app/delivery/http/models"
	ucGlideMessage "glide/internal/app/usecase/glidemessage"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"glide/internal/pkg/handler/handler_errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GlideIdInfoHandler struct {
	sessionClient       session_client.AuthCheckerClient
	glideMessageUsecase ucGlideMessage.Usecase
	bh.BaseHandler
}

func NewGlideIdInfoHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucGlideMessage ucGlideMessage.Usecase) *GlideIdInfoHandler {
	h := &GlideIdInfoHandler{
		sessionClient:       sManager,
		glideMessageUsecase: ucGlideMessage,
		BaseHandler:         *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *GlideIdInfoHandler) GET(w http.ResponseWriter, r *http.Request) {
	msgId, code, err := h.GetInt64FromParam(w, r, handlers.GlideMessageId)
	if err != nil {
		h.Error(w, r, code, err)
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	u, err := h.glideMessageUsecase.Get(msgId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Respond(w, r, http.StatusOK, models_http.ToResponseGlideMessage(*u))
}

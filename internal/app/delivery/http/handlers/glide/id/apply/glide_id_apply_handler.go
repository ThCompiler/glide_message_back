package glide_id_apply_handler

import (
	"github.com/gorilla/mux"
	"glide/internal/app/delivery/http/handlers"
	models_http "glide/internal/app/delivery/http/models"
	"glide/internal/app/middleware"
	ucGlideMessage "glide/internal/app/usecase/glidemessage"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"glide/internal/pkg/handler/handler_errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GlideIdApplyHandler struct {
	sessionClient       session_client.AuthCheckerClient
	glideMessageUsecase ucGlideMessage.Usecase
	bh.BaseHandler
}

func NewGlideIdApplyHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucGlideMessage ucGlideMessage.Usecase) *GlideIdApplyHandler {
	h := &GlideIdApplyHandler{
		sessionClient:       sManager,
		glideMessageUsecase: ucGlideMessage,
		BaseHandler:         *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodPut, h.PUT,
		middleware.NewGlideMessageMiddleware(log, ucGlideMessage).CheckCorrectUserFunc,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *GlideIdApplyHandler) PUT(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

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

	chat, err := h.glideMessageUsecase.Apply(h.Log(r), userID.(string), msgId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	h.Respond(w, r, http.StatusCreated, models_http.ToResponseChat(*chat))
}

package glide_id_picture_handler

import (
	"github.com/gorilla/mux"
	"glide/internal/app/delivery/http/handlers"
	"glide/internal/app/middleware"
	ucGlideMessage "glide/internal/app/usecase/glidemessage"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"glide/internal/pkg/handler/handler_errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GlideIdPictureHandler struct {
	sessionClient       session_client.AuthCheckerClient
	glideMessageUsecase ucGlideMessage.Usecase
	bh.BaseHandler
}

func NewGlideIdPictureHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucGlideMessage ucGlideMessage.Usecase) *GlideIdPictureHandler {
	h := &GlideIdPictureHandler{
		sessionClient:       sManager,
		glideMessageUsecase: ucGlideMessage,
		BaseHandler:         *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodPut, h.PUT,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
		middleware.NewGlideMessageMiddleware(log, ucGlideMessage).CheckCorrectAuthorFunc,
	)

	return h
}

func (h *GlideIdPictureHandler) PUT(w http.ResponseWriter, r *http.Request) {
	msgId, code, err := h.GetInt64FromParam(w, r, handlers.GlideMessageId)
	if err != nil {
		h.Error(w, r, code, err)
		return
	}

	file, filename, code, err := h.GetFilesFromRequest(w, r, handlers.MAX_UPLOAD_SIZE,
		"avatar", []string{"image/png", "image/jpeg", "image/jpg"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	err = h.glideMessageUsecase.UpdatePicture(msgId, file, filename)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

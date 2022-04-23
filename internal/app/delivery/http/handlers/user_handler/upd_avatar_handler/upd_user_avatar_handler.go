package upd_user_avatar_handler

import (
	"glide/internal/app/delivery/http/handlers"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	usecase_user "glide/internal/app/usecase/user"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UpdateUserAvatarHandler struct {
	sessionClient session_client.AuthCheckerClient
	userUsecase   usecase_user.Usecase
	bh.BaseHandler
}

func NewUpdateUserAvatarHandler(log *logrus.Logger,
	sessionClient session_client.AuthCheckerClient, userUsecase usecase_user.Usecase) *UpdateUserAvatarHandler {
	h := &UpdateUserAvatarHandler{
		sessionClient: sessionClient,
		userUsecase:   userUsecase,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMiddleware(session_middleware.NewSessionMiddleware(h.sessionClient, log).Check)
	h.AddMethod(http.MethodPut, h.PUT)
	return h
}

func (h *UpdateUserAvatarHandler) PUT(w http.ResponseWriter, r *http.Request) {
	file, filename, code, err := h.GerFilesFromRequest(w, r, handlers.MAX_UPLOAD_SIZE,
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

	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	err = h.userUsecase.UpdateAvatar(file, filename, userID.(string))
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

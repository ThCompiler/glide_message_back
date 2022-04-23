package upd_avatar_creator_handler

import (
	"net/http"
	csrf_middleware "glide/internal/app/csrf/middleware"
	repository_jwt "glide/internal/app/csrf/repository/jwt"
	usecase_csrf "glide/internal/app/csrf/usecase"
	"glide/internal/app/delivery/http/handlers"
	bh "glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/middleware"
	usecase_creator "glide/internal/app/usecase/creator"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UpdateAvatarCreatorHandler struct {
	sessionClient  session_client.AuthCheckerClient
	creatorUsecase usecase_creator.Usecase
	bh.BaseHandler
}

func NewUpdateAvatarHandler(log *logrus.Logger,
	sessionClient session_client.AuthCheckerClient, creatorUsecase usecase_creator.Usecase) *UpdateAvatarCreatorHandler {
	h := &UpdateAvatarCreatorHandler{
		sessionClient:  sessionClient,
		creatorUsecase: creatorUsecase,
		BaseHandler:    *bh.NewBaseHandler(log),
	}
	h.AddMiddleware(session_middleware.NewSessionMiddleware(h.sessionClient, log).Check,
		middleware.NewCreatorsMiddleware(log).CheckAllowUser)
	h.AddMethod(http.MethodPut, h.PUT,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// PUT AvatarChange
// @Summary set new creator avatar
// @tags creators
// @Accept  image/png, image/jpeg, image/jpg
// @Param avatar formData file true "Avatar file with ext jpeg/png, image/jpeg, image/jpg, max size 4 MB"
// @Success 200 "successfully upload avatar"
// @Failure 400 {object} http_models.ErrResponse "size of file very big", "please upload a JPEG, JPG or PNG files", "invalid form field name"
// @Failure 403 {object} http_models.ErrResponse "csrf token is invalid, get new token"
// @Failure 422 {object} http_models.ErrResponse "this creator id not know"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Router /creators/{creator_id:}/update/avatar [PUT]
func (h *UpdateAvatarCreatorHandler) PUT(w http.ResponseWriter, r *http.Request) {
	file, filename, code, err := h.GerFilesFromRequest(w, r, handlers.MAX_UPLOAD_SIZE,
		"avatar", []string{"image/png", "image/jpeg", "image/jpg"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	creatorId, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	err = h.creatorUsecase.UpdateAvatar(file, filename, creatorId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

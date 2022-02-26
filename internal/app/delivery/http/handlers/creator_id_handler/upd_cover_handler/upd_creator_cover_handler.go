package upd_cover_creator_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	"patreon/internal/app/delivery/http/handlers"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	usecase_creator "patreon/internal/app/usecase/creator"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UpdateCoverCreatorHandler struct {
	sessionClient  session_client.AuthCheckerClient
	creatorUsecase usecase_creator.Usecase
	bh.BaseHandler
}

func NewUpdateCoverHandler(log *logrus.Logger,
	sClient session_client.AuthCheckerClient, creatorUsecase usecase_creator.Usecase) *UpdateCoverCreatorHandler {
	h := &UpdateCoverCreatorHandler{
		sessionClient:  sClient,
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

// PUT CoverChange
// @Summary set new creator cover
// @tags creators
// @Accept  image/png, image/jpeg, image/jpg
// @Param cover formData file true "Cover file with ext jpeg/png, image/jpeg, image/jpg, max size 4 MB"
// @Success 200 "successfully upload cover"
// @Failure 400 {object} http_models.ErrResponse "size of file very big", "invalid form field name", "please upload a JPEG, JPG or PNG files"
// @Failure 403 {object} http_models.ErrResponse "csrf token is invalid, get new token"
// @Failure 422 {object} http_models.ErrResponse "this creator id not know"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Router /creators/{creator_id:}/update/cover [PUT]
func (h *UpdateCoverCreatorHandler) PUT(w http.ResponseWriter, r *http.Request) {
	file, filename, code, err := h.GerFilesFromRequest(w, r, handlers.MAX_UPLOAD_SIZE,
		"cover", []string{"image/png", "image/jpeg", "image/jpg"})
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

	err = h.creatorUsecase.UpdateCover(file, filename, creatorId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

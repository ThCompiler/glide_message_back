package aw_id_handler

import (
	"net/http"
	csrf_middleware "glide/internal/app/csrf/middleware"
	repository_jwt "glide/internal/app/csrf/repository/jwt"
	usecase_csrf "glide/internal/app/csrf/usecase"
	bh "glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/middleware"
	useAwards "glide/internal/app/usecase/awards"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AwardsIdHandler struct {
	awardsUsecase useAwards.Usecase
	bh.BaseHandler
}

func NewAwardsIdHandler(log *logrus.Logger,
	ucAwards useAwards.Usecase, sClient session_client.AuthCheckerClient) *AwardsIdHandler {
	h := &AwardsIdHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		awardsUsecase: ucAwards,
	}

	h.AddMethod(http.MethodDelete, h.DELETE, session_middleware.NewSessionMiddleware(sClient, log).CheckFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)

	return h
}

// DELETE Awards
// @Summary delete current awards
// @tags awards
// @Description delete current awards from current creator
// @Produce json
// @Success 200
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 404 {object} http_models.ErrResponse "award with this id not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this awards not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/awards/{:award_id} [DELETE]
func (h *AwardsIdHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	awardsId, ok := h.GetInt64FromParam(w, r, "award_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	err := h.awardsUsecase.Delete(awardsId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsDELETE)
		return
	}

	w.WriteHeader(http.StatusOK)
}

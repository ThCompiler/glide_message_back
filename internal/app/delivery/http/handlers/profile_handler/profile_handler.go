package profile_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	models_http "patreon/internal/app/delivery/http/models"
	usecase_user "patreon/internal/app/usecase/user"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/sirupsen/logrus"
)

type ProfileHandler struct {
	sessionClient session_client.AuthCheckerClient
	userUsecase   usecase_user.Usecase
	bh.BaseHandler
}

func NewProfileHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucUser usecase_user.Usecase) *ProfileHandler {
	h := &ProfileHandler{
		sessionClient: sManager,
		userUsecase:   ucUser,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)
	return h
}

// GET Profile
// @Summary get information from user for user
// @tags user
// @Description get nickname and avatar for user
// @Accept  json
// @Produce json
// @Success 201 {object} http_models.ProfileResponse "Successfully get user"
// @Failure 404 {object} http_models.ErrResponse "user with this id not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 401 "user are not authorized"
// @Router /user [GET]
func (h *ProfileHandler) GET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	u, err := h.userUsecase.GetProfile(userID.(int64))
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}

	h.Log(r).Debugf("get user %s", u)
	h.Respond(w, r, http.StatusOK, models_http.ToRProfileResponse(*u))
}

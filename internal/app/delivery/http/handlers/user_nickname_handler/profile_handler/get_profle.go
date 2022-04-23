package user_nickname_profile_handler

import (
	"github.com/gorilla/mux"
	"glide/internal/app/delivery/http/handlers"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	models_http "glide/internal/app/delivery/http/models"
	usecase_user "glide/internal/app/usecase/user"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GetProfileHandler struct {
	sessionClient session_client.AuthCheckerClient
	userUsecase   usecase_user.Usecase
	bh.BaseHandler
}

func NewGetProfileHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucUser usecase_user.Usecase) *GetProfileHandler {
	h := &GetProfileHandler{
		sessionClient: sManager,
		userUsecase:   ucUser,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *GetProfileHandler) GET(w http.ResponseWriter, r *http.Request) {
	nickname, status := h.GetStringFromParam(w, r, handlers.UserNickname)

	if status == bh.EmptyQuery {
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	u, err := h.userUsecase.GetProfile(nickname)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Log(r).Debugf("get user %s", u)
	h.Respond(w, r, http.StatusOK, models_http.ToProfileResponse(*u))
}

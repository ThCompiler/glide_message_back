package login_handler

import (
	"context"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/delivery/http/models"
	usecase_user "glide/internal/app/usecase/user"
	sessionclient "glide/internal/microservices/auth/delivery/grpc/client"
	"glide/internal/microservices/auth/sessions/middleware"
	"glide/internal/microservices/auth/sessions/sessions_manager"
	bh "glide/internal/pkg/handler"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoginHandler struct {
	sessionClient sessionclient.AuthCheckerClient
	userUsecase   usecase_user.Usecase
	bh.BaseHandler
}

func NewLoginHandler(log *logrus.Logger, sClient sessionclient.AuthCheckerClient,
	ucUser usecase_user.Usecase) *LoginHandler {
	h := &LoginHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		sessionClient: sClient,
		userUsecase:   ucUser,
	}
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionClient, log).CheckNotAuthorized)
	h.AddMethod(http.MethodPost, h.POST)
	return h
}

func (h *LoginHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestLogin{}
	err := h.GetRequestBody(r, req)
	if err != nil || len(req.Password) == 0 || len(req.Login) == 0 {
		h.Log(r).Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	nickname, err := h.userUsecase.Check(req.Login, req.Password)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrors)
		return
	}

	res, err := h.sessionClient.Create(context.Background(), nickname)
	if err != nil || res.UserID != nickname {
		h.Log(r).Errorf("Error create session %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ErrorCreateSession)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    res.UniqID,
		Path:     "/",
		Expires:  time.Now().Add(sessions_manager.ExpiredCookiesTime),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

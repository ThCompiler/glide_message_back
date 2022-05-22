package logout_handler

import (
	"context"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"glide/internal/pkg/handler/handler_errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LogoutHandler struct {
	sessionClient session_client.AuthCheckerClient
	bh.BaseHandler
}

func NewLogoutHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient) *LogoutHandler {
	h := &LogoutHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		sessionClient: sManager,
	}
	h.AddMethod(http.MethodPost, h.POST,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *LogoutHandler) POST(w http.ResponseWriter, r *http.Request) {
	uniqID := r.Context().Value("session_id")
	if uniqID == nil {
		h.Log(r).Error("can not get session_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	h.Log(r).Debugf("Logout session: %s", uniqID)

	err := h.sessionClient.Delete(context.Background(), uniqID.(string))
	if err != nil {
		h.Log(r).Errorf("can not delete session %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.DeleteCookieFail)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    uniqID.(string),
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, -1),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

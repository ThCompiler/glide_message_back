package csrf_handler

import (
	"net/http"
	usecase_csrf "glide/internal/app/csrf/usecase"
	bh "glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/delivery/http/models"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"

	"github.com/sirupsen/logrus"
)

type CsrfHandler struct {
	csrfUsecase   usecase_csrf.Usecase
	sessionClient session_client.AuthCheckerClient
	bh.BaseHandler
}

func NewCsrfHandler(log *logrus.Logger, sClient session_client.AuthCheckerClient,
	uc usecase_csrf.Usecase) *CsrfHandler {
	h := &CsrfHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		sessionClient: sClient,
		csrfUsecase:   uc,
	}
	h.AddMethod(http.MethodGet, h.GET, session_middleware.NewSessionMiddleware(sClient, log).CheckFunc)
	return h
}

// GET CSRF Token
// @tags utilities
// @Summary get CSRF Token
// @Description generate usecase token and return to client
// @Produce json
// @Success 200 {object} http_models.TokenResponse
// @Failure 500 {object} http_models.ErrResponse "server error"
// @Failure 401 "user are not authorized"
// @Router /token [GET]
func (h *CsrfHandler) GET(w http.ResponseWriter, r *http.Request) {
	sessionId, ok := r.Context().Value("session_id").(string)
	if !ok {
		h.Log(r).Error("invalid conversation session_id from context to string")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Error("invalid conversation userId from context to int64")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	token, err := h.csrfUsecase.Create(sessionId, userId)
	if err != nil {
		h.Log(r).Error("can not create token")
		h.UsecaseError(w, r, err, codeByErrors)
		return
	}
	h.Log(r).Debugf("get token %v", token)
	h.Respond(w, r, http.StatusOK, http_models.TokenResponse{Token: token})
}

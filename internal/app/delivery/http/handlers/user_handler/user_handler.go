package user_handler

import (
	"github.com/microcosm-cc/bluemonday"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/delivery/http/models"
	models_http "glide/internal/app/delivery/http/models"
	usecase_user "glide/internal/app/usecase/user"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"net/http"

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

	h.AddMethod(http.MethodPost, h.POST,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckNotAuthorizedFunc,
	)

	h.AddMethod(http.MethodPut, h.PUT,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *ProfileHandler) GET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	u, err := h.userUsecase.GetProfile(userID.(string))
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorGET)
		return
	}

	h.Log(r).Debugf("get user %s", u)
	h.Respond(w, r, http.StatusOK, models_http.ToProfileResponse(*u))
}

func (h *ProfileHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestRegistration{}

	err := h.GetRequestBody(r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	u := req.ToUser()

	us, err := h.userUsecase.Create(u)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPOST)
		return
	}

	u.MakeEmptyPassword()
	h.Respond(w, r, http.StatusCreated, models_http.ToProfileResponse(*us))
}

func (h *ProfileHandler) PUT(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestUserUpdate{}

	err := h.GetRequestBody(r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	u := req.ToUser()

	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	u.Nickname = userID.(string)

	us, err := h.userUsecase.Update(u)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	h.Respond(w, r, http.StatusCreated, models_http.ToProfileResponse(*us))
}

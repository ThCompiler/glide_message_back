package glide_create_handler

import (
	"github.com/microcosm-cc/bluemonday"
	"glide/internal/app"
	"glide/internal/app/delivery/http/models"
	ucGlideMessage "glide/internal/app/usecase/glidemessage"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"
	bh "glide/internal/pkg/handler"
	"glide/internal/pkg/handler/handler_errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GlideIdApplyHandler struct {
	sessionClient       session_client.AuthCheckerClient
	glideMessageUsecase ucGlideMessage.Usecase
	bh.BaseHandler
}

func NewGlideIdApplyHandler(log *logrus.Logger,
	sManager session_client.AuthCheckerClient, ucGlideMessage ucGlideMessage.Usecase) *GlideIdApplyHandler {
	h := &GlideIdApplyHandler{
		sessionClient:       sManager,
		glideMessageUsecase: ucGlideMessage,
		BaseHandler:         *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodPost, h.POST,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc,
	)

	return h
}

func (h *GlideIdApplyHandler) POST(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	languages, _ := h.GetArrayStringFromQueries(w, r, "languages")
	countries, _ := h.GetArrayStringFromQueries(w, r, "countries")
	age, code, err := h.GetInt64FromQueries(w, r, "age")

	if age == app.InvalidInt {
		h.Error(w, r, code, err)
	}

	req := &http_models.RequestGlideMessage{}

	err = h.GetRequestBody(r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	req.Author = userID.(string)
	msg, err := h.glideMessageUsecase.Create(h.Log(r), req.ToGlideMessage(), languages, countries, age)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPOST)
		return
	}

	h.Respond(w, r, http.StatusCreated, http_models.ToResponseGlideMessage(*msg))
}

package subscribe_handler

import (
	"net/http"
	bh "glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/delivery/http/models"
	usecase_subscribers "glide/internal/app/usecase/subscribers"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	"glide/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type SubscribeHandler struct {
	sessionClient     session_client.AuthCheckerClient
	subscriberUsecase usecase_subscribers.Usecase
	bh.BaseHandler
}

func NewSubscribeHandler(log *logrus.Logger, sClient session_client.AuthCheckerClient,
	ucSubscribers usecase_subscribers.Usecase) *SubscribeHandler {
	h := &SubscribeHandler{
		BaseHandler:       *bh.NewBaseHandler(log),
		subscriberUsecase: ucSubscribers,
		sessionClient:     sClient,
	}
	h.AddMethod(http.MethodGet, h.GET, middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc)
	return h
}

// GET Subscribers
// @Summary subscribers of the creator
// @tags creators
// @Description get subscribers of the creators with id = creator_id
// @Produce json
// @Param creator_id path int true "creator_id"
// @Success 200 {object} http_models.SubscribersCreatorResponse "Successfully get creator subscribers with creator id = creator_id"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 500 {object} http_models.ErrResponse "server error", "can not do bd operation"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/subscribers [GET]
func (h *SubscribeHandler) GET(w http.ResponseWriter, r *http.Request) {
	creatorID, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		h.Log(r).Warnf("invalid creator_id %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parameters %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	subscribers, err := h.subscriberUsecase.GetSubscribers(creatorID)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}
	res := http_models.ToSubscribersCreatorResponse(subscribers)
	h.Log(r).Debugf("get users %v", subscribers)
	h.Respond(w, r, http.StatusOK, res)
}

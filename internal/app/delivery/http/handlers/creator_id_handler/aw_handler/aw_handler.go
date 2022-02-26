package aw_handler

import (
	"image/color"
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	db_models "patreon/internal/app/models"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"

	useAwards "patreon/internal/app/usecase/awards"
)

type AwardsHandler struct {
	awardsUsecase useAwards.Usecase
	bh.BaseHandler
}

func NewAwardsHandler(log *logrus.Logger,
	ucAwards useAwards.Usecase, sClient session_client.AuthCheckerClient) *AwardsHandler {
	h := &AwardsHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		awardsUsecase: ucAwards,
	}
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMethod(http.MethodPost, h.POST, session_middleware.NewSessionMiddleware(sClient, log).CheckFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)

	return h
}

// GET Awards
// @Summary get list of awards of some creator
// @tags awards
// @Description get list of awards which belongs the creator
// @Produce json
// @Success 201 {object} http_models.ResponseAwards
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Router /creators/{:creator_id}/awards [GET]
func (h *AwardsHandler) GET(w http.ResponseWriter, r *http.Request) {
	idInt, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	awards, err := h.awardsUsecase.GetAwards(idInt)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondAwards := make([]http_models.ResponseAward, len(awards))
	for i, aw := range awards {
		respondAwards[i] = http_models.ToResponseAward(aw)
	}

	h.Log(r).Debugf("get creators %v", respondAwards)
	h.Respond(w, r, http.StatusOK, http_models.ResponseAwards{Awards: respondAwards})
}

// POST Create Awards
// @Summary create awards
// @tags awards
// @Description create awards to creator with id from path
// @Param award body http_models.RequestAwards true "Request body for awards"
// @Produce json
// @Success 201 {object} http_models.IdResponse "id awards"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 422 {object} http_models.ErrResponse "empty name in request", "incorrect value of price", "invalid body in request"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 409 {object} http_models.ErrResponse "awards with this price already exists", "awards with this name already exists"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator",  "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/awards [POST]
func (h *AwardsHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestAwards{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	idInt, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}
	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	aw := &db_models.Award{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Color:       color.RGBA{R: req.Color.R, G: req.Color.G, B: req.Color.B, A: req.Color.A},
		CreatorId:   idInt,
	}

	awardsId, err := h.awardsUsecase.Create(aw)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}

	h.Respond(w, r, http.StatusCreated, &http_models.IdResponse{ID: awardsId})
}

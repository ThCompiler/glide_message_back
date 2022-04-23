package creator_id_handler

import (
	"net/http"
	"glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/delivery/http/models"
	usecase_creator "glide/internal/app/usecase/creator"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	"glide/internal/microservices/auth/sessions/middleware"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type CreatorIdHandler struct {
	sessionClient  session_client.AuthCheckerClient
	creatorUsecase usecase_creator.Usecase
	base_handler.BaseHandler
}

func NewCreatorIdHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient, ucCreator usecase_creator.Usecase) *CreatorIdHandler {
	h := &CreatorIdHandler{
		BaseHandler:    *base_handler.NewBaseHandler(log),
		sessionClient:  sManager,
		creatorUsecase: ucCreator,
	}
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionClient, log).AddUserId)
	h.AddMethod(http.MethodGet, h.GET)

	return h
}

// GET Creator
// @Summary get creator
// @Description get creator with id from path
// @Produce json
// @tags creators
// @Param creator_id path int true "Get creator with id"
// @Success 200 {object} http_models.ResponseCreatorWithAwards
// @Failure 404 {object} http_models.ErrResponse "user not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Router /creators/{creator_id:} [GET]
func (s *CreatorIdHandler) GET(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		userId = usecase_creator.NoUser
	}

	creatorId, ok := s.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		s.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		s.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	creator, err := s.creatorUsecase.GetCreator(creatorId, userId)
	if err != nil {
		s.UsecaseError(w, r, err, codesByErrors)
		return
	}

	s.Log(r).Debugf("get creator %v with id %v", creator, creatorId)
	s.Respond(w, r, http.StatusOK, http_models.ToResponseCreatorWithAwards(*creator))
}

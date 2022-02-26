package register_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/models"
	usecase_user "patreon/internal/app/usecase/user"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"

	"github.com/microcosm-cc/bluemonday"

	"github.com/sirupsen/logrus"
)

type RegisterHandler struct {
	sessionClient session_client.AuthCheckerClient
	userUsecase   usecase_user.Usecase
	bh.BaseHandler
}

func NewRegisterHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient,
	ucUser usecase_user.Usecase) *RegisterHandler {
	h := &RegisterHandler{
		sessionClient: sManager,
		userUsecase:   ucUser,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodPost, h.POST)
	return h
}

// POST Registration
// @Summary create new user
// @tags user
// @Description create new account and get cookies
// @Accept  json
// @Produce json
// @Param register_info body http_models.RequestRegistration true "Request body for user registration"
// @Success 201 {object} http_models.IdResponse "Create user successfully"
// @Failure 409 {object} http_models.ErrResponse "user already exist"
// @Failure 422 {object} http_models.ErrResponse "invalid body in request", "nickname already exist", "incorrect email or password", "incorrect nickname"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation"
// @Failure 418 "User are authorized"
// @Router /register [POST]
func (h *RegisterHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestRegistration{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	u := &models.User{
		Login:    req.Login,
		Password: req.Password,
		Nickname: req.Nickname,
	}

	id, err := h.userUsecase.Create(u)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}

	u.MakeEmptyPassword()
	h.Respond(w, r, http.StatusCreated, http_models.IdResponse{ID: id})
}

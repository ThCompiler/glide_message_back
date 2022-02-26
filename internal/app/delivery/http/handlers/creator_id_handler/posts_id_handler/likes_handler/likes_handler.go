package likes_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	"patreon/internal/app/models"
	useLikes "patreon/internal/app/usecase/likes"
	usePosts "patreon/internal/app/usecase/posts"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type LikesHandler struct {
	likesUsecase useLikes.Usecase
	bh.BaseHandler
}

func NewLikesHandler(log *logrus.Logger,
	ucLikes useLikes.Usecase, ucPosts usePosts.Usecase, sClient session_client.AuthCheckerClient) *LikesHandler {
	h := &LikesHandler{
		BaseHandler:  *bh.NewBaseHandler(log),
		likesUsecase: ucLikes,
	}
	postsMiddleware := middleware.NewPostsMiddleware(log, ucPosts)
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)
	h.AddMiddleware(sessionMiddleware.Check, postsMiddleware.CheckCorrectPost)

	h.AddMethod(http.MethodDelete, h.DELETE,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)

	h.AddMethod(http.MethodPut, h.PUT,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)
	return h
}

// DELETE Likes
// @Summary deletes like from the post and return new count of likes
// @tags posts
// @Description deletes like form post id in url
// @Produce json
// @Success 200 {object} http_models.ResponseLike "current count of likes on post"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 409 {object} http_models.ErrResponse "this user not have like for this post"
// @Failure 403 {object} http_models.ErrResponse "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/like [DELETE]
func (h *LikesHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	var postsId, userId int64
	var ok bool
	postsId, ok = h.GetInt64FromParam(w, r, "post_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	userId, ok = r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	res, err := h.likesUsecase.Delete(postsId, userId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsDELETE)
		return
	}
	h.Respond(w, r, http.StatusOK, http_models.ResponseLike{Likes: res})
}

// PUT Likes
// @Summary add like on the post
// @tags posts
// @Description add like on the post with id = post_id and return new count of likes
// @Produce json
// @Success 200 {object} http_models.ResponseLike "current count of likes on post"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 409 {object} http_models.ErrResponse "this user already add like for this post"
// @Failure 403 {object} http_models.ErrResponse "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/like [PUT]
func (h *LikesHandler) PUT(w http.ResponseWriter, r *http.Request) {
	var postsId, userId int64
	var ok bool
	postsId, ok = h.GetInt64FromParam(w, r, "post_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	userId, ok = r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	res, err := h.likesUsecase.Add(&models.Like{PostId: postsId, UserId: userId})
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPUT)
		return
	}

	h.Respond(w, r, http.StatusOK, http_models.ResponseLike{Likes: res})
}

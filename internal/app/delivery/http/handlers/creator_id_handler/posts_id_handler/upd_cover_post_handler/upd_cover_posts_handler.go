package upl_cover_posts_handler

import (
	"net/http"
	csrf_middleware "glide/internal/app/csrf/middleware"
	repository_jwt "glide/internal/app/csrf/repository/jwt"
	usecase_csrf "glide/internal/app/csrf/usecase"
	"glide/internal/app/delivery/http/handlers"
	bh "glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/middleware"
	usePosts "glide/internal/app/usecase/posts"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type PostsUpdateCoverHandler struct {
	postsUsecase usePosts.Usecase
	bh.BaseHandler
}

func NewPostsUpdateCoverHandler(log *logrus.Logger,
	ucPosts usePosts.Usecase, sClient session_client.AuthCheckerClient) *PostsUpdateCoverHandler {
	h := &PostsUpdateCoverHandler{
		BaseHandler:  *bh.NewBaseHandler(log),
		postsUsecase: ucPosts,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)
	h.AddMiddleware(sessionMiddleware.Check, middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)

	h.AddMethod(http.MethodPut, h.PUT,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)
	return h
}

// PUT CoverUpdate
// @Summary set new post cover
// @tags posts
// @Accept  image/png, image/jpeg, image/jpg
// @Param cover formData file true "cover file with ext jpeg/png, image/jpeg, image/jpg, max size 4 MB"
// @Success 200 "successfully upload cover"
// @Failure 400 {object} http_models.ErrResponse "size of file very big", "please upload a JPEG, JPG or PNG files", "invalid form field name for load file"
// @Failure 404 {object} http_models.ErrResponse "post with this id not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/update/cover [PUT]
func (h *PostsUpdateCoverHandler) PUT(w http.ResponseWriter, r *http.Request) {
	var postId int64
	var ok bool
	if postId, ok = h.GetInt64FromParam(w, r, "post_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	file, filename, code, err := h.GerFilesFromRequest(w, r, handlers.MAX_UPLOAD_SIZE,
		"cover", []string{"image/png", "image/jpeg", "image/jpg"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	err = h.postsUsecase.LoadCover(file, filename, postId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}

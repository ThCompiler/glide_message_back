package upl_video_attach_handler

import (
	"net/http"
	csrf_middleware "glide/internal/app/csrf/middleware"
	repository_jwt "glide/internal/app/csrf/repository/jwt"
	usecase_csrf "glide/internal/app/csrf/usecase"
	"glide/internal/app/delivery/http/handlers"
	bh "glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/delivery/http/models"
	"glide/internal/app/middleware"
	useAttaches "glide/internal/app/usecase/attaches"
	usePosts "glide/internal/app/usecase/posts"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type PostsUploadVideoHandler struct {
	attachesUsecase useAttaches.Usecase
	bh.BaseHandler
}

func NewPostsUploadVideoHandler(log *logrus.Logger,
	ucAttaches useAttaches.Usecase,
	ucPosts usePosts.Usecase,
	sClient session_client.AuthCheckerClient) *PostsUploadVideoHandler {
	h := &PostsUploadVideoHandler{
		BaseHandler:     *bh.NewBaseHandler(log),
		attachesUsecase: ucAttaches,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)
	h.AddMiddleware(sessionMiddleware.Check, middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)

	h.AddMethod(http.MethodPost, h.POST,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// POST add video to post
// @Summary add video to post
// @tags attaches
// @Accept 	video/mpeg4-generic, video/mp4,  video/mpeg
// @Param video formData file true "image file with ext video/mpeg4-generic, video/mpeg, video/mp4, max size 80 MB"
// @Success 201 {object} http_models.IdResponse "id attaches"
// @Failure 400 {object} http_models.ErrResponse "size of file very big", "please upload some types", "invalid form field name for load file"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} http_models.ErrResponse "invalid data type", "this post id not know"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/attaches/video [POST]
func (h *PostsUploadVideoHandler) POST(w http.ResponseWriter, r *http.Request) {
	var attachId, postId int64
	var ok bool

	if postId, ok = h.GetInt64FromParam(w, r, "post_id"); !ok {
		return
	}
	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	file, filename, code, err := h.GerFilesFromRequest(w, r, handlers.MAX_UPLOAD_VIDEO_SIZE,
		"video", []string{"video/mpeg4-generic", "video/mpeg", "video/mp4"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	attachId, err = h.attachesUsecase.LoadVideo(file, filename, postId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	h.Respond(w, r, http.StatusCreated, &http_models.IdResponse{ID: attachId})
}

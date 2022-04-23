package upd_audio_attach_handler

import (
	"net/http"
	csrf_middleware "glide/internal/app/csrf/middleware"
	repository_jwt "glide/internal/app/csrf/repository/jwt"
	usecase_csrf "glide/internal/app/csrf/usecase"
	"glide/internal/app/delivery/http/handlers"
	bh "glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/middleware"
	useAttaches "glide/internal/app/usecase/attaches"
	usePosts "glide/internal/app/usecase/posts"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	session_middleware "glide/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type AttachUploadAudioHandler struct {
	attachesUsecase useAttaches.Usecase
	bh.BaseHandler
}

func NewAttachUploadAudioHandler(
	log *logrus.Logger,
	ucAttaches useAttaches.Usecase,
	ucPosts usePosts.Usecase,
	sClient session_client.AuthCheckerClient) *AttachUploadAudioHandler {
	h := &AttachUploadAudioHandler{
		BaseHandler:     *bh.NewBaseHandler(log),
		attachesUsecase: ucAttaches,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)
	h.AddMiddleware(sessionMiddleware.Check,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.
			NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfToken,
		middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost,
		middleware.NewAttachesMiddleware(log, ucAttaches).CheckCorrectAttach)

	h.AddMethod(http.MethodPut, h.PUT)

	return h
}

// PUT update audio to post
// @Summary update audio to post
// @tags attaches
// @Accept  audio/mp3, audio/mpeg, audio/mpeg3
// @Param audio formData file true "image file with ext audio/mp3, audio/mpeg, audio/mpeg3, max size 30 MB"
// @Success 200
// @Failure 400 {object} http_models.ErrResponse "size of file very big", "invalid form field name for load file", "please upload a some types"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} http_models.ErrResponse "invalid data type", "this post id not know"
// @Failure 404 {object} http_models.ErrResponse "attach with this id not found"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "this post not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts/{:post_id}/{:attach_id}/update/audio [PUT]
func (h *AttachUploadAudioHandler) PUT(w http.ResponseWriter, r *http.Request) {
	var attachId int64
	var ok bool

	if attachId, ok = h.GetInt64FromParam(w, r, "attach_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 3 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	file, filename, code, err := h.GerFilesFromRequest(w, r, handlers.MAX_UPLOAD_AUDIO_SIZE,
		"audio", []string{"audio/mp3", "audio/mpeg", "audio/mpeg3"})
	if err != nil {
		h.HandlerError(w, r, code, err)
		return
	}

	err = h.attachesUsecase.UpdateAudio(file, filename, attachId)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}

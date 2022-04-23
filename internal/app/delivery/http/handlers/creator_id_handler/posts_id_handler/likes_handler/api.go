package likes_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/repository"
	usecase_likes "glide/internal/app/usecase/likes"
)

var codesByErrorsDELETE = base_handler.CodeMap{
	usecase_likes.IncorrectDelLike: {
		http.StatusConflict, handler_errors.LikesAlreadyDel, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

var codesByErrorsPUT = base_handler.CodeMap{
	usecase_likes.IncorrectAddLike: {
		http.StatusConflict, handler_errors.LikesAlreadyExists, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

package glide_id_apply_handler

import (
	"glide/internal/app/repository"
	repository_postgresql "glide/internal/app/repository/chat/postgresql"
	"glide/internal/pkg/handler/handler_errors"
	"glide/internal/pkg/utilits/delivery"
	"net/http"

	"github.com/sirupsen/logrus"
)

var codeByErrorPUT = delivery.CodeMap{
	repository_postgresql.ChatAlreadyExists: {
		http.StatusConflict, handler_errors.ChatAlreadyExist, logrus.WarnLevel},
	repository.NotFound: {
		http.StatusNotFound, handler_errors.UserNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

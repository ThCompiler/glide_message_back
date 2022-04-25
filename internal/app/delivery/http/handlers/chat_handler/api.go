package chat_handler

import (
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/repository"
	repository_postgresql "glide/internal/app/repository/chat/postgresql"
	"glide/internal/pkg/utilits/delivery"
	"net/http"

	"github.com/sirupsen/logrus"
)

var codeByErrorGET = delivery.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.UserNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

var codeByErrorPOST = delivery.CodeMap{
	repository_postgresql.ChatAlreadyExists: {
		http.StatusConflict, handler_errors.ChatAlreadyExist, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

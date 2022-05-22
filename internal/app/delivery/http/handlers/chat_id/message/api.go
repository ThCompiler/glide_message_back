package chat_id_message_handler

import (
	"glide/internal/app/repository"
	repository_os "glide/internal/app/repository/files/os"
	usercase_chat "glide/internal/app/usecase/chats"
	"glide/internal/pkg/handler/handler_errors"
	"glide/internal/pkg/utilits"
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
	repository.NotFound: {
		http.StatusNotFound, handler_errors.ChatNotFound, logrus.WarnLevel},
	repository_os.ErrorCreate: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	repository_os.ErrorOpenFile: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	utilits.ConvertErr: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	usercase_chat.FileSystemError: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

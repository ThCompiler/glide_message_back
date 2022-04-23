package attaches_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	repository_postgresql "glide/internal/app/repository/attaches/postgresql"
)

var codesByErrorsPOST = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.AttachNotFound, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	models.IncorrectType: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectType, logrus.WarnLevel},
	models.IncorrectAttachId: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectIdAttach, logrus.WarnLevel},
	models.IncorrectLevel: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	repository_postgresql.UnknownDataFormat: {
		http.StatusInternalServerError, handler_errors.IncorrectType, logrus.ErrorLevel},
}

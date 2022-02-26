package login_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"

	"github.com/sirupsen/logrus"
)

var codesByErrors = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.UserNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	models.IncorrectEmailOrPassword: {
		http.StatusUnauthorized, handler_errors.IncorrectLoginOrPassword, logrus.InfoLevel},
}

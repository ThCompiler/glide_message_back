package login_handler

import (
	"glide/internal/pkg/handler/handler_errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	"glide/internal/pkg/utilits/delivery"
)

var codesByErrors = delivery.CodeMap{
	repository.NotFound: {
		http.StatusUnauthorized, handler_errors.IncorrectLoginOrPassword, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	models.IncorrectNicknameOrPassword: {
		http.StatusUnauthorized, handler_errors.IncorrectLoginOrPassword, logrus.InfoLevel},
}

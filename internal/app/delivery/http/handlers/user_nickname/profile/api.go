package user_nickname_profile_handler

import (
	"glide/internal/app/repository"
	"glide/internal/pkg/handler/handler_errors"
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

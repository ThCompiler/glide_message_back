package glide_id_redirect_handler

import (
	"glide/internal/app/repository"
	"glide/internal/pkg/handler/handler_errors"
	"glide/internal/pkg/utilits/delivery"
	"net/http"

	"github.com/sirupsen/logrus"
)

var codeByErrorPUT = delivery.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.UserNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

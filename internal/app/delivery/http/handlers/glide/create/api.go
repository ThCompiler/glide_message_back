package glide_create_handler

import (
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/repository"
	repository_glidemess "glide/internal/app/repository/glidemessage/postgresql"
	"glide/internal/pkg/utilits/delivery"
	"net/http"

	"github.com/sirupsen/logrus"
)

var codeByErrorPOST = delivery.CodeMap{
	repository_glidemess.IncorrectCountry: {
		http.StatusBadRequest, handler_errors.InvalidUserLanguage, logrus.WarnLevel},
	repository_glidemess.IncorrectLanguage: {
		http.StatusBadRequest, handler_errors.InvalidUserCounty, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

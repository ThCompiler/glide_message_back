package aw_id_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/repository"
)

var codesByErrorsDELETE = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.AwardNotFound, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

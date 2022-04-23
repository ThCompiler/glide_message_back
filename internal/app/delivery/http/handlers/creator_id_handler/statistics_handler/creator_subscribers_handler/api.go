package statistics_count_subscribers_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/repository"
	"glide/internal/app/usecase/statistics"
)

var codeByErrorGet = base_handler.CodeMap{
	statistics.CreatorDoesNotExists: {
		http.StatusNotFound, handler_errors.CreatorNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

package creator_payments_handler

import (
	"net/http"
	"glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/repository"

	"github.com/sirupsen/logrus"
)

var codeByErrorGET = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNoContent, handler_errors.CreatorPaymentsNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}

package glide_id_picture_handler

import (
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/repository"
	repository_os "glide/internal/app/repository/files/os"
	"glide/internal/pkg/utilits"
	"glide/internal/pkg/utilits/delivery"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var codeByError = delivery.CodeMap{
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, log.ErrorLevel},
	repository.NotFound: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectCreatorId, log.WarnLevel},
	repository_os.ErrorCreate: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	repository_os.ErrorOpenFile: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	utilits.ConvertErr: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	utilits.UnknownExtOfFileName: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
}

package upd_avatar_creator_handler

import (
	"glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/repository"
	"glide/internal/app/repository/files/os"
	"glide/internal/pkg/utils"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var codeByError = base_handler.CodeMap{
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, log.ErrorLevel},
	repository.NotFound: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectCreatorId, log.WarnLevel},
	repository_os.ErrorCreate: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	repository_os.ErrorCreate: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	utils.ConvertErr: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	utils.UnknownExtOfFileName: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
}

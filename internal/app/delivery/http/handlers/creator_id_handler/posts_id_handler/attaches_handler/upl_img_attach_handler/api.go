package upl_img_attach_handler

import (
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_os "patreon/internal/microservices/files/files/repository/files/os"

	repository_postgresql "patreon/internal/app/repository/attaches/postgresql"
	"glide/internal/pkg/utils"

	"github.com/sirupsen/logrus"
)

var codeByErrorPUT = base_handler.CodeMap{
	repository_postgresql.UnknownDataFormat: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectDataType, logrus.WarnLevel},
	models.InvalidType: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectDataType, logrus.WarnLevel},
	models.InvalidPostId: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectPostId, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	app.UnknownError: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	repository_os.ErrorCopyFile: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	repository_os.ErrorCreate: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	utils.ConvertErr: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	utils.UnknownExtOfFileName: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}

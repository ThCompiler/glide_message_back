package user_handler

import (
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	repository_user "glide/internal/app/repository/user/postgresql"
	usercase_user "glide/internal/app/usecase/user"
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

var codeByErrorPOST = delivery.CodeMap{
	models.IncorrectAge: {
		http.StatusBadRequest, handler_errors.InvalidUserAge, logrus.InfoLevel},
	models.IncorrectNicknameOrPassword: {
		http.StatusBadRequest, handler_errors.InvalidUserNickname, logrus.InfoLevel},
	models.EmptyPassword: {
		http.StatusBadRequest, handler_errors.InvalidBody, logrus.InfoLevel},
	repository_user.NicknameAlreadyExist: {
		http.StatusConflict, handler_errors.NicknameAlreadyExist, logrus.WarnLevel},
	repository_user.IncorrectLanguage: {
		http.StatusBadRequest, handler_errors.InvalidUserLanguage, logrus.WarnLevel},
	repository_user.IncorrectCounty: {
		http.StatusBadRequest, handler_errors.InvalidUserCounty, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	usercase_user.BadEncrypt: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

var codeByErrorPUT = delivery.CodeMap{
	models.IncorrectAge: {
		http.StatusBadRequest, handler_errors.InvalidUserAge, logrus.InfoLevel},
	repository_user.NicknameAlreadyExist: {
		http.StatusConflict, handler_errors.NicknameAlreadyExist, logrus.WarnLevel},
	repository_user.IncorrectLanguage: {
		http.StatusBadRequest, handler_errors.InvalidUserLanguage, logrus.WarnLevel},
	repository_user.IncorrectCounty: {
		http.StatusBadRequest, handler_errors.InvalidUserCounty, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

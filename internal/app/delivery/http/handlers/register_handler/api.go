package register_handler

import (
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	repository_user "glide/internal/app/repository/user/postgresql"
	useUser "glide/internal/app/usecase/user"
	"glide/internal/pkg/utilits/delivery"
	"net/http"

	"github.com/sirupsen/logrus"
)

var codeByError = delivery.CodeMap{
	models.IncorrectNickname: {
		http.StatusUnprocessableEntity, handler_errors.InvalidUserNickname, logrus.InfoLevel},
	models.EmptyPassword: {
		http.StatusUnprocessableEntity, handler_errors.InvalidBody, logrus.InfoLevel},
	repository_user.LoginAlreadyExist: {
		http.StatusConflict, handler_errors.UserAlreadyExist, logrus.InfoLevel},
	useUser.UserExist: {
		http.StatusConflict, handler_errors.UserAlreadyExist, logrus.InfoLevel},
	repository_user.NicknameAlreadyExist: {
		http.StatusUnprocessableEntity, handler_errors.NicknameAlreadyExist, logrus.InfoLevel},
	models.IncorrectEmailOrPassword: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectLoginOrPassword, logrus.InfoLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

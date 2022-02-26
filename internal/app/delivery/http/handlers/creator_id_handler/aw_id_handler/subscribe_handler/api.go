package aw_subscribe_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	repository_redis "patreon/internal/app/repository/pay_token/redis"
	usecase_pay_token "patreon/internal/app/usecase/pay_token"
	usecase_subscribers "patreon/internal/app/usecase/subscribers"

	"github.com/sirupsen/logrus"
)

var codesByErrorsPOST = base_handler.CodeMap{
	usecase_pay_token.InvalidUserToken: {
		http.StatusBadRequest, handler_errors.InvalidUserPayToken, logrus.WarnLevel},
	repository_redis.SetError: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	usecase_subscribers.SubscriptionAlreadyExists: {
		http.StatusConflict, handler_errors.UserAlreadySubscribe, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
var codesByErrorsDELETE = base_handler.CodeMap{
	usecase_subscribers.SubscriptionsNotFound: {
		http.StatusConflict, handler_errors.SubscribesNotFound, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

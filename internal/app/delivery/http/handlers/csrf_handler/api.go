package csrf_handler

import (
	"net/http"
	repository_jwt "glide/internal/app/csrf/repository/jwt"
	"glide/internal/app/delivery/http/handlers/base_handler"
	"glide/internal/app/delivery/http/handlers/handler_errors"

	"github.com/sirupsen/logrus"
)

var codeByErrors = base_handler.CodeMap{
	repository_jwt.ErrorSignedToken: {http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}

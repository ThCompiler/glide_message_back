package delivery

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"glide/internal/app"
)

type RespondError struct {
	Code  int
	Error error
	Level logrus.Level
}

type CodeMap map[error]RespondError

type ErrorConvertor struct {
	Responder
}

func (h *ErrorConvertor) UsecaseError(w http.ResponseWriter, r *http.Request, usecaseErr error, codeByErr CodeMap) {
	var generalError *app.GeneralError
	//orginalError := usecaseErr
	if errors.As(usecaseErr, &generalError) {
		usecaseErr = errors.Cause(usecaseErr).(*app.GeneralError).Err
	}

	respond := RespondError{http.StatusServiceUnavailable,
		errors.New("UnknownError"), logrus.ErrorLevel}

	for err, respondErr := range codeByErr {
		if errors.Is(usecaseErr, err) {
			respond = respondErr
			break
		}
	}

	//h.Log(w, r).Logf(respond.Level, "Gotted error: %v", orginalError)
	h.Error(w, r, respond.Code, respond.Error)
}

func (h *ErrorConvertor) HandlerError(w http.ResponseWriter, r *http.Request, code int, err error) {
	//h.Log(w, r).Errorf("Gotted error: %v", err)

	var generalError *app.GeneralError
	if errors.As(err, &generalError) {
		err = errors.Cause(err).(*app.GeneralError).Err
	}
	h.Error(w, r, code, err)
}

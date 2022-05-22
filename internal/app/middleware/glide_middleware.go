package middleware

import (
	"errors"
	"glide/internal/app/delivery/http/handlers"
	"glide/internal/app/repository"
	ucGlideMessage "glide/internal/app/usecase/glidemessage"
	"glide/internal/pkg/handler/handler_errors"
	hf "glide/internal/pkg/handler/handler_interfaces"
	"glide/internal/pkg/utilits"
	"glide/internal/pkg/utilits/delivery"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type GlideMessageMiddleware struct {
	log                 utilits.LogObject
	usecaseGlideMessage ucGlideMessage.Usecase
}

func NewGlideMessageMiddleware(log *logrus.Logger, usecaseGlideMessage ucGlideMessage.Usecase) *GlideMessageMiddleware {
	return &GlideMessageMiddleware{
		log:                 utilits.NewLogObject(log),
		usecaseGlideMessage: usecaseGlideMessage,
	}
}

// CheckCorrectUserFunc Errors
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 409 middleware.InvalidUserForApplyGlideMessage
func (mw *GlideMessageMiddleware) CheckCorrectUserFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := delivery.Responder{LogObject: mw.log}
		vars := mux.Vars(r)
		id, ok := vars[handlers.GlideMessageId]
		msgId, err := strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Info(vars)
			respond.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
			return
		}

		userID := r.Context().Value("user_id")
		if userID == nil {
			respond.Log(r).Error("can not get user_id from context")
			respond.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
			return
		}

		err = mw.usecaseGlideMessage.CheckAllowUser(msgId, userID.(string))

		if err != nil {
			if err != nil && !errors.Is(err, repository.NotFound) {
				mw.log.Log(r).Errorf("some error of bd awards %v", err)
				respond.Error(w, r, http.StatusInternalServerError, handler_errors.BDError)
				return
			}
			mw.log.Log(r).Warnf("this glide message %d not sent to this user %s", msgId, userID.(string))
			respond.Error(w, r, http.StatusForbidden, handler_errors.InvalidUserForApplyGlideMessage)
			return
		}

		next(w, r)
	}
}

func (mw *GlideMessageMiddleware) CheckCorrectUserId(handler http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckCorrectUserFunc(handler.ServeHTTP))
}

// CheckCorrectAuthorFunc Errors
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 409 middleware.InvalidAuthorForGlideMessage
func (mw *GlideMessageMiddleware) CheckCorrectAuthorFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := delivery.Responder{LogObject: mw.log}
		vars := mux.Vars(r)
		id, ok := vars[handlers.GlideMessageId]
		msgId, err := strconv.ParseInt(id, 10, 64)
		if !ok || err != nil {
			mw.log.Log(r).Info(vars)
			respond.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
			return
		}

		userID := r.Context().Value("user_id")
		if userID == nil {
			respond.Log(r).Error("can not get user_id from context")
			respond.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
			return
		}

		err = mw.usecaseGlideMessage.CheckAllowAuthor(msgId, userID.(string))

		if err != nil {
			if err != nil && !errors.Is(err, repository.NotFound) {
				mw.log.Log(r).Errorf("some error of bd awards %v", err)
				respond.Error(w, r, http.StatusInternalServerError, handler_errors.BDError)
				return
			}
			mw.log.Log(r).Warnf("this glide message %d are notcreated by this user %s", msgId, userID.(string))
			respond.Error(w, r, http.StatusForbidden, handler_errors.InvalidAuthorForGlideMessage)
			return
		}

		next(w, r)
	}
}

func (mw *GlideMessageMiddleware) CheckCorrectAuthor(handler http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckCorrectAuthorFunc(handler.ServeHTTP))
}

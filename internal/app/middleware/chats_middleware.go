package middleware

import (
	"errors"
	"glide/internal/app/delivery/http/handlers"
	"glide/internal/app/repository"
	usecase_chats "glide/internal/app/usecase/chats"
	"glide/internal/pkg/handler/handler_errors"
	hf "glide/internal/pkg/handler/handler_interfaces"
	"glide/internal/pkg/utilits"
	"glide/internal/pkg/utilits/delivery"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ChatsMiddleware struct {
	log          utilits.LogObject
	usecaseChats usecase_chats.Usecase
}

func NewChatsMiddleware(log *logrus.Logger, usecaseChats usecase_chats.Usecase) *ChatsMiddleware {
	return &ChatsMiddleware{
		log:          utilits.NewLogObject(log),
		usecaseChats: usecaseChats,
	}
}

// CheckCorrectChatIdFunc Errors
//		Status 400 middleware.InvalidParameters
//		Status 500 middleware.BDError
//		Status 409 middleware.IncorrectUserForChat
func (mw *ChatsMiddleware) CheckCorrectChatIdFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond := delivery.Responder{LogObject: mw.log}
		vars := mux.Vars(r)
		id, ok := vars[handlers.ChatId]
		chatId, err := strconv.ParseInt(id, 10, 64)
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

		err = mw.usecaseChats.CheckAllow(userID.(string), chatId)

		if err != nil {
			if err != nil && !errors.Is(err, repository.NotFound) {
				mw.log.Log(r).Errorf("some error of bd awards %v", err)
				respond.Error(w, r, http.StatusInternalServerError, handler_errors.BDError)
				return
			}
			mw.log.Log(r).Warnf("this chat %d not belongs to this user %s", chatId, userID.(string))
			respond.Error(w, r, http.StatusForbidden, handler_errors.IncorrectUserForChat)
			return
		}

		next(w, r)
	}
}

func (mw *ChatsMiddleware) CheckCorrectChatId(handler http.Handler) http.Handler {
	return http.HandlerFunc(mw.CheckCorrectChatIdFunc(handler.ServeHTTP))
}

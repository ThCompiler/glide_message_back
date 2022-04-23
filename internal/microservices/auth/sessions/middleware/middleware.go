package middleware

import (
	"context"
	"glide/internal/microservices/auth/delivery/grpc/client"
	"glide/internal/microservices/auth/sessions/sessions_manager"
	hf "glide/internal/pkg/handler/handler_interfaces"
	"glide/internal/pkg/utilits"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type SessionMiddleware struct {
	SessionClient client.AuthCheckerClient
	utilits.LogObject
}

func NewSessionMiddleware(authClient client.AuthCheckerClient, log *logrus.Logger) *SessionMiddleware {
	return &SessionMiddleware{
		SessionClient: authClient,
		LogObject:     utilits.NewLogObject(log),
	}
}

func (m *SessionMiddleware) updateCookie(w http.ResponseWriter, cook *http.Cookie) {
	cook.Expires = time.Now().Add(sessions_manager.ExpiredCookiesTime)
	cook.Path = "/"
	cook.HttpOnly = true
	http.SetCookie(w, cook)
}

func (m *SessionMiddleware) clearCookie(w http.ResponseWriter, cook *http.Cookie) {
	cook.Expires = time.Now().AddDate(0, 0, -1)
	cook.Path = "/"
	cook.HttpOnly = true
	http.SetCookie(w, cook)
}

// CheckFunc Errors:
//		Status 401 "not authorized user"
func (m *SessionMiddleware) CheckFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			m.Log(r).Warnf("in parsing cookie: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		uniqID := sessionID.Value
		if res, err := m.SessionClient.Check(context.Background(), uniqID); err != nil {
			m.Log(r).Warnf("Error in checking session: %v", err)
			m.clearCookie(w, sessionID)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			m.Log(r).Debugf("Get session for user: %d", res.UserID)
			r = r.WithContext(context.WithValue(r.Context(), "user_id", res.UserID))
			r = r.WithContext(context.WithValue(r.Context(), "session_id", res.UniqID))
			m.updateCookie(w, sessionID)
		}
		next(w, r)
	}
}

// Check Errors:
//		Status 401 "not authorized user"
func (m *SessionMiddleware) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(m.CheckFunc(next.ServeHTTP))
}

// CheckNotAuthorizedFunc Errors:
//		Status 418 "user already authorized"
func (m *SessionMiddleware) CheckNotAuthorizedFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			m.Log(r).Debug("User not Authorized")
			next.ServeHTTP(w, r)
			return
		}

		uniqID := sessionID.Value
		if res, err := m.SessionClient.Check(context.Background(), uniqID); err != nil {
			m.Log(r).Debug("User not Authorized")
			m.clearCookie(w, sessionID)
			next.ServeHTTP(w, r)
			return
		} else {
			m.Log(r).Warnf("UserAuthorized: %d", res.UserID)
			m.updateCookie(w, sessionID)
		}
		w.WriteHeader(http.StatusTeapot)
	}
}

// CheckNotAuthorized Errors:
//		Status 418 "user already authorized"
func (m *SessionMiddleware) CheckNotAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(m.CheckNotAuthorizedFunc(next.ServeHTTP))
}

// AddUserIdFunc Errors:
//		Nothing return only add user_id and session_id to context
func (m *SessionMiddleware) AddUserIdFunc(next hf.HandlerFunc) hf.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err == nil {
			uniqID := sessionID.Value
			if res, err := m.SessionClient.Check(context.Background(), uniqID); err == nil {
				m.Log(r).Debugf("Get session for user: %d", res.UserID)
				r = r.WithContext(context.WithValue(r.Context(), "user_id", res.UserID))
				r = r.WithContext(context.WithValue(r.Context(), "session_id", res.UniqID))
				m.updateCookie(w, sessionID)
			} else {
				m.clearCookie(w, sessionID)
			}
		}
		next(w, r)
	}
}

// AddUserId Errors:
//		Nothing return only add user_id and session_id to context
func (m *SessionMiddleware) AddUserId(next http.Handler) http.Handler {
	return http.HandlerFunc(m.AddUserIdFunc(next.ServeHTTP))
}

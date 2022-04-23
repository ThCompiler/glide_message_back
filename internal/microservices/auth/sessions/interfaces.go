package sessions

import "glide/internal/microservices/auth/sessions/models"

type SessionRepository interface {
	Set(session *models.Session) error
	GetUserId(key string, updExpiration int) (string, error)
	Del(session *models.Session) error
}

type SessionsManager interface {
	Check(uniqID string) (models.Result, error)
	Create(userID string) (models.Result, error)
	Delete(uniqID string) error
}

package sessions_manager

import (
	"encoding/hex"
	"fmt"
	"glide/internal/microservices/auth/sessions"
	"glide/internal/microservices/auth/sessions/models"
	"time"

	"golang.org/x/crypto/sha3"

	uuid "github.com/satori/go.uuid"
)

const (
	ExpiredCookiesTime = 48 * time.Hour
	UnknownUser        = ""
)

type SessionManager struct {
	sessionRepository sessions.SessionRepository
}

func NewSessionManager(sessionRep sessions.SessionRepository) *SessionManager {
	return &SessionManager{
		sessionRepository: sessionRep,
	}
}

func (manager *SessionManager) Create(nickname string) (models.Result, error) {
	strUserID := fmt.Sprintf("%d", nickname)

	session := &models.Session{
		UserID:     strUserID,
		UniqID:     generateUniqID(strUserID),
		Expiration: int(ExpiredCookiesTime.Milliseconds()),
	}
	if err := manager.sessionRepository.Set(session); err != nil {
		return models.Result{UserID: UnknownUser}, err
	}
	return models.Result{UserID: nickname, UniqID: session.UniqID}, nil
}

func (manager *SessionManager) Delete(uniqID string) error {
	session := &models.Session{
		UniqID: uniqID,
	}
	return manager.sessionRepository.Del(session)
}

func (manager *SessionManager) Check(uniqID string) (models.Result, error) {
	userID, err := manager.sessionRepository.GetUserId(uniqID, int(ExpiredCookiesTime.Milliseconds()))
	if err != nil {
		return models.Result{UserID: UnknownUser, UniqID: uniqID}, err
	}

	if err != nil {
		return models.Result{UserID: UnknownUser, UniqID: uniqID}, err
	}
	return models.Result{UserID: userID, UniqID: uniqID}, nil
}

func generateUniqID(userID string) string {
	value := append([]byte(userID), uuid.NewV4().Bytes()...)
	hash := sha3.Sum512(value)

	return hex.EncodeToString(hash[:])
}

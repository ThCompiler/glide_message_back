package models

type Session struct {
	UserID     string
	UniqID     string
	Expiration int
}

func (session *Session) String() string {
	return "User_id : " + session.UserID + "; session_id : " + session.UniqID
}

type Result struct {
	UserID string
	UniqID string
}

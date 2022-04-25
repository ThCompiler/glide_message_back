package push_client

type Pusher interface {
	NewMessage(messageId int64, companion string) error
	NewGlideMessage(companion int64, glideId int64) error
}

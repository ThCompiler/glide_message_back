package push

//go:generate easyjson -all -disallow_unknown_fields models.go

const (
	MessagePush = "Message"
	GlidePush   = "Glide"
)

//easyjson:json
type MessageInfo struct {
	Companion string `json:"companion"`
	MessageId int64  `json:"message_id"`
}

//easyjson:json
type GlideInfo struct {
	Companion string `json:"companion"`
	GlideId   int64  `json:"glide_id"`
}

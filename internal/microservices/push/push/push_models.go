package push_models

//go:generate easyjson -all -disallow_unknown_fields push_models.go

//easyjson:json
type MessagePush struct {
	ChatId          int64  `json:"chat_id"`
	Companion       string `json:"companion"`
	CompanionAvatar string `json:"companion_avatar"`
	MessageId       int64  `json:"message_id"`
	Text            string `json:"text"`
}

//easyjson:json
type GlidePush struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	Message      string `json:"message"`
	Country      string `json:"country"`
	Author       string `json:"author"`
	AuthorAvatar string `json:"author_avatar"`
}

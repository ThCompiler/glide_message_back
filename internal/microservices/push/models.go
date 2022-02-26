package push

import "time"

//go:generate easyjson -all -disallow_unknown_fields models.go

const (
	CommentPush = "Comment"
	PostPush    = "Post"
	NewSubPush  = "Subscriber"
)

//easyjson:json
type PostInfo struct {
	CreatorId int64     `json:"creator_id"`
	PostId    int64     `json:"post_id"`
	PostTitle string    `json:"post_title"`
	Date      time.Time `json:"date"`
}

//easyjson:json
type SubInfo struct {
	AwardsId  int64     `json:"awards_id"`
	UserId    int64     `json:"user_id"`
	CreatorId int64     `json:"creator_id"`
	Date      time.Time `json:"date"`
}

//easyjson:json
type CommentInfo struct {
	CommentId int64     `json:"comment_id"`
	AuthorId  int64     `json:"author_id"`
	PostId    int64     `json:"post_id"`
	Date      time.Time `json:"date"`
}

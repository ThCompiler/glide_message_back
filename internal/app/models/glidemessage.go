package models

import (
	"fmt"
	"time"
)

type GlideMessage struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	Picture        string    `json:"picture,omitempty"`
	Created        time.Time `json:"created"`
	Author         string    `json:"author"`
	AuthorFullname string    `json:"author_fullname"`
	AuthorAvatar   string    `json:"author_avatar,omitempty"`
	Country        string    `json:"country"`
}

func (ps *GlideMessage) String() string {
	return fmt.Sprintf("{ID: %d, Title: %s, Author: %s}", ps.ID,
		ps.Title, ps.Author)
}

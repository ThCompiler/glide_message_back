package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	models_utilits "glide/internal/pkg/utilits/models"
	"time"
)

type Message struct {
	ID       int64     `json:"id"`
	Text     string    `json:"text"`
	Picture  string    `json:"picture,omitempty"`
	Created  time.Time `json:"created"`
	Author   string    `json:"author"`
	IsViewed bool      `json:"is_viewed"`
}

type GlideMessage struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Picture string    `json:"picture,omitempty"`
	Created time.Time `json:"created"`
	Author  string    `json:"author"`
	Country string    `json:"country"`
}

func (ps *GlideMessage) String() string {
	return fmt.Sprintf("{ID: %d, Title: %s, Author: %s}", ps.ID,
		ps.Title, ps.Author)
}

type Chat struct {
	ID              int64   `json:"id"`
	Companion       string  `json:"companion"`
	CompanionAvatar string  `json:"companion_avatar"`
	LastMessage     Message `json:"last_message"`
}

func (ps *Message) String() string {
	return fmt.Sprintf("{ID: %d, Text: %s, Author: %s}", ps.ID,
		ps.Text, ps.Author)
}

func (ps *Chat) String() string {
	return fmt.Sprintf("{ID: %d, companion: %s, LastMessage: %s}", ps.ID,
		ps.Companion, ps.LastMessage.String())
}

// Validate Errors:
//		EmptyTitle
//		InvalidAwardsId
// Important can return some other error
func (ps *Message) Validate() error {
	err := validation.Errors{
		"text": validation.Validate(ps.Text, validation.Required),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = models_utilits.ExtractValidateError(messageValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}

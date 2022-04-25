package http_models

import (
	"glide/internal/app/csrf/csrf_models"
	"glide/internal/app/models"
	"time"
)

//go:generate easyjson -all -disallow_unknown_fields response_models.go

//easyjson:json
type TokenResponse struct {
	Token csrf_models.Token `json:"token"`
}

//easyjson:json
type ErrResponse struct {
	Err string `json:"error"`
}

//easyjson:json
type OkResponse struct {
	Ok string `json:"OK"`
}

//easyjson:json
type IdResponse struct {
	ID int64 `json:"id"`
}

//easyjson:json
type ProfileResponse struct {
	Nickname  string   `json:"nickname"`
	Fullname  string   `json:"fullname"`
	Avatar    string   `json:"avatar"`
	About     string   `json:"about,omitempty"`
	Age       int64    `json:"age"`
	Country   string   `json:"country"`
	Languages []string `json:"languages"`
}

//easyjson:json
type ResponseInfo struct {
	models.Info
}

//easyjson:json
type ResponseMessage struct {
	ID       int64     `json:"id"`
	Text     string    `json:"text"`
	Picture  string    `json:"picture,omitempty"`
	Created  time.Time `json:"created"`
	Author   string    `json:"author"`
	IsViewed bool      `json:"is_viewed"`
}

//easyjson:json
type ResponseChat struct {
	ID              int64           `json:"id"`
	Companion       string          `json:"companion"`
	CompanionAvatar string          `json:"companion_avatar"`
	LastMessage     ResponseMessage `json:"last_message"`
}

//easyjson:json
type ResponseChats []ResponseChat

//easyjson:json
type ResponseMessages []ResponseMessage

func ToProfileResponse(us models.User) ProfileResponse {
	return ProfileResponse{
		Nickname:  us.Nickname,
		Avatar:    us.Avatar,
		Fullname:  us.Fullname,
		Languages: us.Languages,
		Country:   us.Country,
		About:     us.About,
		Age:       us.Age,
	}
}

func ToResponseInfo(info models.Info) ResponseInfo {
	return ResponseInfo{
		info,
	}
}

func ToResponseMessage(msg models.Message) ResponseMessage {
	return ResponseMessage{
		ID:       msg.ID,
		Author:   msg.Author,
		Text:     msg.Text,
		Picture:  msg.Picture,
		Created:  msg.Created,
		IsViewed: msg.IsViewed,
	}
}

func ToResponseMessages(msgs []models.Message) ResponseMessages {
	respondMessages := make([]ResponseMessage, len(msgs))
	for i, msg := range msgs {
		respondMessages[i] = ToResponseMessage(msg)
	}
	return respondMessages
}

func ToResponseChat(cht models.Chat) ResponseChat {
	return ResponseChat{
		ID:              cht.ID,
		Companion:       cht.Companion,
		CompanionAvatar: cht.CompanionAvatar,
		LastMessage:     ToResponseMessage(cht.LastMessage),
	}
}

func ToResponseChats(crs []models.Chat) ResponseChats {
	respondChats := make([]ResponseChat, len(crs))
	for i, cr := range crs {
		respondChats[i] = ToResponseChat(cr)
	}
	return respondChats
}

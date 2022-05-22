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
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

//easyjson:json
type ResponseInfos []ResponseInfo

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
	ID              int64            `json:"id"`
	Companion       string           `json:"companion"`
	CompanionAvatar string           `json:"companion_avatar"`
	CountNotViewed  int64            `json:"count_not_viewed"`
	LastMessage     *ResponseMessage `json:"last_message,omitempty"`
}

//easyjson:json
type ResponseGlideMessage struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Picture string    `json:"picture,omitempty"`
	Created time.Time `json:"created"`
	Author  string    `json:"author"`
	Country string    `json:"country"`
}

//easyjson:json
type ResponseGlideMessages []ResponseGlideMessage

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

func LanguageToResponseInfo(info models.InfoLanguage) ResponseInfo {
	return ResponseInfo{
		Name:    info.Language,
		Picture: info.Picture,
	}
}

func CountryToResponseInfo(info models.InfoCountry) ResponseInfo {
	return ResponseInfo{
		Name:    info.CountryName,
		Picture: info.Picture,
	}
}

func CountriesToInfos(msgs []models.InfoCountry) ResponseInfos {
	respondInfos := make([]ResponseInfo, len(msgs))
	for i, msg := range msgs {
		respondInfos[i] = CountryToResponseInfo(msg)
	}
	return respondInfos
}

func LanguagesToInfos(msgs []models.InfoLanguage) ResponseInfos {
	respondInfos := make([]ResponseInfo, len(msgs))
	for i, msg := range msgs {
		respondInfos[i] = LanguageToResponseInfo(msg)
	}
	return respondInfos
}

func ToResponseMessage(msg models.Message) *ResponseMessage {
	return &ResponseMessage{
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
		respondMessages[i] = *ToResponseMessage(msg)
	}
	return respondMessages
}

func ToResponseChat(cht models.Chat) ResponseChat {
	if cht.LastMessage == nil {
		return ResponseChat{
			ID:              cht.ID,
			Companion:       cht.Companion,
			CompanionAvatar: cht.CompanionAvatar,
			CountNotViewed:  cht.CountNotViewed,
			LastMessage:     nil,
		}
	}
	return ResponseChat{
		ID:              cht.ID,
		Companion:       cht.Companion,
		CompanionAvatar: cht.CompanionAvatar,
		CountNotViewed:  cht.CountNotViewed,
		LastMessage:     ToResponseMessage(*cht.LastMessage),
	}
}

func ToResponseGlideMessage(msg models.GlideMessage) ResponseGlideMessage {
	return ResponseGlideMessage{
		ID:      msg.ID,
		Author:  msg.Author,
		Message: msg.Message,
		Title:   msg.Title,
		Created: msg.Created,
		Picture: msg.Picture,
		Country: msg.Country,
	}
}

func ToResponseGlideMessages(msgs []models.GlideMessage) ResponseGlideMessages {
	respondGlideMessage := make([]ResponseGlideMessage, len(msgs))
	for i, msg := range msgs {
		respondGlideMessage[i] = ToResponseGlideMessage(msg)
	}
	return respondGlideMessage
}

func ToResponseChats(crs []models.Chat) ResponseChats {
	respondChats := make([]ResponseChat, len(crs))
	for i, cr := range crs {
		respondChats[i] = ToResponseChat(cr)
	}
	return respondChats
}

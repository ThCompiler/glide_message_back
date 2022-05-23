package http_models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"glide/internal/app/models"
	"glide/internal/pkg/handler/handler_errors"
	models_utilits "glide/internal/pkg/utilits/models"
	"image/color"
)

//go:generate easyjson -all -disallow_unknown_fields request_models.go

//easyjson:json
type RequestMessageIds []int64

func (req *RequestMessageIds) ToArray() []int64 {
	return *req
}

//easyjson:json
type RequestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//easyjson:json
type RequestChangePassword struct {
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
}

//easyjson:json
type RequestChangeNickname struct {
	OldNickname string `json:"old"`
	NewNickname string `json:"new"`
}

//easyjson:json
type RequestRegistration struct {
	Nickname  string   `json:"nickname"`
	Fullname  string   `json:"fullname"`
	About     string   `json:"about,omitempty"`
	Age       int64    `json:"age"`
	Country   string   `json:"country,omitempty"`
	Languages []string `json:"languages,omitempty"`
	Password  string   `json:"password"`
}

//easyjson:json
type RequestGlideMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Author  string `json:"author,omitempty"`
}

func (rr *RequestGlideMessage) ToGlideMessage() *models.GlideMessage {
	return &models.GlideMessage{
		Title:   rr.Title,
		Message: rr.Message,
		Author:  rr.Author,
	}
}

func (rr *RequestRegistration) ToUser() *models.User {
	return &models.User{
		Nickname:  rr.Nickname,
		Fullname:  rr.Fullname,
		Languages: rr.Languages,
		Country:   rr.Country,
		About:     rr.About,
		Age:       rr.Age,
		Password:  rr.Password,
	}
}

//easyjson:json
type RequestUserUpdate struct {
	Fullname  string   `json:"fullname,omitempty"`
	About     string   `json:"about,omitempty"`
	Age       int64    `json:"age,omitempty"`
	Country   string   `json:"country,omitempty"`
	Languages []string `json:"languages,omitempty"`
}

func (rr *RequestUserUpdate) ToUser() *models.User {
	return &models.User{
		Fullname:  rr.Fullname,
		Languages: rr.Languages,
		Country:   rr.Country,
		About:     rr.About,
		Age:       rr.Age,
	}
}

//easyjson:json
type Color struct {
	R uint8 `json:"red"`
	G uint8 `json:"green"`
	B uint8 `json:"blue"`
	A uint8 `json:"alpha"`
}

func NewColor(rgba color.RGBA) Color {
	return Color{
		R: rgba.R,
		G: rgba.G,
		B: rgba.B,
		A: rgba.A,
	}
}

func (req *RequestChangeNickname) Validate() error {
	err := validation.Errors{
		"old_nickname": validation.Validate(req.OldNickname, validation.Required,
			validation.Length(models.MIN_NICKNAME_LENGTH, models.MAX_NICKNAME_LENGTH)),
		"new_nickname": validation.Validate(req.NewNickname, validation.Required,
			validation.Length(models.MIN_NICKNAME_LENGTH, models.MAX_NICKNAME_LENGTH)),
	}.Filter()
	if err != nil {
		return NicknameValidateError
	}
	return nil
}

// requestAttachValidError Errors:
//		handler_errors.IncorrectType
//		handler_errors.IncorrectIdAttach
//      handler_errors.IncorrectStatus
func requestAttachValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"type":   handler_errors.IncorrectType,
		"id":     handler_errors.IncorrectIdAttach,
		"status": handler_errors.IncorrectStatus,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}

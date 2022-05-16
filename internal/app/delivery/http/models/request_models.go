package http_models

import (
	"glide/internal/app/delivery/http/handlers/handler_errors"
	"glide/internal/app/models"
	models_utilits "glide/internal/pkg/utilits/models"
	"image/color"

	validation "github.com/go-ozzo/ozzo-validation"
)

//go:generate easyjson -all -disallow_unknown_fields request_models.go

//easyjson:json
type RequestCreator struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}

//easyjson:json
type RequestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//easyjson:json
type RequestComment struct {
	Body      string `json:"body"`
	AsCreator bool   `json:"as_creator,omitempty"`
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
	Fullname  string   `json:"fullname"`
	About     string   `json:"about,omitempty"`
	Age       int64    `json:"age"`
	Country   string   `json:"country"`
	Languages []string `json:"languages"`
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

//easyjson:json
type RequestAwards struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price"`
	Color       Color  `json:"color,omitempty"`
}

//easyjson:json
type RequestPosts struct {
	Title       string `json:"title,omitempty"`
	AwardsId    int64  `json:"awards_id,omitempty"`
	Description string `json:"description,omitempty"`
	IsDraft     bool   `json:"is_draft,omitempty"`
}

//easyjson:json
type RequestText struct {
	Text string `json:"text"`
}

type SubscribeRequest struct {
	Token string `json:"pay_token"`
}

func (req *SubscribeRequest) Validate() error {
	err := validation.Errors{
		"pay_token": validation.Validate(req.Token, validation.Required, validation.Length(1, 0)),
	}.Filter()
	if err != nil {
		return TokenValidateError
	}
	return nil
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

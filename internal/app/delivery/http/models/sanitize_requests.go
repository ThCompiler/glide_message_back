package http_models

import (
	"github.com/microcosm-cc/bluemonday"
)

func (req *RequestCreator) Sanitize(sanitizer bluemonday.Policy) {
	req.Category = sanitizer.Sanitize(req.Category)
	req.Description = sanitizer.Sanitize(req.Description)
}

func (req *RequestLogin) Sanitize(sanitizer bluemonday.Policy) {
	req.Login = sanitizer.Sanitize(req.Login)
	req.Password = sanitizer.Sanitize(req.Password)
}

func (req *RequestMessageIds) Sanitize(_ bluemonday.Policy) {}

func (req *RequestChangePassword) Sanitize(sanitizer bluemonday.Policy) {
	req.OldPassword = sanitizer.Sanitize(req.OldPassword)
	req.NewPassword = sanitizer.Sanitize(req.NewPassword)
}

func (req *RequestRegistration) Sanitize(sanitizer bluemonday.Policy) {
	req.Nickname = sanitizer.Sanitize(req.Nickname)
	req.Password = sanitizer.Sanitize(req.Password)
	req.Fullname = sanitizer.Sanitize(req.Fullname)
	req.About = sanitizer.Sanitize(req.About)
	req.Country = sanitizer.Sanitize(req.Country)
	for id, lang := range req.Languages {
		req.Languages[id] = sanitizer.Sanitize(lang)
	}
}

func (req *RequestUserUpdate) Sanitize(sanitizer bluemonday.Policy) {
	req.Fullname = sanitizer.Sanitize(req.Fullname)
	req.About = sanitizer.Sanitize(req.About)
	req.Country = sanitizer.Sanitize(req.Country)
	for id, lang := range req.Languages {
		req.Languages[id] = sanitizer.Sanitize(lang)
	}
}

func (req *RequestComment) Sanitize(sanitizer bluemonday.Policy) {
	req.Body = sanitizer.Sanitize(req.Body)
}

func (req *RequestAwards) Sanitize(sanitizer bluemonday.Policy) {
	req.Name = sanitizer.Sanitize(req.Name)
	req.Description = sanitizer.Sanitize(req.Description)
}

func (req *RequestPosts) Sanitize(sanitizer bluemonday.Policy) {
	req.Title = sanitizer.Sanitize(req.Title)
	req.Description = sanitizer.Sanitize(req.Description)
}

func (req *RequestText) Sanitize(sanitizer bluemonday.Policy) {
	req.Text = sanitizer.Sanitize(req.Text)
}

func (req *SubscribeRequest) Sanitize(sanitizer bluemonday.Policy) {
	req.Token = sanitizer.Sanitize(req.Token)
}
func (req *RequestChangeNickname) Sanitize(sanitizer bluemonday.Policy) {
	req.OldNickname = sanitizer.Sanitize(req.OldNickname)
	req.NewNickname = sanitizer.Sanitize(req.NewNickname)
}

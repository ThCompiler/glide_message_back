package http_models

import (
	"github.com/microcosm-cc/bluemonday"
)

func (req *RequestLogin) Sanitize(sanitizer bluemonday.Policy) {
	req.Login = sanitizer.Sanitize(req.Login)
	req.Password = sanitizer.Sanitize(req.Password)
}

func (req *RequestGlideMessage) Sanitize(sanitizer bluemonday.Policy) {
	req.Title = sanitizer.Sanitize(req.Title)
	req.Message = sanitizer.Sanitize(req.Message)
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

func (req *RequestChangeNickname) Sanitize(sanitizer bluemonday.Policy) {
	req.OldNickname = sanitizer.Sanitize(req.OldNickname)
	req.NewNickname = sanitizer.Sanitize(req.NewNickname)
}

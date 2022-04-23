package models

import (
	"fmt"
	"github.com/pkg/errors"
	"glide/internal/pkg/utilits/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

const (
	MIN_NICKNAME_LENGTH = 4
	MAX_NICKNAME_LENGTH = 25
	MIN_PASSWORD_LENGTH = 6
	MAX_PASSWORD_LENGTH = 50

	EmptyAge = 0
)

type User struct {
	Nickname          string   `json:"nickname"`
	Password          string   `json:"password,omitempty"`
	Fullname          string   `json:"fullname"`
	About             string   `json:"about,omitempty"`
	EncryptedPassword string   `json:",omitempty"`
	Avatar            string   `json:"avatar,omitempty"`
	Age               int64    `json:"age"`
	Country           string   `json:"country,omitempty"`
	Languages         []string `json:"languages,omitempty"`
}

func (u *User) String() string {
	return fmt.Sprintf("{Login: %s}", u.Nickname)
}

// ValidateUpdate Errors:
//		IncorrectAge
// Important can return some other error
func (u *User) ValidateUpdate() error {
	err := validation.Errors{
		"age": validation.Validate(u.Nickname, validation.Min(EmptyAge)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(err, "failed error getting in validate user")
	}

	if knowError = models_utilits.ExtractValidateError(userValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}

// Validate Errors:
//		IncorrectNicknameOrPassword
//		IncorrectAge
// Important can return some other error
func (u *User) Validate() error {
	err := validation.Errors{
		"password": validation.Validate(u.Password, validation.By(models_utilits.RequiredIf(u.EncryptedPassword == "")),
			validation.Length(MIN_PASSWORD_LENGTH, MAX_PASSWORD_LENGTH)),
		"nickname": validation.Validate(u.Nickname, validation.Required, validation.Length(MIN_NICKNAME_LENGTH, MAX_NICKNAME_LENGTH)),
		"age":      validation.Validate(u.Nickname, validation.Min(EmptyAge)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(err, "failed error getting in validate user")
	}

	if knowError = models_utilits.ExtractValidateError(userValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}

func (u *User) MakeEmptyPassword() {
	u.Password = ""
	u.EncryptedPassword = ""
}

// Encrypt Errors:
// 		EmptyPassword
// Important can return some other error
func (u *User) Encrypt() error {
	if len(u.Password) == 0 {
		return EmptyPassword
	}
	enc, err := u.encryptString(u.Password)
	if err != nil {
		return err
	}
	u.EncryptedPassword = enc
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func (u *User) encryptString(s string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(enc), nil
}

package http_models

import (
	"fmt"
	"glide/internal/app/models"

	"github.com/pkg/errors"
)

var (
	TokenValidateError    = errors.New("invalid pay_token")
	NicknameValidateError = errors.New(fmt.Sprintf("invalid nickname in body len must be from %v to %v",
		models.MIN_NICKNAME_LENGTH, models.MAX_NICKNAME_LENGTH))
)

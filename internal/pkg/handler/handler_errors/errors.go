package handler_errors

import (
	"errors"
	"fmt"
	http_models "glide/internal/app/delivery/http/handlers"
	"glide/internal/app/models"
)

/// NOT FOUND
var (
	ChatNotFound = errors.New("chat not found")
	UserNotFound = errors.New("user not found")
)

/// File parse error
var (
	IncorrectType = errors.New(
		fmt.Sprintf("Not allow type, allowed type is: image"))
	IncorrectIdAttach = errors.New("Not valid attach id")
	IncorrectStatus   = errors.New(fmt.Sprintf("Not allow status, allowed status is: %s, %s",
		http_models.AddStatus, http_models.UpdateStatus))
)

/// Fields Incorrect
var (
	IncorrectCreatorId       = errors.New("this creator id not know")
	IncorrectLoginOrPassword = errors.New("incorrect nickname or password")
)

// BD Error
var (
	ChatAlreadyExist     = errors.New("chat already exist")
	NicknameAlreadyExist = errors.New("nickname already exist")
	BDError              = errors.New("can not do bd operation")
)

// Session Error
var (
	ErrorCreateSession = errors.New("can not create session")
	DeleteCookieFail   = errors.New("can not delete cookie from session store")
)

// Request Error
var (
	InvalidBody          = errors.New("invalid body in request")
	InvalidUserAge       = errors.New("invalid user age in request")
	InvalidUserCounty    = errors.New("invalid user county in request")
	InvalidUserLanguage  = errors.New("invalid user language in request")
	InvalidParameters    = errors.New("invalid parameters")
	InvalidQueries       = errors.New("invalid parameters in query")
	FileSizeError        = errors.New("size of file very big")
	InvalidFormFieldName = errors.New("invalid form field name for load file")
	InvalidExt           = errors.New("please upload: ")
	InvalidUserNickname  = errors.New(fmt.Sprintf("invalid nickname in body len must be from %v to %v",
		models.MIN_NICKNAME_LENGTH, models.MAX_NICKNAME_LENGTH))
	InvalidUserForApplyGlideMessage = errors.New("this user was not gotten this glide message")
	InvalidAuthorForGlideMessage    = errors.New("this user was not author of this glide message")
	IncorrectUserForChat            = errors.New("this chat not belongs this user")
)

var InternalError = errors.New("server error")

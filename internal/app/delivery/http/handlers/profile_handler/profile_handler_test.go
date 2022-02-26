package profile_handler

import (
	"bytes"
	"context"
	"github.com/mailru/easyjson"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	models_http "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProfileTestSuite struct {
	handlers.SuiteHandler
	handler *ProfileHandler
}

func (s *ProfileTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewProfileHandler(s.Logger, s.MockSessionsManager, s.MockUserUsecase)
}

func (s *ProfileTestSuite) TestProfileHandler_GET_Correct() {
	userID := int64(1)
	test := handlers.TestTable{
		Name:              "correct",
		Data:              &models_http.ProfileResponse{ID: userID, Login: "some", Nickname: "done"},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}

	recorder := httptest.NewRecorder()
	b := bytes.Buffer{}

	ctx := context.WithValue(context.Background(), "user_id", userID)
	reader, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/user", &b)

	s.MockUserUsecase.
		EXPECT().
		GetProfile(userID).
		Times(test.ExpectedMockTimes).
		Return(&models.User{ID: userID, Login: "some", Nickname: "done"}, nil)
	s.handler.GET(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)

	req := &models_http.ProfileResponse{}
	err := easyjson.UnmarshalFromReader(recorder.Body, req)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), req, &models_http.ProfileResponse{ID: userID, Login: test.Data.(*models_http.ProfileResponse).Login,
		Nickname: test.Data.(*models_http.ProfileResponse).Nickname,
		Avatar:   test.Data.(*models_http.ProfileResponse).Avatar})
}

func (s *ProfileTestSuite) TestProfileHandler_GET_NotFound() {
	userID := int64(1)
	s.Tb = handlers.TestTable{
		Name:              "with not found",
		Data:              nil,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusNotFound,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	ctx := context.WithValue(context.Background(), "user_id", userID)
	reader, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/user", &b)

	s.MockUserUsecase.EXPECT().
		GetProfile(userID).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, repository.NotFound)
	s.handler.GET(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func (s *ProfileTestSuite) TestProfileHandler_GET_WithoutContext() {
	userID := int64(1)
	s.Tb = handlers.TestTable{
		Name:              "without context",
		Data:              nil,
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusInternalServerError,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	reader, _ := http.NewRequest(http.MethodGet, "/user", &b)

	s.MockUserUsecase.
		EXPECT().
		GetProfile(userID).
		Times(s.Tb.ExpectedMockTimes)
	s.handler.GET(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func TestProfileHandler(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

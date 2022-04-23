package usecase

import (
	"io"
	mock_repository "glide/internal/app/repository/access/mocks"
	mock_repository_attaches "glide/internal/app/repository/attaches/mocks"
	mock_repository_awards "glide/internal/app/repository/awards/mocks"
	mock_repository_creator "glide/internal/app/repository/creator/mocks"
	mock_repository_info "glide/internal/app/repository/info/mocks"
	mock_repository_likes "glide/internal/app/repository/likes/mocks"
	mock_repository_posts "glide/internal/app/repository/posts/mocks"
	mock_repository_subscribers "glide/internal/app/repository/subscribers/mocks"
	mock_repository_user "glide/internal/app/repository/user/mocks"
	mock_files "glide/internal/microservices/files/delivery/grpc/client/mocks"
	mock_utils "glide/internal/pkg/utils/mocks"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type TestTable struct {
	Name              string
	Data              interface{}
	ExpectedMockTimes int
	ExpectedError     error
}
type SuiteUsecase struct {
	suite.Suite
	Mock                      *gomock.Controller
	MockCreatorRepository     *mock_repository_creator.CreatorRepository
	MockUserRepository        *mock_repository_user.UserRepository
	MockSubscribersRepository *mock_repository_subscribers.SubscribersRepository
	MockAwardsRepository      *mock_repository_awards.AwardsRepository
	MockPostsRepository       *mock_repository_posts.PostsRepository
	MockLikesRepository       *mock_repository_likes.LikesRepository
	MockAccessRepository      *mock_repository.AccessRepository
	MockInfoRepository        *mock_repository_info.InfoRepository
	MockAttachesRepository    *mock_repository_attaches.AttachesRepository
	MockFileClient            *mock_files.MockFileServiceClient
	MockConvector             *mock_utils.MockImageConverter
	MockSubscriberRepository  *mock_repository_subscribers.SubscribersRepository
	Logger                    *logrus.Logger
	Tb                        TestTable
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockCreatorRepository = mock_repository_creator.NewCreatorRepository(s.Mock)
	s.MockUserRepository = mock_repository_user.NewUserRepository(s.Mock)
	s.MockSubscribersRepository = mock_repository_subscribers.NewSubscribersRepository(s.Mock)
	s.MockFileClient = mock_files.NewMockFileServiceClient(s.Mock)
	s.MockAttachesRepository = mock_repository_attaches.NewAttachesRepository(s.Mock)
	s.MockPostsRepository = mock_repository_posts.NewPostsRepository(s.Mock)
	s.MockAwardsRepository = mock_repository_awards.NewAwardsRepository(s.Mock)
	s.MockLikesRepository = mock_repository_likes.NewLikesRepository(s.Mock)
	s.MockInfoRepository = mock_repository_info.NewInfoRepository(s.Mock)
	s.MockConvector = mock_utils.NewMockImageConverter(s.Mock)
	s.MockAccessRepository = mock_repository.NewAccessRepository(s.Mock)

	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}

package handler_factory

import (
	useCsrf "patreon/internal/app/csrf/usecase"
	useAttaches "patreon/internal/app/usecase/attaches"
	useAwards "patreon/internal/app/usecase/awards"
	useComments "patreon/internal/app/usecase/comments"
	useCreator "patreon/internal/app/usecase/creator"
	useInfo "patreon/internal/app/usecase/info"
	useLikes "patreon/internal/app/usecase/likes"
	usePayToken "patreon/internal/app/usecase/pay_token"
	usePayments "patreon/internal/app/usecase/payments"
	usePosts "patreon/internal/app/usecase/posts"
	useStats "patreon/internal/app/usecase/statistics"
	useSubscr "patreon/internal/app/usecase/subscribers"
	useUser "patreon/internal/app/usecase/user"
)

//go:generate mockgen -destination=mocks/mock_usecase_factory.go -package=mock_usecase_factory . UsecaseFactory

type UsecaseFactory interface {
	GetUserUsecase() useUser.Usecase
	GetCreatorUsecase() useCreator.Usecase
	GetCsrfUsecase() useCsrf.Usecase
	GetAwardsUsecase() useAwards.Usecase
	GetPostsUsecase() usePosts.Usecase
	GetSubscribersUsecase() useSubscr.Usecase
	GetLikesUsecase() useLikes.Usecase
	GetAttachesUsecase() useAttaches.Usecase
	GetPaymentsUsecase() usePayments.Usecase
	GetInfoUsecase() useInfo.Usecase
	GetCommentsUsecase() useComments.Usecase
	GetStatsUsecase() useStats.Usecase
	GetPayTokenUsecase() usePayToken.Usecase
}

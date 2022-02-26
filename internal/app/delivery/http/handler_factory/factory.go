package handler_factory

import (
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/creator_handler"
	search_creators_handler "patreon/internal/app/delivery/http/handlers/creator_handler/search_creators"
	"patreon/internal/app/delivery/http/handlers/creator_handler/subscribe_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_id_handler"
	aw_subscribe_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_id_handler/subscribe_handler"
	aw_upd_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_id_handler/upd_aw_handler"
	upd_cover_awards_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/aw_id_handler/upd_cover_awards"
	creator_payments_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/payments_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_handler/upl_audio_attach_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_handler/upl_img_attach_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_handler/upl_text_attach_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_handler/upl_video_attach_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_id_handler"
	upd_audio_attach_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_id_handler/upd_audio_post_handler"
	upd_img_data_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_id_handler/upd_image_post_handler"
	upd_text_data_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_id_handler/upd_text_post_handler"
	upd_video_attach_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/attaches_id_handler/upd_video_post_handler"
	comments_id_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/comment_id_handler"
	comments_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/comments_handler"
	"patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/likes_handler"
	upl_cover_posts_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/upd_cover_post_handler"
	posts_upd_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/posts_id_handler/upd_handler"
	statistics_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/statistics_handler/creator_subscribers_handler"
	statistics_total_income_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/statistics_handler/creator_total_income_handler"
	statistics_count_posts_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/statistics_handler/posts_handler/creator_count_posts_handler"
	statistics_count_posts_views_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/statistics_handler/posts_handler/creator_count_posts_views_handler"
	upd_avatar_creator_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/upd_avatar_handler"
	upd_cover_creator_handler "patreon/internal/app/delivery/http/handlers/creator_id_handler/upd_cover_handler"
	"patreon/internal/app/delivery/http/handlers/csrf_handler"
	"patreon/internal/app/delivery/http/handlers/info_handler"
	"patreon/internal/app/delivery/http/handlers/login_handler"
	"patreon/internal/app/delivery/http/handlers/logout_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/payments_handler"
	pay_account_handler "patreon/internal/app/delivery/http/handlers/profile_handler/payments_handler/account_handler"
	pay_token_handler "patreon/internal/app/delivery/http/handlers/profile_handler/payments_handler/token_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/subscriptions_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/avatar_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/nickname_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/password_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/user_comments_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/user_posts_handler"
	"patreon/internal/app/delivery/http/handlers/register_handler"
	"patreon/internal/microservices/auth/delivery/grpc/client"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
)

const (
	ROOT = iota
	INFO
	REGISTER
	LOGIN
	LOGOUT
	PROFILE
	CREATORS
	CREATOR_WITH_ID
	CREATOR_AVATAR
	CREATOR_COVER
	SEARCH_CREATORS
	UPDATE_PASSWORD
	UPDATE_AVATAR
	UPDATE_NICKNAME
	AWARDS
	AWARDS_WITH_ID
	AWARDS_UPDATE
	AWARDS_COVER
	POSTS
	POSTS_AVAILABLE
	POSTS_WITH_ID
	POSTS_UPD
	POSTS_LIKES
	GET_CSRF_TOKEN
	GET_USER_SUBSCRIPTIONS
	POST_UPD_COVER
	ATTACH_ADD_TEXT
	ATTACH_ADD_IMAGE
	ATTACH_ADD_AUDIO
	ATTACH_ADD_VIDEO
	ATTACH_UPD_TEXT
	ATTACH_UPD_IMAGE
	ATTACH_ID
	SUBSCRIBES
	ATTACHES
	AWARDS_CREATOR_SUBSCRIBE
	USER_PAYMENTS
	ATTACH_UPD_VIDEO
	ATTACH_UPD_AUDIO
	STATS_COUNT_POSTS
	STATS_POSTS_VIEWS
	STATS_COUNT_SUBSCRIBERS
	STATS_TOTAL_INCOMES
	POST_COMMENTS
	COMMENTS_ID
	USER_COMMENTS
	USER_PAYMENTS_TOKEN
	PAYMENTS_ACCOUNT
	CREATOR_PAYMENTS
)

type HandlerFactory struct {
	usecaseFactory    UsecaseFactory
	sessionClientConn *grpc.ClientConn
	logger            *logrus.Logger
	urlHandler        *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, usecaseFactory UsecaseFactory, sClientConn *grpc.ClientConn) *HandlerFactory {
	return &HandlerFactory{
		usecaseFactory:    usecaseFactory,
		logger:            logger,
		sessionClientConn: sClientConn,
	}
}

func (f *HandlerFactory) initAllHandlers() map[int]app.Handler {
	ucUser := f.usecaseFactory.GetUserUsecase()
	ucCreator := f.usecaseFactory.GetCreatorUsecase()
	ucCsrf := f.usecaseFactory.GetCsrfUsecase()
	ucAwards := f.usecaseFactory.GetAwardsUsecase()
	ucPosts := f.usecaseFactory.GetPostsUsecase()
	ucLikes := f.usecaseFactory.GetLikesUsecase()
	ucSubscr := f.usecaseFactory.GetSubscribersUsecase()
	ucAttaches := f.usecaseFactory.GetAttachesUsecase()
	ucPayments := f.usecaseFactory.GetPaymentsUsecase()
	ucInfo := f.usecaseFactory.GetInfoUsecase()
	ucComment := f.usecaseFactory.GetCommentsUsecase()
	sManager := client.NewSessionClient(f.sessionClientConn)
	ucStats := f.usecaseFactory.GetStatsUsecase()
	ucPayToken := f.usecaseFactory.GetPayTokenUsecase()

	return map[int]app.Handler{
		INFO:                     info_handler.NewInfoHandler(f.logger, ucInfo),
		REGISTER:                 register_handler.NewRegisterHandler(f.logger, sManager, ucUser),
		LOGIN:                    login_handler.NewLoginHandler(f.logger, sManager, ucUser),
		LOGOUT:                   logout_handler.NewLogoutHandler(f.logger, sManager),
		PROFILE:                  profile_handler.NewProfileHandler(f.logger, sManager, ucUser),
		CREATORS:                 creator_handler.NewCreatorHandler(f.logger, sManager, ucCreator, ucUser),
		CREATOR_WITH_ID:          creator_id_handler.NewCreatorIdHandler(f.logger, sManager, ucCreator),
		SEARCH_CREATORS:          search_creators_handler.NewCreatorHandler(f.logger, sManager, ucCreator),
		UPDATE_PASSWORD:          password_handler.NewUpdatePasswordHandler(f.logger, sManager, ucUser),
		UPDATE_AVATAR:            avatar_handler.NewUpdateAvatarHandler(f.logger, sManager, ucUser),
		UPDATE_NICKNAME:          nickname_handler.NewUpdateNicknameHandler(f.logger, sManager, ucUser),
		AWARDS:                   aw_handler.NewAwardsHandler(f.logger, ucAwards, sManager),
		AWARDS_WITH_ID:           aw_id_handler.NewAwardsIdHandler(f.logger, ucAwards, sManager),
		AWARDS_UPDATE:            aw_upd_handler.NewAwardsUpdHandler(f.logger, ucAwards, sManager),
		POSTS:                    posts_handler.NewPostsHandler(f.logger, ucPosts, sManager),
		POSTS_WITH_ID:            posts_id_handler.NewPostsIDHandler(f.logger, ucPosts, ucUser, sManager),
		POSTS_UPD:                posts_upd_handler.NewPostsUpdateHandler(f.logger, ucPosts, sManager),
		POSTS_LIKES:              likes_handler.NewLikesHandler(f.logger, ucLikes, ucPosts, sManager),
		GET_CSRF_TOKEN:           csrf_handler.NewCsrfHandler(f.logger, sManager, ucCsrf),
		GET_USER_SUBSCRIPTIONS:   subscriptions_handler.NewSubscriptionsHandler(f.logger, sManager, ucSubscr),
		SUBSCRIBES:               subscribe_handler.NewSubscribeHandler(f.logger, sManager, ucSubscr),
		POST_UPD_COVER:           upl_cover_posts_handler.NewPostsUpdateCoverHandler(f.logger, ucPosts, sManager),
		ATTACH_ADD_TEXT:          upl_text_attach_handler.NewAttachesUploadTextHandler(f.logger, ucAttaches, ucPosts, sManager),
		ATTACH_ADD_IMAGE:         upl_img_attach_handler.NewPostsUploadImageHandler(f.logger, ucAttaches, ucPosts, sManager),
		ATTACH_ID:                attaches_id_handler.NewAttachesIDHandler(f.logger, ucAttaches, ucPosts, sManager),
		CREATOR_AVATAR:           upd_avatar_creator_handler.NewUpdateAvatarHandler(f.logger, sManager, ucCreator),
		CREATOR_COVER:            upd_cover_creator_handler.NewUpdateCoverHandler(f.logger, sManager, ucCreator),
		AWARDS_COVER:             upd_cover_awards_handler.NewUpdateCoverAwardsHandler(f.logger, sManager, ucAwards),
		ATTACH_UPD_IMAGE:         upd_img_data_handler.NewAttachUploadImageHandler(f.logger, ucAttaches, ucPosts, sManager),
		ATTACH_UPD_TEXT:          upd_text_data_handler.NewAttachesUpdateTextHandler(f.logger, ucAttaches, ucPosts, sManager),
		AWARDS_CREATOR_SUBSCRIBE: aw_subscribe_handler.NewAwardsSubscribeHandler(f.logger, sManager, ucSubscr, ucAwards, ucPayToken),
		USER_PAYMENTS:            payments_handler.NewPaymentsHandler(f.logger, sManager, ucPayments),
		ATTACHES:                 attaches_handler.NewAttachesHandler(f.logger, ucAttaches, ucPosts, sManager),
		ATTACH_ADD_VIDEO:         upl_video_attach_handler.NewPostsUploadVideoHandler(f.logger, ucAttaches, ucPosts, sManager),
		ATTACH_ADD_AUDIO:         upl_audio_attach_handler.NewPostsUploadAudioHandler(f.logger, ucAttaches, ucPosts, sManager),
		ATTACH_UPD_VIDEO:         upd_video_attach_handler.NewAttachUploadVideoHandler(f.logger, ucAttaches, ucPosts, sManager),
		ATTACH_UPD_AUDIO:         upd_audio_attach_handler.NewAttachUploadAudioHandler(f.logger, ucAttaches, ucPosts, sManager),
		POSTS_AVAILABLE:          user_posts_handler.NewPostsHandler(f.logger, sManager, ucPosts),
		STATS_COUNT_SUBSCRIBERS:  statistics_handler.NewCreatorCountSubscribersHandler(f.logger, ucStats),
		STATS_COUNT_POSTS:        statistics_count_posts_handler.NewCreatorCountPostsHandler(f.logger, ucStats),
		STATS_POSTS_VIEWS:        statistics_count_posts_views_handler.NewCreatorViewsHandler(f.logger, ucStats),
		STATS_TOTAL_INCOMES:      statistics_total_income_handler.NewCreatorTotalIncomeHandler(f.logger, ucStats),
		POST_COMMENTS:            comments_handler.NewCommentsHandler(f.logger, ucComment, ucPosts, sManager),
		COMMENTS_ID:              comments_id_handler.NewCommentsIdHandler(f.logger, ucComment, ucPosts, sManager),
		USER_COMMENTS:            user_comments_handler.NewUserCommentsHandler(f.logger, ucComment, sManager),
		USER_PAYMENTS_TOKEN:      pay_token_handler.NewTokenHandler(f.logger, sManager, ucPayToken, ucPayments),
		PAYMENTS_ACCOUNT:         pay_account_handler.NewAccountHandler(f.logger, ucPayToken),
		CREATOR_PAYMENTS:         creator_payments_handler.NewPaymentsHandler(f.logger, sManager, ucPayments),
	}
}

func (f *HandlerFactory) GetHandleUrls() *map[string]app.Handler {
	if f.urlHandler != nil {
		return f.urlHandler
	}

	hs := f.initAllHandlers()
	f.urlHandler = &map[string]app.Handler{
		//"/":                     "I am a joke?",
		"/login":    hs[LOGIN],
		"/info":     hs[INFO],
		"/logout":   hs[LOGOUT],
		"/register": hs[REGISTER],
		// /user     ---------------------------------------------------------////
		"/user":                  hs[PROFILE],
		"/user/update/password":  hs[UPDATE_PASSWORD],
		"/user/update/avatar":    hs[UPDATE_AVATAR],
		"/user/update/nickname":  hs[UPDATE_NICKNAME],
		"/user/subscriptions":    hs[GET_USER_SUBSCRIPTIONS],
		"/user/payments":         hs[USER_PAYMENTS],
		"/user/payments/token":   hs[USER_PAYMENTS_TOKEN],
		"/user/payments/account": hs[PAYMENTS_ACCOUNT],
		"/user/comments":         hs[USER_COMMENTS],
		"/user/posts":            hs[POSTS_AVAILABLE],
		// /creators ---------------------------------------------------------////
		"/creators":                                   hs[CREATORS],
		"/creators/{creator_id:[0-9]+}":               hs[CREATOR_WITH_ID],
		"/creators/{creator_id:[0-9]+}/subscribers":   hs[SUBSCRIBES],
		"/creators/{creator_id:[0-9]+}/update/avatar": hs[CREATOR_AVATAR],
		"/creators/{creator_id:[0-9]+}/update/cover":  hs[CREATOR_COVER],
		"/creators/{creator_id:[0-9]+}/payments":      hs[CREATOR_PAYMENTS],
		"/creators/search":                            hs[SEARCH_CREATORS],
		// ../awards ---------------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/awards":                                hs[AWARDS],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}":              hs[AWARDS_WITH_ID],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}/update":       hs[AWARDS_UPDATE],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}/update/cover": hs[AWARDS_COVER],
		"/creators/{creator_id:[0-9]+}/awards/{award_id:[0-9]+}/subscribe":    hs[AWARDS_CREATOR_SUBSCRIBE],
		// ../posts  ---------------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/posts":                               hs[POSTS],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}":              hs[POSTS_WITH_ID],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/update":       hs[POSTS_UPD],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/update/cover": hs[POST_UPD_COVER],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/like":         hs[POSTS_LIKES],
		// ../comments -----------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/comments":                     hs[POST_COMMENTS],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/comments/{comment_id:[0-9]+}": hs[COMMENTS_ID],
		// ../attaches  ----------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/attaches":                        hs[ATTACHES],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/attaches/text":                   hs[ATTACH_ADD_TEXT],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/attaches/image":                  hs[ATTACH_ADD_IMAGE],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/attaches/video":                  hs[ATTACH_ADD_VIDEO],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/attaches/audio":                  hs[ATTACH_ADD_AUDIO],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/{attach_id:[0-9]+}":              hs[ATTACH_ID],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/{attach_id:[0-9]+}/update/text":  hs[ATTACH_UPD_TEXT],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/{attach_id:[0-9]+}/update/image": hs[ATTACH_UPD_IMAGE],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/{attach_id:[0-9]+}/update/video": hs[ATTACH_UPD_VIDEO],
		"/creators/{creator_id:[0-9]+}/posts/{post_id:[0-9]+}/{attach_id:[0-9]+}/update/audio": hs[ATTACH_UPD_AUDIO],
		// ../statistics -----------------------------------------------------////
		"/creators/{creator_id:[0-9]+}/statistics/posts/views":  hs[STATS_POSTS_VIEWS],
		"/creators/{creator_id:[0-9]+}/statistics/posts/count":  hs[STATS_COUNT_POSTS],
		"/creators/{creator_id:[0-9]+}/statistics/total_income": hs[STATS_TOTAL_INCOMES],
		"/creators/{creator_id:[0-9]+}/statistics/subscribers":  hs[STATS_COUNT_SUBSCRIBERS],

		//   /token  ---------------------------------------------------------////
		"/token": hs[GET_CSRF_TOKEN],
	}
	return f.urlHandler
}

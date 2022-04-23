package usecase_likes

import (
	"glide/internal/app"
	"glide/internal/app/models"
	"glide/internal/app/repository"
	repoLikes "glide/internal/app/repository/likes"

	"github.com/pkg/errors"
)

type LikesUsecase struct {
	repository repoLikes.Repository
}

func NewLikesUsecase(repository repoLikes.Repository) *LikesUsecase {
	return &LikesUsecase{
		repository: repository,
	}
}

// Add Errors:
//		usecase_likes.IncorrectAddLike
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *LikesUsecase) Add(like *models.Like) (int64, error) {
	_, err := usecase.repository.GetLikeId(like.UserId, like.PostId)
	if err != nil {
		if errors.Is(err, repository.NotFound) {
			like.Value = 1
			return usecase.repository.Add(like)
		}
		return app.InvalidInt, err
	}
	return app.InvalidInt, IncorrectAddLike
}

// Delete Errors:
//		usecase_likes.IncorrectDelLike
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *LikesUsecase) Delete(postId int64, userId int64) (int64, error) {
	likeId, err := usecase.repository.GetLikeId(userId, postId)
	if err != nil {
		if errors.Is(err, repository.NotFound) {
			return app.InvalidInt, IncorrectDelLike
		}
		return app.InvalidInt, err
	}
	return usecase.repository.Delete(likeId)
}

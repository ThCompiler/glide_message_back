package usecase_factory

import (
	repFiles "glide/internal/app/repository/files"
	repUser "glide/internal/app/repository/user"
)

//go:generate mockgen -destination=mocks/mock_repository_factory.go -package=mock_repository_factory . RepositoryFactory

type RepositoryFactory interface {
	GetUserRepository() repUser.Repository
	GetFileRepository() repFiles.Repository
}

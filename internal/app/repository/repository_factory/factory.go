package repository_factory

import (
	"github.com/sirupsen/logrus"
	"glide/internal/app"
	repoFiles "glide/internal/app/repository/files"
	repository_os "glide/internal/app/repository/files/os"
	repUser "glide/internal/app/repository/user"
	repUserPsql "glide/internal/app/repository/user/postgresql"
)

type RepositoryFactory struct {
	expectedConnections app.ExpectedConnections
	logger              *logrus.Logger
	userRepository      repUser.Repository
	fileRepository      repoFiles.Repository
}

func NewRepositoryFactory(logger *logrus.Logger, expectedConnections app.ExpectedConnections) *RepositoryFactory {
	return &RepositoryFactory{
		expectedConnections: expectedConnections,
		logger:              logger,
	}
}

func (f *RepositoryFactory) GetUserRepository() repUser.Repository {
	if f.userRepository == nil {
		f.userRepository = repUserPsql.NewUserRepository(f.expectedConnections.SqlConnection)
	}
	return f.userRepository
}

func (f *RepositoryFactory) GetFileRepository() repoFiles.Repository {
	if f.fileRepository == nil {
		f.fileRepository = repository_os.NewFileRepository(f.expectedConnections.PathFiles)
	}
	return f.fileRepository
}

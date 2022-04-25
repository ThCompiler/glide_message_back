package repository_postgresql

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	postgresql_utilits "glide/internal/pkg/utilits/postgresql"
)

var (
	ChatAlreadyExists = errors.New("chat already exist")
)

const (
	codeForeignKeyVal = "23505"
	chatConstraint    = "unique_companions_of_chat"
)

func parsePQError(err *pq.Error) error {
	switch {
	case err.Code == codeForeignKeyVal && err.Constraint == chatConstraint:
		return ChatAlreadyExists
	default:
		return postgresql_utilits.NewDBError(err)
	}
}

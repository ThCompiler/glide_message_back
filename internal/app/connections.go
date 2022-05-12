package app

import (
	"github.com/jmoiron/sqlx"
	"glide/internal/pkg/rabbit"

	"google.golang.org/grpc"
)

type ExpectedConnections struct {
	SessionGrpcConnection *grpc.ClientConn
	SqlConnection         *sqlx.DB
	PathFiles             string
	RabbitSession         *rabbit.Session
}

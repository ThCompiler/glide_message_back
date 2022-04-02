package client

import (
	proto "glide/internal/microservices/auth/delivery/grpc/protobuf"
	"glide/internal/microservices/auth/sessions/models"
)

func ConvertAuthServerRespond(result *proto.Result) models.Result {
	if result == nil {
		return models.Result{}
	}
	res := models.Result{
		UserID: result.UserID,
		UniqID: result.SessionID,
	}
	return res
}

package server

import (
	"glide/internal/app"
)

type HandlerFactory interface {
	GetHandleUrls() *map[string]app.Handler
}

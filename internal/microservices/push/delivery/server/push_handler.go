package push_server

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	session_client "glide/internal/microservices/auth/delivery/grpc/client"
	"glide/internal/microservices/auth/sessions/middleware"
	"glide/internal/microservices/push/utils"
	hf "glide/internal/pkg/handler"
	"log"
	"net/http"
)

type PushHandler struct {
	sessionClient session_client.AuthCheckerClient
	hub           *utils.SendHub
	upgrader      *websocket.Upgrader
	hf.BaseHandler
}

func NewPushHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient, hub *utils.SendHub,
	upgrader *websocket.Upgrader) *PushHandler {
	h := &PushHandler{
		BaseHandler:   *hf.NewBaseHandler(log),
		sessionClient: sManager,
		hub:           hub,
		upgrader:      upgrader,
	}
	h.AddMethod(http.MethodGet, h.GET, middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc)
	return h
}

// GET create websocket with push
// @Summary create websocket with push
// @Description create websocket with send push about new comment or post, or subscriber
// @Produce json
// @tags utilities
// @Success 200 {object} utils.PushResponse Type can be "Comment", "Post", "Subscriber"
// @Success 200 {object} push_models.PostPush
// @Success 200 {object} push_models.CommentPush
// @Success 200 {object} push_models.SubPush
// @Failure 500 "server error"
// @Failure 401 "user are not authorized"
// @Router /user/push [GET]
func (h *PushHandler) GET(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Errorf("not found user_id in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := utils.NewClient(h.hub, userId, conn, h.Log(r))
	h.hub.RegisterClient(client)
	go client.SenderProcesses()
}

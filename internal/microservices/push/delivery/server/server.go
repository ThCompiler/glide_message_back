package push_server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"glide/internal/microservices/auth/delivery/grpc/client"
	"glide/internal/microservices/push"
	push_models "glide/internal/microservices/push/push"
	"glide/internal/microservices/push/push/repository"
	"glide/internal/microservices/push/push/usecase"
	"glide/internal/microservices/push/utils"
	"net/http"
	"time"

	"google.golang.org/grpc/connectivity"

	"glide/internal/app/middleware"

	"glide/internal/app"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config      *push.Config
	logger      *log.Logger
	connections app.ExpectedConnections
}

func New(config *push.Config, connections app.ExpectedConnections, logger *log.Logger) *Server {
	return &Server{
		config:      config,
		logger:      logger,
		connections: connections,
	}
}

func (s *Server) checkConnection() error {
	if err := s.connections.SqlConnection.Ping(); err != nil {
		return fmt.Errorf("Can't check connection to sql with error %v ", err)
	}

	s.logger.Info("Success check connection to sql db")

	state := s.connections.SessionGrpcConnection.GetState()
	if state != connectivity.Ready {
		return fmt.Errorf("Session connection not ready, status is: %s ", state)
	}

	if !s.connections.RabbitSession.CheckConnection() {
		return fmt.Errorf("Rabbit connection not ready ")
	}
	return nil
}

func (s *Server) Start() error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	router := mux.NewRouter()
	sManager := client.NewSessionClient(s.connections.SessionGrpcConnection)
	routerApi := router.PathPrefix("/api/v1/").Subrouter()

	senderHub := utils.NewHub()
	defer senderHub.StopHub()
	go senderHub.Run()

	h := NewPushHandler(s.logger, sManager, senderHub, upgrader)
	h.Connect(routerApi.Path("/push"))

	utilitsMiddleware := middleware.NewUtilitiesMiddleware(s.logger)
	routerApi.Use(utilitsMiddleware.CheckPanic, utilitsMiddleware.UpgradeLogger)

	cors := middleware.NewCorsMiddleware(&s.config.Cors, router)
	routerCors := cors.SetCors(router)

	pushUsecase := usecase.NewPushUsecase(repository.NewPushRepository(s.connections.SqlConnection))
	processingPush := utils.NewProcessingPush(s.logger.WithField("service", "push_proccessing"),
		s.connections.RabbitSession, senderHub, pushUsecase)

	defer processingPush.Stop()
	go processingPush.RunProcessMessage()
	go processingPush.RunProcessGlide()

	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(25 * time.Second)

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				keys := make([]string, len(h.hub.Clients))
				i := 0
				for k := range h.hub.Clients {
					keys[i] = k
					i++
				}
				h.hub.SendMessage(keys, &utils.PushResponse{Type: push.MessagePush, Push: &push_models.MessagePush{
					ChatId:          1,
					Companion:       "Человек",
					MessageId:       2,
					Text:            "Привет",
					CompanionAvatar: "tude",
				}})
			}
		}
	}()

	defer func() {
		done <- true
	}()
	s.logger.Info("start no production http server")
	return http.ListenAndServe(s.config.Port, routerCors)
}

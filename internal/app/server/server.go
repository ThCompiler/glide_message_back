package server

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/connectivity"

	"glide/internal/app"
	"glide/internal/app/delivery/http/handler_factory"
	"glide/internal/app/middleware"
	"glide/internal/app/repository/repository_factory"
	"glide/internal/app/usecase/usecase_factory"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config      *app.Config
	logger      *log.Logger
	connections app.ExpectedConnections
}

func New(config *app.Config, connections app.ExpectedConnections, logger *log.Logger) *Server {
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
		return fmt.Errorf("Push service connection not ready")
	}

	return nil
}

func (s *Server) Start(config *app.Config) error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	router := mux.NewRouter()

	routerApi := router.PathPrefix("/api/v1/").Subrouter()

	fileServer := http.FileServer(http.Dir(config.MediaDir + "/"))
	routerApi.PathPrefix("/" + app.LoadFileUrl).Handler(http.StripPrefix("/api/v1/"+app.LoadFileUrl, fileServer))

	repositoryFactory := repository_factory.NewRepositoryFactory(s.logger, s.connections)

	usecaseFactory := usecase_factory.NewUsecaseFactory(repositoryFactory)
	factory := handler_factory.NewFactory(s.logger, usecaseFactory, s.connections.SessionGrpcConnection)
	hs := factory.GetHandleUrls()

	for apiUrl, h := range *hs {
		h.Connect(routerApi.Path(apiUrl))
	}
	utilitsMiddleware := middleware.NewUtilitiesMiddleware(s.logger)
	routerApi.Use(utilitsMiddleware.CheckPanic, utilitsMiddleware.UpgradeLogger)

	cors := middleware.NewCorsMiddleware(&config.Cors, router)
	routerCors := cors.SetCors(router)

	s.logger.Info("start no production http server")
	return http.ListenAndServe(config.Port, routerCors)
}

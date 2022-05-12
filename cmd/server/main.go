package main

import (
	"flag"
	"glide/internal/app"
	main_server "glide/internal/app/server"
	"glide/internal/pkg/utilits"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
}

func main() {
	flag.Parse()
	logrus.Info(os.Args[:])

	config := app.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		logrus.Fatal(err)
	}

	logger, closeResource := utilits.NewLogger(&config.Config, false, "")

	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeResource, logger)

	repositoryConfig := &config.Repository

	db, closeResource := utilits.NewPostgresConnection(repositoryConfig.DataBaseUrl)

	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeResource, logger)

	rabbit, closeResource := utilits.NewRabbitSession(logger, repositoryConfig.RabbitUrl)

	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeResource, logger)

	sessionConn, err := utilits.NewGrpcConnection(config.Microservices.SessionServerUrl)
	if err != nil {
		logger.Fatal(err)
	}

	server := main_server.New(config,
		app.ExpectedConnections{
			SessionGrpcConnection: sessionConn,
			SqlConnection:         db,
			PathFiles:             config.MediaDir,
			RabbitSession:         rabbit,
		},
		logger,
	)

	if err = server.Start(config); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Server was stopped")
}

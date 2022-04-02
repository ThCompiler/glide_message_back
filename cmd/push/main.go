package main

import (
	"flag"
	"glide/internal/app"
	"glide/internal/microservices/push"
	push_server "glide/internal/microservices/push/delivery/server"
	"glide/internal/pkg/utilits"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/push.toml", "path to config file")
}

func main() {
	config := &push.Config{}

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		logrus.Fatal(err)
	}
	logger, CloseLogger := utilits.NewLogger(&config.Config, true, "push_microservice")
	defer CloseLogger()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		os.Exit(1)
	}
	logger.SetLevel(level)

	db, closeResource := utilits.NewPostgresConnection(config.SqlUrl)

	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeResource, logger)

	rabbit, closeResource := utilits.NewRabbitSession(logger, config.RabbitUrl)

	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeResource, logger)

	sessionConn, err := utilits.NewGrpcConnection(config.SessionUrl)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Push-service was start")
	server := push_server.New(config, app.ExpectedConnections{
		SessionGrpcConnection: sessionConn,
		RabbitSession:         rabbit,
		SqlConnection:         db,
	}, logger)
	if err = server.Start(); err != nil {
		logger.Fatalln(err)
	}
	logger.Info("Push-service was stopped")
}

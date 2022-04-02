package main

import (
	"flag"
	"os"
	_ "glide/docs"
	"glide/internal/app"
	main_server "glide/internal/app/server"
	"glide/internal/pkg/utilits"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	configPath          string
	useServerRepository bool
	runHttps            bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
	flag.BoolVar(&useServerRepository, "server-run", false, "true if it server run, false if it local run")
	flag.BoolVar(&runHttps, "run-https", false, "run https serve with certificates")

}

// @title Patreon
// @version 1.0
// @description Server for Patreon application.

// @tag.name user
// @tag.description "Some methods for work with user"

// @tag.name creators
// @tag.description "Some methods for work with creators"

// @tag.name attaches
// @tag.description "Some methods for work with attaches of post"

// @tag.name posts
// @tag.description "Some methods for work with posts"

// @tag.name awards
// @tag.description "Some methods for work with posts"

// @tag.name payments
// @tag.description "Some methods for work with payments"

// @tag.name comments
// @tag.description "Some methods for work with comments"

// @tag.name utilities
// @tag.description "Some methods for front work"

// @host api.pyaterochka-team.site
// @BasePath /api/v1

// @x-extension-openapi {"example": "value on a json format"}

func main() {
	flag.Parse()
	logrus.Info(os.Args[:])

	config := app.NewConfig()
	config.IsHTTPSServer = runHttps

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

	repositoryConfig := &config.LocalRepository
	if useServerRepository {
		repositoryConfig = &config.ServerRepository
	}

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
	filesConn, err := utilits.NewGrpcConnection(config.Microservices.FilesUrl)
	if err != nil {
		logger.Fatal(err)
	}
	server := main_server.New(config,
		app.ExpectedConnections{
			SessionGrpcConnection: sessionConn,
			FilesGrpcConnection:   filesConn,
			AccessRedisPool:       utilits.NewRedisPool(repositoryConfig.AccessRedisUrl),
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

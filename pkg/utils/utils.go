package utils

import (
	"fmt"
	"os"
	"patreon/internal"
	"glide/internal/pkg/rabbit"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"google.golang.org/grpc"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const MAX_GRPC_SIZE = 1024 * 1024 * 100

func NewLogger(config *internal.Config, isService bool, serviceName string) (log *logrus.Logger, closeResource func() error) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	currentTime := time.Now().In(time.UTC)
	var servicePath string
	if isService {
		servicePath = serviceName
	}
	formatted := config.LogAddr + fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second()) + "__" + servicePath + ".log"

	f, err := os.OpenFile(formatted, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	logger.SetOutput(f)
	logger.Writer()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger, f.Close
}

func NewPostgresConnection(databaseUrl string) (db *sqlx.DB, closeResource func() error) {
	db, err := sqlx.Open("postgres", databaseUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	return db, db.Close
}

func NewRabbitSession(logger *logrus.Logger, url string) (session *rabbit.Session, closeResource func() error) {
	session = rabbit.New(logger.WithField("service", "rabbit"), "rabbit", url)
	return session, session.Close
}

func NewRedisPool(redisUrl string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisUrl)
		},
	}
}

func NewGrpcConnection(grpcUrl string) (*grpc.ClientConn, error) {
	return grpc.Dial(grpcUrl, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MAX_GRPC_SIZE),
		grpc.MaxCallSendMsgSize(MAX_GRPC_SIZE)), grpc.WithBlock())
}

func StringsToLowerCase(array []string) []string {
	res := make([]string, len(array))
	for i, str := range array {
		res[i] = strings.ToLower(str)
	}
	return res
}

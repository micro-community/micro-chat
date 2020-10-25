package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"

	"github.com/micro/micro/v3/service/config"

	"github.com/micro-community/micro-chat/handler"
	_ "github.com/micro-community/micro-chat/profile"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("chat"),
		service.Version("latest"),
	)

	srv.Init(service.BeforeStart(func() error {
		redisHostValue, err := config.Get("Redis.Host")
		redisHost := redisHostValue.String("")

		if err != nil && redisHost != "" {
			logger.Info("config:", redisHost)
		}

		return nil
	}))

	srv.Handle(new(handler.Chat))

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

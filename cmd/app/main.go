package main

import (
	"context"
	"log"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/app"
	"github.com/SeiFlow-3P2/board_service/pkg/env"
)

func main() {
	if err := env.LoadEnv(); err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}

	cfg := &app.Config{
		AppName:      env.GetAppName(),
		Port:         env.GetPort(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		MongoURI:     env.GetMongoURI(),
		MongoDB:      env.GetMongoName(),
	}

	app := app.New(cfg)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

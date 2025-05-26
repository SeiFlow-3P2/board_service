package main

import (
	"context"
	"log"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/app"
	"github.com/SeiFlow-3P2/board_service/internal/config"
)

func main() {
	config.LoadEnv()

	cfg := &app.Config{
		Port:         "9090",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		MongoURI:     config.GetMongoURI(),
		MongoDB:      config.GetMongoDB(),
	}

	app := app.New(cfg)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

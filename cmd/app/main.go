package main

import (
	"context"
	"log"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/app"
)

func main() {
	cfg := &app.Config{
		Port:         "8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		MongoURI:     "mongodb://localhost:27017",
		MongoDB:      "board_service",
	}

	app := app.New(cfg)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

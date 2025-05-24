package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	application := app.New(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := application.Start(); err != nil {
			log.Printf("Error starting server: %v\n", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
	}

	log.Println("Server stopped")
}

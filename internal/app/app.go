package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/api"
	"github.com/SeiFlow-3P2/board_service/internal/config"
	"github.com/SeiFlow-3P2/board_service/internal/interceptor"
	"github.com/SeiFlow-3P2/board_service/internal/repository"
	"github.com/SeiFlow-3P2/board_service/internal/service"
	"github.com/SeiFlow-3P2/board_service/pkg/env"
	pb "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"github.com/SeiFlow-3P2/shared/kafka"
	"github.com/SeiFlow-3P2/shared/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	AppName      string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	MongoURI     string
	MongoDB      string
}

type App struct {
	config *Config
}

func New(cfg *Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Start(ctx context.Context) error {
	client, err := config.NewMongoClient(a.config.MongoURI)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	db := client.Database(a.config.MongoDB)

	shutdownTracer, err := telemetry.NewTracerProvider(
		ctx,
		a.config.AppName,
		env.GetOtelEndpoint(),
	)
	if err != nil {
		return fmt.Errorf("failed to create tracer provider: %w", err)
	}
	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}()

	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	boardService := service.NewBoardService(boardRepo)
	columnService := service.NewColumnService(columnRepo, boardRepo)

	p, err := kafka.NewProducer(
		env.GetKafkaBrokers(),
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %v", err)
	}
	defer p.Close()

	taskService := service.NewTaskService(taskRepo, columnRepo, p)

	boardServiceHandler := api.NewBoardServiceHandler(boardService)
	columnServiceHandler := api.NewColumnServiceHandler(columnService)
	taskServiceHandler := api.NewTaskServiceHandler(taskService)

	handler := api.NewHandler(boardServiceHandler, columnServiceHandler, taskServiceHandler)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthUnaryServerInterceptor()),
	)

	pb.RegisterBoardServiceServer(grpcServer, handler)

	reflection.Register(grpcServer)

	l, err := net.Listen("tcp", ":"+a.config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	serverError := make(chan error, 1)
	go func() {
		log.Printf("Starting gRPC server on port %s", a.config.Port)
		serverError <- grpcServer.Serve(l)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverError:
		return fmt.Errorf("grpc server error: %v", err)
	case <-shutdown:
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")
		return nil
	}
}

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/config"
	"github.com/SeiFlow-3P2/board_service/internal/repository"
	"github.com/SeiFlow-3P2/board_service/internal/service"
	"github.com/google/uuid"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	MongoURI     string
	MongoDB      string
}

type App struct {
	config        *Config
	server        *http.Server
	boardService  *service.BoardService
	columnService *service.ColumnService
	taskService   *service.TaskService
}

func New(cfg *Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Start() error {
	client, err := config.NewMongoClient(a.config.MongoURI)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database(a.config.MongoDB)

	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	a.boardService = service.NewBoardService(boardRepo)
	a.columnService = service.NewColumnService(columnRepo, boardRepo)
	a.taskService = service.NewTaskService(taskRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", a.healthCheckHandler)
	mux.HandleFunc("/", a.homeHandler)

	mux.HandleFunc("/api/boards", a.handleBoards)
	mux.HandleFunc("/api/boards/", a.handleBoardByID)

	mux.HandleFunc("/api/columns", a.handleColumns)
	mux.HandleFunc("/api/columns/", a.handleColumnByID)

	mux.HandleFunc("/api/tasks", a.handleTasks)
	mux.HandleFunc("/api/tasks/", a.handleTaskByID)

	a.server = &http.Server{
		Addr:         ":" + a.config.Port,
		Handler:      mux,
		ReadTimeout:  a.config.ReadTimeout,
		WriteTimeout: a.config.WriteTimeout,
		IdleTimeout:  a.config.IdleTimeout,
	}

	log.Printf("Server starting on port %s", a.config.Port)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to the Board Service")
}

func (a *App) handleBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, "userID", "test-user")

	switch r.Method {
	case http.MethodGet:
		boards, err := a.boardService.GetBoards(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(boards)
	case http.MethodPost:
		var input service.CreateBoardInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		board, err := a.boardService.CreateBoard(ctx, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(board)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) handleBoardByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "userID", "test-user")

	idStr := r.URL.Path[len("/api/boards/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		board, err := a.boardService.GetBoardInfo(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(board)
	case http.MethodPut:
		var input service.UpdateBoardInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		input.ID = id
		board, err := a.boardService.UpdateBoard(ctx, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(board)
	case http.MethodDelete:
		err := a.boardService.DeleteBoard(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) handleColumns(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (a *App) handleColumnByID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (a *App) handleTasks(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (a *App) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

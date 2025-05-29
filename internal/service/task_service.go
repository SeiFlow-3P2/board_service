package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/interceptor"
	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/SeiFlow-3P2/board_service/internal/repository"
	"github.com/SeiFlow-3P2/shared/kafka"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrEmptyTitle        = errors.New("task title cannot be empty")
	ErrEmptyDescription  = errors.New("task description cannot be empty")
	ErrInvalidDeadline   = errors.New("task deadline must be in the future")
	ErrEmptyID           = errors.New("ID cannot be empty")
	ErrTaskNotFound      = errors.New("task not found")
	ErrNewColumnNotFound = errors.New("new column not found")
	ErrGetColumnInfo     = errors.New("failed to get column info")
)

type TaskService struct {
	taskRepo   repository.TaskRepository
	columnRepo repository.ColumnRepository
	producer   *kafka.Producer
}

func NewTaskService(
	taskRepo repository.TaskRepository,
	columnRepo repository.ColumnRepository,
	producer *kafka.Producer,
) *TaskService {
	return &TaskService{
		taskRepo:   taskRepo,
		columnRepo: columnRepo,
		producer:   producer,
	}
}

type CreateTaskInput struct {
	Title       string
	Description string
	Deadline    *time.Time
	ColumnID    uuid.UUID
	InCalendar  bool
}

type MoveTaskInput struct {
	TaskID      uuid.UUID
	NewColumnID uuid.UUID
}

type UpdateTaskInput struct {
	TaskID      uuid.UUID
	Title       *string
	Description *string
	Deadline    *time.Time
}

type DeleteTaskInput struct {
	TaskID uuid.UUID
}

func (s *TaskService) CreateTask(ctx context.Context, input CreateTaskInput) (*models.Task, error) {
	userID, ok := ctx.Value(interceptor.UserIDKey).(string)
	if !ok {
		return nil, ErrUserNotInContext
	}

	task := &models.Task{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		Deadline:    *input.Deadline,
		Column_id:   input.ColumnID,
		In_Calendar: input.InCalendar,
	}

	task, err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	if input.InCalendar {
		go func() {
			msg := models.BoardEvent{
				EventType:   "create",
				Title:       task.Title,
				Description: task.Description,
				Deadline:    task.Deadline,
				UserID:      userID,
			}

			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Printf("failed to marshal message: %v", err)
				return
			}

			err = s.producer.Produce(
				ctx,
				string(jsonMsg),
				"board.event",
				userID,
				time.Second*10,
			)
			if err != nil {
				log.Printf("failed to produce message: %v", err)
				return
			}
		}()
	}

	fmt.Println("task", task)
	return task, nil
}

func (s *TaskService) MoveTask(ctx context.Context, input MoveTaskInput) (*models.Task, error) {

	_, err := s.taskRepo.GetTask(ctx, input.TaskID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	_, err = s.columnRepo.GetColumnInfo(ctx, input.NewColumnID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNewColumnNotFound
		}
		return nil, ErrGetColumnInfo
	}

	err = s.taskRepo.MoveTask(ctx, input.TaskID, input.NewColumnID)
	if err != nil {
		return nil, err
	}

	return s.taskRepo.GetTask(ctx, input.TaskID)
}

func (s *TaskService) UpdateTask(ctx context.Context, input UpdateTaskInput) (*models.Task, error) {

	_, err := s.taskRepo.GetTask(ctx, input.TaskID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	updates := &repository.TaskUpdates{
		Title:       input.Title,
		Description: input.Description,
		Deadline:    input.Deadline,
	}

	return s.taskRepo.UpdateTask(ctx, input.TaskID, updates)
}

func (s *TaskService) DeleteTask(ctx context.Context, input DeleteTaskInput) error {

	err := s.taskRepo.DeleteTask(ctx, input.TaskID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrTaskNotFound
		}
		return err
	}

	return nil
}

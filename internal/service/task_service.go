package service

import (
	"context"
	"errors"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/SeiFlow-3P2/board_service/internal/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrEmptyTitle       = errors.New("task title cannot be empty")
	ErrEmptyDescription = errors.New("task description cannot be empty")
	ErrInvalidDeadline  = errors.New("task deadline must be in the future")
	ErrEmptyID          = errors.New("ID cannot be empty")
	ErrTaskNotFound     = errors.New("task not found")
)

type TaskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

type CreateTaskInput struct {
	Title       string
	Description string
	Deadline    *time.Time
	ColumnID    uuid.UUID
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

	task := &models.Task{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		Deadline:    *input.Deadline,
		Column_id:   input.ColumnID,
		In_Calendar: false,
	}

	task, err := s.taskRepo.CreateTask(ctx, task)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) MoveTask(ctx context.Context, input MoveTaskInput) (*models.Task, error) {

	_, err := s.taskRepo.GetTask(ctx, input.TaskID.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	err = s.taskRepo.MoveTask(ctx, input.TaskID.String(), input.NewColumnID.String())
	if err != nil {
		return nil, err
	}

	return s.taskRepo.GetTask(ctx, input.TaskID.String())
}

func (s *TaskService) UpdateTask(ctx context.Context, input UpdateTaskInput) (*models.Task, error) {

	_, err := s.taskRepo.GetTask(ctx, input.TaskID.String())
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

	return s.taskRepo.UpdateTask(ctx, input.TaskID.String(), updates)
}

func (s *TaskService) DeleteTask(ctx context.Context, input DeleteTaskInput) error {

	err := s.taskRepo.DeleteTask(ctx, input.TaskID.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrTaskNotFound
		}
		return err
	}

	return nil
}

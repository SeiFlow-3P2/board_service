package api

import (
	"context"
	"strings"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/service"
	pb "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskServiceHandler struct {
	taskService *service.TaskService
}

func NewTaskServiceHandler(taskService *service.TaskService) *TaskServiceHandler {
	return &TaskServiceHandler{taskService: taskService}
}

func (h *TaskServiceHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if strings.TrimSpace(req.Description) == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}

	columnID, err := uuid.Parse(req.ColumnId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid column ID")
	}

	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid deadline format")
	}

	if deadline.Before(time.Now()) {
		return nil, status.Error(codes.InvalidArgument, "deadline must be in the future")
	}

	task, err := h.taskService.CreateTask(ctx, service.CreateTaskInput{
		Title:       req.Name,
		Description: req.Description,
		ColumnID:    columnID,
		Deadline:    &deadline,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TaskResponse{
		Id:          task.ID.String(),
		Name:        task.Title,
		Description: task.Description,
		Deadline:    task.Deadline.Format(time.RFC3339),
		InCalendar:  task.In_Calendar,
		ColumnId:    task.Column_id.String(),
	}, nil
}

func (h *TaskServiceHandler) MoveTask(ctx context.Context, req *pb.MoveTaskRequest) (*pb.MoveTaskResponse, error) {
	if req.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task ID is required")
	}

	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task ID")
	}

	newColumnID, err := uuid.Parse(req.NewColumnId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid column ID")
	}

	task, err := h.taskService.MoveTask(ctx, service.MoveTaskInput{
		TaskID:      taskID,
		NewColumnID: newColumnID,
	})
	if err != nil {
		if err == service.ErrTaskNotFound {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.MoveTaskResponse{
		TaskId:      task.ID.String(),
		NewColumnId: newColumnID.String(),
	}, nil
}

func (h *TaskServiceHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "task ID is required")
	}

	taskID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task ID")
	}

	if req.Name != nil {
		if strings.TrimSpace(req.Name.Value) == "" {
			return nil, status.Error(codes.InvalidArgument, "name is required")
		}
	}

	if req.Description != nil {
		if strings.TrimSpace(req.Description.Value) == "" {
			return nil, status.Error(codes.InvalidArgument, "description is required")
		}
	}

	var title, description *string
	if req.Name != nil {
		title = &req.Name.Value
	}
	if req.Description != nil {
		description = &req.Description.Value
	}
	task, err := h.taskService.UpdateTask(ctx, service.UpdateTaskInput{
		TaskID:      taskID,
		Title:       title,
		Description: description,
	})
	if err != nil {
		if err == service.ErrTaskNotFound {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TaskResponse{
		Id:          task.ID.String(),
		Name:        task.Title,
		Description: task.Description,
		Deadline:    task.Deadline.Format(time.RFC3339),
		InCalendar:  task.In_Calendar,
		ColumnId:    task.Column_id.String(),
	}, nil
}

func (h *TaskServiceHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "task ID is required")
	}

	taskID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task ID")
	}

	err = h.taskService.DeleteTask(ctx, service.DeleteTaskInput{
		TaskID: taskID,
	})
	if err != nil {
		if err == service.ErrTaskNotFound {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

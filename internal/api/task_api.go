package api

import (
	"context"
	"strings"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/service"
	pb "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"github.com/SeiFlow-3P2/shared/telemetry"
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
	ctx, span := telemetry.StartSpan(ctx, "TaskHandler.CreateTask")
	defer span.End()

	if strings.TrimSpace(req.Name) == "" {
		err := status.Error(codes.InvalidArgument, "name is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	if strings.TrimSpace(req.Description) == "" {
		err := status.Error(codes.InvalidArgument, "description is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	columnID, err := uuid.Parse(req.ColumnId)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid column ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid deadline format")
		telemetry.RecordError(span, err)
		return nil, err
	}

	if deadline.Before(time.Now()) {
		err := status.Error(codes.InvalidArgument, "deadline must be in the future")
		telemetry.RecordError(span, err)
		return nil, err
	}

	task, err := h.taskService.CreateTask(ctx, service.CreateTaskInput{
		Title:       req.Name,
		Description: req.Description,
		ColumnID:    columnID,
		Deadline:    &deadline,
		InCalendar:  req.InCalendar,
	})
	if err != nil {
		err := status.Error(codes.Internal, err.Error())
		telemetry.RecordError(span, err)
		return nil, err
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
	ctx, span := telemetry.StartSpan(ctx, "TaskHandler.MoveTask")
	defer span.End()

	if req.TaskId == "" {
		err := status.Error(codes.InvalidArgument, "task ID is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid task ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	newColumnID, err := uuid.Parse(req.NewColumnId)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid column ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	task, err := h.taskService.MoveTask(ctx, service.MoveTaskInput{
		TaskID:      taskID,
		NewColumnID: newColumnID,
	})
	if err != nil {
		switch {
		case err == service.ErrTaskNotFound:
			err := status.Error(codes.NotFound, "task not found")
			telemetry.RecordError(span, err)
			return nil, err
		case err.Error() == "new column not found":
			err := status.Error(codes.NotFound, "new column not found")
			telemetry.RecordError(span, err)
			return nil, err
		default:
			err := status.Error(codes.Internal, err.Error())
			telemetry.RecordError(span, err)
			return nil, err
		}
	}

	return &pb.MoveTaskResponse{
		TaskId:      task.ID.String(),
		NewColumnId: newColumnID.String(),
	}, nil
}

func (h *TaskServiceHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	ctx, span := telemetry.StartSpan(ctx, "TaskHandler.UpdateTask")
	defer span.End()

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "task ID is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	taskID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid task ID")
		telemetry.RecordError(span, err)
		return nil, err
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
			err := status.Error(codes.NotFound, "task not found")
			telemetry.RecordError(span, err)
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		telemetry.RecordError(span, err)
		return nil, err
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
	ctx, span := telemetry.StartSpan(ctx, "TaskHandler.DeleteTask")
	defer span.End()

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "task ID is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	taskID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid task ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	err = h.taskService.DeleteTask(ctx, service.DeleteTaskInput{
		TaskID: taskID,
	})
	if err != nil {
		if err == service.ErrTaskNotFound {
			err := status.Error(codes.NotFound, "task not found")
			telemetry.RecordError(span, err)
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		telemetry.RecordError(span, err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

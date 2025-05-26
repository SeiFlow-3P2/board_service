package api

import (
	"context"

	pb "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedBoardServiceServer
	boardHandler  *BoardServiceHandler
	columnHandler *ColumnServiceHandler
	taskHandler   *TaskServiceHandler
}

func NewHandler(
	boardHandler *BoardServiceHandler,
	columnHandler *ColumnServiceHandler,
	taskHandler *TaskServiceHandler,
) *Handler {
	return &Handler{
		boardHandler:  boardHandler,
		columnHandler: columnHandler,
		taskHandler:   taskHandler,
	}
}

// Board methods
func (h *Handler) CreateBoard(ctx context.Context, req *pb.CreateBoardRequest) (*pb.GetBoardInfoResponse, error) {
	return h.boardHandler.CreateBoard(ctx, req)
}

func (h *Handler) GetBoardInfo(ctx context.Context, req *pb.GetBoardInfoRequest) (*pb.GetBoardInfoResponse, error) {
	return h.boardHandler.GetBoardInfo(ctx, req)
}

func (h *Handler) GetBoards(ctx context.Context, req *pb.GetBoardsRequest) (*pb.BoardsListResponse, error) {
	return h.boardHandler.GetBoards(ctx, req)
}

func (h *Handler) UpdateBoard(ctx context.Context, req *pb.UpdateBoardRequest) (*pb.GetBoardInfoResponse, error) {
	return h.boardHandler.UpdateBoard(ctx, req)
}

func (h *Handler) DeleteBoard(ctx context.Context, req *pb.DeleteBoardRequest) (*emptypb.Empty, error) {
	return h.boardHandler.DeleteBoard(ctx, req)
}

// Column methods
func (h *Handler) CreateColumn(ctx context.Context, req *pb.CreateColumnRequest) (*pb.ColumnResponse, error) {
	return h.columnHandler.CreateColumn(ctx, req)
}

func (h *Handler) UpdateColumn(ctx context.Context, req *pb.UpdateColumnRequest) (*pb.ColumnResponse, error) {
	return h.columnHandler.UpdateColumn(ctx, req)
}

func (h *Handler) DeleteColumn(ctx context.Context, req *pb.DeleteColumnRequest) (*emptypb.Empty, error) {
	return h.columnHandler.DeleteColumn(ctx, req)
}

// Task methods
func (h *Handler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	return h.taskHandler.CreateTask(ctx, req)
}

func (h *Handler) MoveTask(ctx context.Context, req *pb.MoveTaskRequest) (*pb.MoveTaskResponse, error) {
	return h.taskHandler.MoveTask(ctx, req)
}

func (h *Handler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	return h.taskHandler.UpdateTask(ctx, req)
}

func (h *Handler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*emptypb.Empty, error) {
	return h.taskHandler.DeleteTask(ctx, req)
}

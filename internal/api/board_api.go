package api

import (
	"context"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/SeiFlow-3P2/board_service/internal/service"
	pb "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"github.com/SeiFlow-3P2/shared/telemetry"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BoardServiceHandler struct {
	boardService *service.BoardService
}

func NewBoardServiceHandler(boardService *service.BoardService) *BoardServiceHandler {
	return &BoardServiceHandler{boardService: boardService}
}

func (h *BoardServiceHandler) boardToGetInfoResponse(board *models.Board) *pb.GetBoardInfoResponse {
	var columns []*pb.ColumnInfo
	for _, col := range board.Columns {
		var tasks []*pb.TaskInfo
		for _, task := range col.Tasks {
			tasks = append(tasks, &pb.TaskInfo{
				Id:          task.ID.String(),
				Name:        task.Title,
				Description: task.Description,
				Deadline:    task.Deadline.Format(time.RFC3339),
				InCalendar:  task.In_Calendar,
				ColumnId:    task.Column_id.String(),
			})
		}

		columns = append(columns, &pb.ColumnInfo{
			Id:          col.ID.String(),
			Name:        col.Name,
			BoardId:     col.Desk_id.String(),
			OrderNumber: int64(col.Order_number),
			Tasks:       tasks,
		})
	}

	return &pb.GetBoardInfoResponse{
		Board: &pb.BoardInfo{
			Id:            board.ID.String(),
			Name:          board.Title,
			Description:   board.Description,
			Methodology:   board.Metodology,
			Category:      board.Category,
			Progress:      int64(board.Progress),
			Favorite:      board.Favorite,
			UpdatedAt:     timestamppb.New(board.Updated_at),
			CreatedAt:     timestamppb.New(board.Created_at),
			ColumnsAmount: int64(board.Columns_amount),
			UserId:        board.User_id,
			Columns:       columns,
		},
	}
}

func (h *BoardServiceHandler) CreateBoard(ctx context.Context, req *pb.CreateBoardRequest) (*pb.GetBoardInfoResponse, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardHandler.CreateBoard")
	defer span.End()

	if req.Name == "" {
		err := status.Error(codes.InvalidArgument, "title is required")
		telemetry.RecordError(span, err)
		return nil, err
	}
	if req.Description == "" {
		err := status.Error(codes.InvalidArgument, "description is required")
		telemetry.RecordError(span, err)
		return nil, err
	}
	if req.Methodology == "" {
		err := status.Error(codes.InvalidArgument, "methodology is required")
		telemetry.RecordError(span, err)
		return nil, err
	}
	if req.Category == "" {
		err := status.Error(codes.InvalidArgument, "category is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	if req.Methodology != "kanban" && req.Methodology != "simple" {
		err := status.Error(codes.InvalidArgument, "methotodology must be kanban or simple")
		telemetry.RecordError(span, err)
		return nil, err
	}

	board, err := h.boardService.CreateBoard(ctx, service.CreateBoardInput{
		Title:       req.Name,
		Description: req.Description,
		Metodology:  req.Methodology,
		Category:    req.Category,
	})
	if err != nil {
		switch {
		case err == service.ErrUserNotInContext:
			err := status.Error(codes.InvalidArgument, err.Error())
			telemetry.RecordError(span, err)
			return nil, err
		case err == service.ErrBoardExists:
			err := status.Error(codes.AlreadyExists, err.Error())
			telemetry.RecordError(span, err)
			return nil, err
		default:
			err := status.Error(codes.Internal, err.Error())
			telemetry.RecordError(span, err)
			return nil, err
		}
	}
	return h.boardToGetInfoResponse(board), nil
}

func (h *BoardServiceHandler) GetBoardInfo(ctx context.Context, req *pb.GetBoardInfoRequest) (*pb.GetBoardInfoResponse, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardHandler.GetBoardInfo")
	defer span.End()

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "board ID is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	boardID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid board ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	board, err := h.boardService.GetBoardInfo(ctx, boardID)
	if err != nil {
		if err == service.ErrBoardNotFound {
			err := status.Error(codes.NotFound, "board not found")
			telemetry.RecordError(span, err)
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		telemetry.RecordError(span, err)
		return nil, err
	}

	return h.boardToGetInfoResponse(board), nil
}

func (h *BoardServiceHandler) GetBoards(ctx context.Context, req *pb.GetBoardsRequest) (*pb.BoardsListResponse, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardHandler.GetBoards")
	defer span.End()

	boards, err := h.boardService.GetBoards(ctx)
	if err != nil {
		if err == service.ErrUserNotInContext {
			err := status.Error(codes.InvalidArgument, err.Error())
			telemetry.RecordError(span, err)
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		telemetry.RecordError(span, err)
		return nil, err
	}

	response := &pb.BoardsListResponse{
		Boards: make([]*pb.BoardResponse, 0, len(boards)),
	}

	for _, board := range boards {
		pbBoard := &pb.BoardResponse{
			Id:          board.ID.String(),
			Name:        board.Title,
			Description: board.Description,
			Methodology: board.Metodology,
			Category:    board.Category,
			Progress:    int64(board.Progress),
			Favorite:    board.Favorite,
			UpdatedAt:   timestamppb.New(board.Updated_at),
		}
		response.Boards = append(response.Boards, pbBoard)
	}

	return response, nil
}

func (h *BoardServiceHandler) UpdateBoard(ctx context.Context, req *pb.UpdateBoardRequest) (*pb.GetBoardInfoResponse, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardHandler.UpdateBoard")
	defer span.End()

	if req.Id == "" {
		err := status.Error(codes.InvalidArgument, "board ID is required")
		telemetry.RecordError(span, err)
		return nil, err
	}
	boardID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid board ID")
		telemetry.RecordError(span, err)
		return nil, err
	}
	if req.Name == nil && req.Description == nil && req.Progress == nil && req.Favorite == nil {
		err := status.Error(codes.InvalidArgument, "at least one field is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	var title *string
	if req.Name != nil {
		title = &req.Name.Value
	}

	var description *string
	if req.Description != nil {
		description = &req.Description.Value
	}

	var progress *int
	if req.Progress != nil {
		pbProgress := int(req.Progress.Value)
		progress = &pbProgress
	}

	var favorite *bool
	if req.Favorite != nil {
		favorite = &req.Favorite.Value
	}

	updates := service.UpdateBoardInput{
		ID:          boardID,
		Title:       title,
		Description: description,
		Progress:    progress,
		Favorite:    favorite,
	}

	board, err := h.boardService.UpdateBoard(ctx, updates)
	if err != nil {
		if err == service.ErrBoardNotFound {
			err := status.Error(codes.NotFound, "board not found")
			telemetry.RecordError(span, err)
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		telemetry.RecordError(span, err)
		return nil, err
	}

	return h.boardToGetInfoResponse(board), nil
}

func (h *BoardServiceHandler) DeleteBoard(ctx context.Context, req *pb.DeleteBoardRequest) (*emptypb.Empty, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardHandler.DeleteBoard")
	defer span.End()

	boardID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid board ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	err = h.boardService.DeleteBoard(ctx, boardID)
	if err != nil {
		if err == service.ErrBoardNotFound {
			err := status.Error(codes.NotFound, "board not found")
			telemetry.RecordError(span, err)
			return nil, err
		}
		err := status.Error(codes.Internal, err.Error())
		telemetry.RecordError(span, err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

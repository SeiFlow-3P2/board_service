package api

import (
	"context"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/SeiFlow-3P2/board_service/internal/service"
	pb "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
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
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description is required")
	}
	if req.Methodology == "" {
		return nil, status.Error(codes.InvalidArgument, "methodology is required")
	}
	if req.Category == "" {
		return nil, status.Error(codes.InvalidArgument, "category is required")
	}

	if req.Methodology != "kanban" && req.Methodology != "simple" {
		return nil, status.Error(codes.InvalidArgument, "methotodology must be kanban or simple")
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
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case err == service.ErrBoardExists:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
	}}
	return h.boardToGetInfoResponse(board), nil
}

func (h *BoardServiceHandler) GetBoardInfo(ctx context.Context, req *pb.GetBoardInfoRequest) (*pb.GetBoardInfoResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "board ID is required")
	}

	boardID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid board ID")
	}

	board, err := h.boardService.GetBoardInfo(ctx, boardID)
	if err != nil {
		if err == service.ErrBoardNotFound {
			return nil, status.Error(codes.NotFound, "board not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.boardToGetInfoResponse(board), nil
}

func (h *BoardServiceHandler) GetBoards(ctx context.Context, req *pb.GetBoardsRequest) (*pb.BoardsListResponse, error) {
	boards, err := h.boardService.GetBoards(ctx)
	if err != nil {
		if err == service.ErrUserNotInContext {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
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
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "board ID is required")
	}
	boardID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid board ID")
	}
	if req.Name == nil && req.Description == nil && req.Progress == nil && req.Favorite == nil {
		return nil, status.Error(codes.InvalidArgument, "at least one field is required")
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
			return nil, status.Error(codes.NotFound, "board not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.boardToGetInfoResponse(board), nil
}

func (h *BoardServiceHandler) DeleteBoard(ctx context.Context, req *pb.DeleteBoardRequest) (*emptypb.Empty, error) {
	boardID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid board ID")
	}

	err = h.boardService.DeleteBoard(ctx, boardID)
	if err != nil {
		if err == service.ErrBoardNotFound {
			return nil, status.Error(codes.NotFound, "board not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

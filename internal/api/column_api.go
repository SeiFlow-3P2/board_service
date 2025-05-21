package api

import (
	"context"

	"github.com/SeiFlow-3P2/board_service/internal/service"
	pb "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ColumnServiceHandler struct {
	columnService service.ColumnService
}

func NewColumnServiceHandler(columnService service.ColumnService) *ColumnServiceHandler {
	return &ColumnServiceHandler{columnService: columnService}
}

func (h *ColumnServiceHandler) CreateColumn(ctx context.Context, req *pb.CreateColumnRequest) (*pb.ColumnResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	boardID, err := uuid.Parse(req.BoardId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid board ID")
	}

	column, err := h.columnService.CreateColumn(ctx, service.CreateColumnInput{
		Name:   req.Name,
		DeskID: boardID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ColumnResponse{
		Id:          column.ID.String(),
		Name:        column.Name,
		BoardId:     column.Desk_id.String(),
		OrderNumber: int64(column.Order_number),
	}, nil
}

func (h *ColumnServiceHandler) UpdateColumn(ctx context.Context, req *pb.UpdateColumnRequest) (*pb.ColumnResponse, error) {
	if req.Name == nil || req.Name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	columnID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid column ID")
	}

	column, err := h.columnService.UpdateColumn(ctx, service.UpdateColumnInput{
		ID:          columnID,
		Name:        &req.Name.Value,
		OrderNumber: nil,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ColumnResponse{
		Id:          column.ID.String(),
		Name:        column.Name,
		BoardId:     column.Desk_id.String(),
		OrderNumber: int64(column.Order_number),
	}, nil
}

func (h *ColumnServiceHandler) DeleteColumn(ctx context.Context, req *pb.DeleteColumnRequest) (*emptypb.Empty, error) {
	columnID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid column ID")
	}

	err = h.columnService.DeleteColumn(ctx, service.DeleteColumnInput{
		ID: columnID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

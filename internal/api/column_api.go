package api

import (
	"context"
	"fmt"

	"github.com/SeiFlow-3P2/board_service/internal/service"
	pb "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"github.com/SeiFlow-3P2/shared/telemetry"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ColumnServiceHandler struct {
	columnService *service.ColumnService
}

func NewColumnServiceHandler(columnService *service.ColumnService) *ColumnServiceHandler {
	return &ColumnServiceHandler{columnService: columnService}
}

func (h *ColumnServiceHandler) CreateColumn(ctx context.Context, req *pb.CreateColumnRequest) (*pb.ColumnResponse, error) {
	ctx, span := telemetry.StartSpan(ctx, "ColumnHandler.CreateColumn")
	defer span.End()

	fmt.Println("CreateColumnRequest", req)
	if req.Name == "" {
		err := status.Error(codes.InvalidArgument, "name is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	boardID, err := uuid.Parse(req.BoardId)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid board ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	column, err := h.columnService.CreateColumn(ctx, service.CreateColumnInput{
		Name:   req.Name,
		DeskID: boardID,
	})
	if err != nil {
		switch {
		case err.Error() == "board not found":
			err := status.Error(codes.NotFound, "board not found")
			telemetry.RecordError(span, err)
			return nil, err
		case err.Error() == service.ErrColumnExists.Error():
			err := status.Error(codes.AlreadyExists, "column with this name already exists")
			telemetry.RecordError(span, err)
			return nil, err
		default:
			err := status.Error(codes.Internal, err.Error())
			telemetry.RecordError(span, err)
			return nil, err
		}
	}

	return &pb.ColumnResponse{
		Id:          column.ID.String(),
		Name:        column.Name,
		BoardId:     column.Desk_id.String(),
		OrderNumber: int64(column.Order_number),
	}, nil
}

func (h *ColumnServiceHandler) UpdateColumn(ctx context.Context, req *pb.UpdateColumnRequest) (*pb.ColumnResponse, error) {
	ctx, span := telemetry.StartSpan(ctx, "ColumnHandler.UpdateColumn")
	defer span.End()

	if req.Name == nil || req.Name.Value == "" {
		err := status.Error(codes.InvalidArgument, "name is required")
		telemetry.RecordError(span, err)
		return nil, err
	}

	columnID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid column ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	column, err := h.columnService.UpdateColumn(ctx, service.UpdateColumnInput{
		ID:          columnID,
		Name:        &req.Name.Value,
		OrderNumber: nil,
	})
	if err != nil {
		switch {
		case err.Error() == "column not found":
			err := status.Error(codes.NotFound, "column not found")
			telemetry.RecordError(span, err)
			return nil, err
		case err.Error() == service.ErrColumnExists.Error():
			err := status.Error(codes.AlreadyExists, "column with this name already exists")
			telemetry.RecordError(span, err)
			return nil, err
		default:
			err := status.Error(codes.Internal, err.Error())
			telemetry.RecordError(span, err)
			return nil, err
		}
	}

	return &pb.ColumnResponse{
		Id:          column.ID.String(),
		Name:        column.Name,
		BoardId:     column.Desk_id.String(),
		OrderNumber: int64(column.Order_number),
	}, nil
}

func (h *ColumnServiceHandler) DeleteColumn(ctx context.Context, req *pb.DeleteColumnRequest) (*emptypb.Empty, error) {
	ctx, span := telemetry.StartSpan(ctx, "ColumnHandler.DeleteColumn")
	defer span.End()

	columnID, err := uuid.Parse(req.Id)
	if err != nil {
		err := status.Error(codes.InvalidArgument, "invalid column ID")
		telemetry.RecordError(span, err)
		return nil, err
	}

	err = h.columnService.DeleteColumn(ctx, service.DeleteColumnInput{
		ID: columnID,
	})
	if err != nil {
		switch {
		case err.Error() == "column not found":
			err := status.Error(codes.NotFound, "column not found")
			telemetry.RecordError(span, err)
			return nil, err
		default:
			err := status.Error(codes.Internal, err.Error())
			telemetry.RecordError(span, err)
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

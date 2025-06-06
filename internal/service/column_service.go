package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/SeiFlow-3P2/board_service/internal/repository"
	"github.com/SeiFlow-3P2/shared/telemetry"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrEmptyName        = errors.New("column name cannot be empty")
	ErrEmptyDeskID      = errors.New("desk ID cannot be empty")
	ErrEmptyOrderNumber = errors.New("order number cannot be empty")
	ErrColumnExists     = errors.New("column with this name already exists in the board")
	ErrColumnNotFound   = errors.New("column not found")
)

type ColumnService struct {
	columnRepo repository.ColumnRepository
	boardRepo  repository.BoardRepository
}

func NewColumnService(columnRepo repository.ColumnRepository, boardRepo repository.BoardRepository) *ColumnService {
	return &ColumnService{
		columnRepo: columnRepo,
		boardRepo:  boardRepo,
	}
}

type CreateColumnInput struct {
	Name        string
	DeskID      uuid.UUID
	OrderNumber int
}

type UpdateColumnInput struct {
	ID          uuid.UUID
	Name        *string
	OrderNumber *int
}

type DeleteColumnInput struct {
	ID     uuid.UUID
	DeskID uuid.UUID
}

func (s *ColumnService) CreateColumn(ctx context.Context, input CreateColumnInput) (*models.Column, error) {
	ctx, span := telemetry.StartSpan(ctx, "ColumnService.CreateColumn")
	defer span.End()

	_, err := s.boardRepo.GetBoardInfo(ctx, input.DeskID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			telemetry.RecordError(span, ErrBoardNotFound)
			return nil, ErrBoardNotFound
		}
		telemetry.RecordError(span, err)
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	existColumns, err := s.columnRepo.GetColumns(ctx, input.DeskID)
	if err != nil {
		telemetry.RecordError(span, err)
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	for _, col := range existColumns {
		if strings.EqualFold(col.Name, input.Name) {
			telemetry.RecordError(span, ErrColumnExists)
			return nil, ErrColumnExists
		}
	}

	newOrderNumber, err := s.boardRepo.IncrementColumnsAmount(ctx, input.DeskID)
	if err != nil {
		telemetry.RecordError(span, err)
		return nil, fmt.Errorf("failed to increment columns amount: %w", err)
	}

	column := &models.Column{
		ID:           uuid.New(),
		Name:         input.Name,
		Desk_id:      input.DeskID,
		Order_number: newOrderNumber,
	}

	column, err = s.columnRepo.CreateColumn(ctx, column)
	if err != nil {
		_ = s.boardRepo.DecrementColumnsAmount(ctx, input.DeskID)
		telemetry.RecordError(span, err)
		return nil, fmt.Errorf("failed to create column: %w", err)
	}

	return column, nil
}

func (s *ColumnService) UpdateColumn(ctx context.Context, input UpdateColumnInput) (*models.Column, error) {
	ctx, span := telemetry.StartSpan(ctx, "ColumnService.UpdateColumn")
	defer span.End()

	column, err := s.columnRepo.GetColumnInfo(ctx, input.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			telemetry.RecordError(span, ErrColumnNotFound)
			return nil, ErrColumnNotFound
		}
		telemetry.RecordError(span, err)
		return nil, fmt.Errorf("failed to get column info: %w", err)
	}

	updates := &repository.ColumnUpdates{}

	if input.Name != nil {
		existColumns, err := s.columnRepo.GetColumns(ctx, column.Desk_id)
		if err != nil {
			telemetry.RecordError(span, err)
			return nil, fmt.Errorf("failed to get columns: %w", err)
		}

		for _, col := range existColumns {
			if col.ID != column.ID && strings.EqualFold(col.Name, *input.Name) {
				telemetry.RecordError(span, ErrColumnExists)
				return nil, ErrColumnExists
			}
		}
		updates.Name = input.Name
	}

	return s.columnRepo.UpdateColumn(ctx, input.ID, updates)
}

func (s *ColumnService) DeleteColumn(ctx context.Context, input DeleteColumnInput) error {
	ctx, span := telemetry.StartSpan(ctx, "ColumnService.DeleteColumn")
	defer span.End()

	if input.ID == uuid.Nil {
		telemetry.RecordError(span, ErrEmptyID)
		return ErrEmptyID
	}

	column, err := s.columnRepo.GetColumnInfo(ctx, input.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			telemetry.RecordError(span, ErrColumnNotFound)
			return ErrColumnNotFound
		}
		telemetry.RecordError(span, err)
		return fmt.Errorf("failed to get column info: %w", err)
	}

	err = s.columnRepo.DeleteColumn(ctx, input.ID)
	if err != nil {
		telemetry.RecordError(span, err)
		return fmt.Errorf("failed to delete column: %w", err)
	}

	err = s.columnRepo.DecrementOrderNumbers(ctx, column.Desk_id, column.Order_number)
	if err != nil {
		telemetry.RecordError(span, err)
		return fmt.Errorf("failed to decrement order numbers: %w", err)
	}

	err = s.boardRepo.DecrementColumnsAmount(ctx, column.Desk_id)
	if err != nil {
		telemetry.RecordError(span, err)
		return fmt.Errorf("failed to decrement columns amount: %w", err)
	}

	return nil
}

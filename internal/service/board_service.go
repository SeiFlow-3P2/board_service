package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/interceptor"
	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/SeiFlow-3P2/board_service/internal/repository"
	"github.com/SeiFlow-3P2/shared/telemetry"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrBoardExists      = errors.New("board already exists")
	ErrUserNotInContext = errors.New("user ID not found in context")
	ErrBoardNotFound    = errors.New("board not found")
)

type BoardService struct {
	boardRepo repository.BoardRepository
}

func NewBoardService(boardRepo repository.BoardRepository) *BoardService {
	return &BoardService{
		boardRepo: boardRepo,
	}
}

type CreateBoardInput struct {
	Title       string
	Description string
	Metodology  string
	Category    string
}

type UpdateBoardInput struct {
	ID          uuid.UUID
	Title       *string
	Description *string
	Progress    *int
	Favorite    *bool
}

func (s *BoardService) CreateBoard(ctx context.Context, input CreateBoardInput) (*models.Board, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardService.CreateBoard")
	defer span.End()

	userID, ok := ctx.Value(interceptor.UserIDKey).(string)
	if !ok {
		telemetry.RecordError(span, ErrUserNotInContext)
		return nil, ErrUserNotInContext
	}

	boards, err := s.boardRepo.GetBoards(ctx, userID)
	if err != nil {
		telemetry.RecordError(span, err)
		return nil, err
	}
	for _, board := range boards {
		if board.Title == input.Title {
			telemetry.RecordError(span, ErrBoardExists)
			return nil, ErrBoardExists
		}
	}

	var columns []models.Column
	columnsAmount := 0

	now := time.Now()
	boardID := uuid.New()

	switch input.Metodology {
	case "kanban":
		columnsAmount = 3
		columns = []models.Column{
			{
				ID:           uuid.New(),
				Name:         "To Do",
				Order_number: 1,
				Desk_id:      boardID,
				Tasks:        []models.Task{},
			},
			{
				ID:           uuid.New(),
				Name:         "In Progress",
				Order_number: 2,
				Desk_id:      boardID,
				Tasks:        []models.Task{},
			},
			{
				ID:           uuid.New(),
				Name:         "Done",
				Order_number: 3,
				Desk_id:      boardID,
				Tasks:        []models.Task{},
			},
		}
	case "simple":
		columnsAmount = 1
		columns = []models.Column{
			{
				ID:           uuid.New(),
				Name:         "",
				Order_number: 1,
				Desk_id:      boardID,
				Tasks:        []models.Task{},
			},
		}
	}

	board := &models.Board{
		ID:             boardID,
		Title:          input.Title,
		Description:    input.Description,
		Category:       input.Category,
		Progress:       0,
		Favorite:       false,
		Metodology:     input.Metodology,
		Columns_amount: columnsAmount,
		Created_at:     now,
		Updated_at:     now,
		User_id:        userID,
		Columns:        columns,
	}

	createdBoard, err := s.boardRepo.CreateBoard(ctx, board)
	if err != nil {
		return nil, err
	}
	return createdBoard, nil
}

func (s *BoardService) GetBoardInfo(ctx context.Context, id uuid.UUID) (*models.Board, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardService.GetBoardInfo")
	defer span.End()

	board, err := s.boardRepo.GetBoardInfo(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			telemetry.RecordError(span, ErrBoardNotFound)
			return nil, ErrBoardNotFound
		}
		return nil, err
	}
	return board, nil
}

func (s *BoardService) GetBoards(ctx context.Context) ([]*models.Board, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardService.GetBoards")
	defer span.End()

	userID, ok := ctx.Value(interceptor.UserIDKey).(string)
	if !ok {
		return nil, ErrUserNotInContext
	}

	boards, err := s.boardRepo.GetBoards(ctx, userID)
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (s *BoardService) UpdateBoard(ctx context.Context, input UpdateBoardInput) (*models.Board, error) {
	ctx, span := telemetry.StartSpan(ctx, "BoardService.UpdateBoard")
	defer span.End()

	board, err := s.boardRepo.GetBoardInfo(ctx, input.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			telemetry.RecordError(span, ErrBoardNotFound)
			return nil, ErrBoardNotFound
		}
		return nil, err
	}

	updates := &repository.BoardUpdates{}
	now := time.Now()

	if input.Title != nil {
		existBoards, err := s.boardRepo.GetBoards(ctx, board.User_id)
		if err != nil {
			telemetry.RecordError(span, err)
			return nil, err
		}
		for _, b := range existBoards {
			if b.ID != input.ID && strings.EqualFold(b.Title, *input.Title) {
				telemetry.RecordError(span, ErrBoardExists)
				return nil, ErrBoardExists
			}
		}
		updates.Title = input.Title
	}
	updates.Description = input.Description
	updates.Progress = input.Progress
	updates.Favorite = input.Favorite
	updates.Updated_at = &now

	return s.boardRepo.UpdateBoard(ctx, input.ID, updates)
}

func (s *BoardService) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BoardService.DeleteBoard")
	defer span.End()

	_, err := s.boardRepo.GetBoardInfo(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			telemetry.RecordError(span, ErrBoardNotFound)
			return ErrBoardNotFound
		}
		return err
	}

	err = s.boardRepo.DeleteBoard(ctx, id)
	if err != nil {
		telemetry.RecordError(span, err)
		return err
	}

	return nil
}

package board_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	boardpb "board_service/pkg/proto/board/v1"
)

// mockBoardService implements boardpb.BoardServiceServer for testing
type mockBoardService struct {
	boardpb.UnimplementedBoardServiceServer
}

func (m *mockBoardService) CreateBoard(ctx context.Context, req *boardpb.CreateBoardRequest) (*boardpb.GetBoardInfoResponse, error) {
	if req.GetName() == "" {
		return nil, errors.New("board name is required")
	}

	now := timestamppb.Now()
	return &boardpb.GetBoardInfoResponse{
		Board: &boardpb.BoardInfo{
			Id:          "board-123",
			Name:        req.GetName(),
			Description: req.GetDescription(),
			Methodology: req.GetMethodology(),
			Category:    req.GetCategory(),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, nil
}

func (m *mockBoardService) GetBoards(ctx context.Context, req *boardpb.GetBoardsRequest) (*boardpb.BoardsListResponse, error) {
	now := timestamppb.Now()
	return &boardpb.BoardsListResponse{
		Boards: []*boardpb.BoardResponse{
			{
				Id:          "board-1",
				Name:        "Test Board 1",
				Description: "Description 1",
				UpdatedAt:   now,
			},
			{
				Id:          "board-2",
				Name:        "Test Board 2",
				Description: "Description 2",
				UpdatedAt:   now,
			},
		},
	}, nil
}

func (m *mockBoardService) GetBoardInfo(ctx context.Context, req *boardpb.GetBoardInfoRequest) (*boardpb.GetBoardInfoResponse, error) {
	if req.GetId() == "" {
		return nil, errors.New("board ID is required")
	}

	now := timestamppb.Now()
	return &boardpb.GetBoardInfoResponse{
		Board: &boardpb.BoardInfo{
			Id:          req.GetId(),
			Name:        "Test Board",
			Description: "Test Description",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}, nil
}

func (m *mockBoardService) UpdateBoard(ctx context.Context, req *boardpb.UpdateBoardRequest) (*boardpb.GetBoardInfoResponse, error) {
	if req.GetId() == "" {
		return nil, errors.New("board ID is required")
	}

	now := timestamppb.Now()
	board := &boardpb.BoardInfo{
		Id:        req.GetId(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if name := req.GetName(); name != nil {
		board.Name = name.GetValue()
	}
	if desc := req.GetDescription(); desc != nil {
		board.Description = desc.GetValue()
	}
	if progress := req.GetProgress(); progress != nil {
		board.Progress = int64(progress.GetValue())
	}
	if fav := req.GetFavorite(); fav != nil {
		board.Favorite = fav.GetValue()
	}

	return &boardpb.GetBoardInfoResponse{Board: board}, nil
}

func (m *mockBoardService) DeleteBoard(ctx context.Context, req *boardpb.DeleteBoardRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New("board ID is required")
	}
	return &emptypb.Empty{}, nil
}

func (m *mockBoardService) CreateColumn(ctx context.Context, req *boardpb.CreateColumnRequest) (*boardpb.ColumnResponse, error) {
	if req.GetName() == "" || req.GetBoardId() == "" {
		return nil, errors.New("column name and board ID are required")
	}

	return &boardpb.ColumnResponse{
		Id:        "column-123",
		Name:      req.GetName(),
		BoardId:   req.GetBoardId(),
		OrderNumber: 1,
	}, nil
}

func (m *mockBoardService) UpdateColumn(ctx context.Context, req *boardpb.UpdateColumnRequest) (*boardpb.ColumnResponse, error) {
	if req.GetId() == "" {
		return nil, errors.New("column ID is required")
	}

	resp := &boardpb.ColumnResponse{
		Id: req.GetId(),
	}
	if name := req.GetName(); name != nil {
		resp.Name = name.GetValue()
	}
	return resp, nil
}

func (m *mockBoardService) DeleteColumn(ctx context.Context, req *boardpb.DeleteColumnRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New("column ID is required")
	}
	return &emptypb.Empty{}, nil
}

func (m *mockBoardService) CreateTask(ctx context.Context, req *boardpb.CreateTaskRequest) (*boardpb.TaskResponse, error) {
	if req.GetName() == "" || req.GetColumnId() == "" {
		return nil, errors.New("task name and column ID are required")
	}

	return &boardpb.TaskResponse{
		Id:          "task-123",
		Name:        req.GetName(),
		Description: req.GetDescription(),
		ColumnId:    req.GetColumnId(),
	}, nil
}

func (m *mockBoardService) MoveTask(ctx context.Context, req *boardpb.MoveTaskRequest) (*boardpb.MoveTaskResponse, error) {
	if req.GetTaskId() == "" || req.GetNewColumnId() == "" {
		return nil, errors.New("task ID and new column ID are required")
	}

	return &boardpb.MoveTaskResponse{
		TaskId:      req.GetTaskId(),
		NewColumnId: req.GetNewColumnId(),
	}, nil
}

func (m *mockBoardService) UpdateTask(ctx context.Context, req *boardpb.UpdateTaskRequest) (*boardpb.TaskResponse, error) {
	if req.GetId() == "" {
		return nil, errors.New("task ID is required")
	}

	resp := &boardpb.TaskResponse{
		Id: req.GetId(),
	}
	if name := req.GetName(); name != nil {
		resp.Name = name.GetValue()
	}
	if desc := req.GetDescription(); desc != nil {
		resp.Description = desc.GetValue()
	}
	return resp, nil
}

func (m *mockBoardService) DeleteTask(ctx context.Context, req *boardpb.DeleteTaskRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New("task ID is required")
	}
	return &emptypb.Empty{}, nil
}

func TestBoardService(t *testing.T) {
	ctx := context.Background()
	service := &mockBoardService{}

	t.Run("BoardOperations", func(t *testing.T) {
		t.Run("CreateBoard", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.CreateBoardRequest
				wantErr bool
			}{
				{
					name: "success",
					req: &boardpb.CreateBoardRequest{
						Name:        "Test Board",
						Description: "Test Description",
					},
					wantErr: false,
				},
				{
					name:    "missing name",
					req:     &boardpb.CreateBoardRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.CreateBoard(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("CreateBoard() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && resp.GetBoard().GetName() != tt.req.GetName() {
						t.Errorf("expected board name %q, got %q", tt.req.GetName(), resp.GetBoard().GetName())
					}
				})
			}
		})

		t.Run("GetBoards", func(t *testing.T) {
			resp, err := service.GetBoards(ctx, &boardpb.GetBoardsRequest{})
			if err != nil {
				t.Fatalf("GetBoards() error = %v", err)
			}
			if len(resp.GetBoards()) != 2 {
				t.Errorf("expected 2 boards, got %d", len(resp.GetBoards()))
			}
		})

		t.Run("GetBoardInfo", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.GetBoardInfoRequest
				wantErr bool
			}{
				{
					name: "success",
					req:  &boardpb.GetBoardInfoRequest{Id: "board-123"},
				},
				{
					name:    "missing id",
					req:     &boardpb.GetBoardInfoRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.GetBoardInfo(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("GetBoardInfo() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && resp.GetBoard().GetId() != tt.req.GetId() {
						t.Errorf("expected board ID %q, got %q", tt.req.GetId(), resp.GetBoard().GetId())
					}
				})
			}
		})

		t.Run("UpdateBoard", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.UpdateBoardRequest
				wantErr bool
			}{
				{
					name: "update name",
					req: &boardpb.UpdateBoardRequest{
						Id:   "board-123",
						Name: wrapperspb.String("New Name"),
					},
				},
				{
					name:    "missing id",
					req:     &boardpb.UpdateBoardRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.UpdateBoard(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("UpdateBoard() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && name := tt.req.GetName(); name != nil && resp.GetBoard().GetName() != name.GetValue() {
						t.Errorf("expected board name %q, got %q", name.GetValue(), resp.GetBoard().GetName())
					}
				})
			}
		})

		t.Run("DeleteBoard", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.DeleteBoardRequest
				wantErr bool
			}{
				{
					name: "success",
					req:  &boardpb.DeleteBoardRequest{Id: "board-123"},
				},
				{
					name:    "missing id",
					req:     &boardpb.DeleteBoardRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					_, err := service.DeleteBoard(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("DeleteBoard() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
			}
		})
	})

	t.Run("ColumnOperations", func(t *testing.T) {
		t.Run("CreateColumn", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.CreateColumnRequest
				wantErr bool
			}{
				{
					name: "success",
					req: &boardpb.CreateColumnRequest{
						Name:    "To Do",
						BoardId: "board-123",
					},
				},
				{
					name:    "missing name",
					req:     &boardpb.CreateColumnRequest{BoardId: "board-123"},
					wantErr: true,
				},
				{
					name:    "missing board_id",
					req:     &boardpb.CreateColumnRequest{Name: "To Do"},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.CreateColumn(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("CreateColumn() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && resp.GetName() != tt.req.GetName() {
						t.Errorf("expected column name %q, got %q", tt.req.GetName(), resp.GetName())
					}
				})
			}
		})

		t.Run("UpdateColumn", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.UpdateColumnRequest
				wantErr bool
			}{
				{
					name: "update name",
					req: &boardpb.UpdateColumnRequest{
						Id:   "column-123",
						Name: wrapperspb.String("New Name"),
					},
				},
				{
					name:    "missing id",
					req:     &boardpb.UpdateColumnRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.UpdateColumn(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("UpdateColumn() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && name := tt.req.GetName(); name != nil && resp.GetName() != name.GetValue() {
						t.Errorf("expected column name %q, got %q", name.GetValue(), resp.GetName())
					}
				})
			}
		})

		t.Run("DeleteColumn", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.DeleteColumnRequest
				wantErr bool
			}{
				{
					name: "success",
					req:  &boardpb.DeleteColumnRequest{Id: "column-123"},
				},
				{
					name:    "missing id",
					req:     &boardpb.DeleteColumnRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					_, err := service.DeleteColumn(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("DeleteColumn() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
			}
		})
	})

	t.Run("TaskOperations", func(t *testing.T) {
		t.Run("CreateTask", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.CreateTaskRequest
				wantErr bool
			}{
				{
					name: "success",
					req: &boardpb.CreateTaskRequest{
						Name:     "Task 1",
						ColumnId: "column-123",
					},
				},
				{
					name:    "missing name",
					req:     &boardpb.CreateTaskRequest{ColumnId: "column-123"},
					wantErr: true,
				},
				{
					name:    "missing column_id",
					req:     &boardpb.CreateTaskRequest{Name: "Task 1"},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.CreateTask(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && resp.GetName() != tt.req.GetName() {
						t.Errorf("expected task name %q, got %q", tt.req.GetName(), resp.GetName())
					}
				})
			}
		})

		t.Run("MoveTask", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.MoveTaskRequest
				wantErr bool
			}{
				{
					name: "success",
					req: &boardpb.MoveTaskRequest{
						TaskId:      "task-123",
						NewColumnId: "column-456",
					},
				},
				{
					name:    "missing task_id",
					req:     &boardpb.MoveTaskRequest{NewColumnId: "column-456"},
					wantErr: true,
				},
				{
					name:    "missing new_column_id",
					req:     &boardpb.MoveTaskRequest{TaskId: "task-123"},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.MoveTask(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("MoveTask() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && resp.GetNewColumnId() != tt.req.GetNewColumnId() {
						t.Errorf("expected new column ID %q, got %q", tt.req.GetNewColumnId(), resp.GetNewColumnId())
					}
				})
			}
		})

		t.Run("UpdateTask", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.UpdateTaskRequest
				wantErr bool
			}{
				{
					name: "update name",
					req: &boardpb.UpdateTaskRequest{
						Id:   "task-123",
						Name: wrapperspb.String("New Name"),
					},
				},
				{
					name:    "missing id",
					req:     &boardpb.UpdateTaskRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.UpdateTask(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && name := tt.req.GetName(); name != nil && resp.GetName() != name.GetValue() {
						t.Errorf("expected task name %q, got %q", name.GetValue(), resp.GetName())
					}
				})
			}
		})

		t.Run("DeleteTask", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *boardpb.DeleteTaskRequest
				wantErr bool
			}{
				{
					name: "success",
					req:  &boardpb.DeleteTaskRequest{Id: "task-123"},
				},
				{
					name:    "missing id",
					req:     &boardpb.DeleteTaskRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					_, err := service.DeleteTask(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
			}
		})
	})
}
package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Board  BoardRepository
	Task   TaskRepository
	Column ColumnRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Board:  NewBoardRepository(db),
		Task:   NewTaskRepository(db),
		Column: NewColumnRepository(db),
	}
}

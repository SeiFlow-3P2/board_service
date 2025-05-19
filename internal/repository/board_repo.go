package repository

import (
	"context"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BoardRepository interface {
	CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error)
	GetBoardInfo(ctx context.Context, id uuid.UUID) (*models.Board, error)
	GetBoards(ctx context.Context, userID string) ([]*models.Board, error)
	UpdateBoard(ctx context.Context, id uuid.UUID, updates *BoardUpdates) (*models.Board, error)
	DeleteBoard(ctx context.Context, id uuid.UUID) error
	IncrementColumnsAmount(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	DecrementColumnsAmount(ctx context.Context, id uuid.UUID) error
}

type BoardUpdates struct {
	Title       *string    `bson:"title,omitempty"`
	Description *string    `bson:"description,omitempty"`
	Progress    *int       `bson:"progress,omitempty"`
	Favorite    *bool      `bson:"favorite,omitempty"`
	Updated_at  *time.Time `bson:"updated_at,omitempty"`
}

type boardRepository struct {
	db *mongo.Database
}

func NewBoardRepository(db *mongo.Database) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error) {
	collection := r.db.Collection("Boards")
	_, err := collection.InsertOne(ctx, board)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (r *boardRepository) GetBoardInfo(ctx context.Context, id uuid.UUID) (*models.Board, error) {
	collection := r.db.Collection("Boards")
	var board models.Board
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&board)
	if err != nil {
		return nil, err
	}

	columnsCursor, err := r.db.Collection("Columns").Find(ctx, bson.M{"desk_id": id})
	if err != nil {
		return nil, err
	}
	defer columnsCursor.Close(ctx)

	for columnsCursor.Next(ctx) {
		var column models.Column
		if err := columnsCursor.Decode(&column); err != nil {
			return nil, err
		}
		board.Columns = append(board.Columns, column)
	}
	if err := columnsCursor.Err(); err != nil {
		return nil, err
	}

	for i := range board.Columns {
		tasksCursor, err := r.db.Collection("Tasks").Find(ctx, bson.M{"column_id": board.Columns[i].ID})
		if err != nil {
			return nil, err
		}
		defer tasksCursor.Close(ctx)
		for tasksCursor.Next(ctx) {
			var task models.Task
			if err := tasksCursor.Decode(&task); err != nil {
				return nil, err
			}
			board.Columns[i].Tasks = append(board.Columns[i].Tasks, task)
		}
		if err := tasksCursor.Err(); err != nil {
			return nil, err
		}
	}

	return &board, nil
}

func (r *boardRepository) GetBoards(ctx context.Context, userID string) ([]*models.Board, error) {
	collection := r.db.Collection("Boards")
	var boards []*models.Board
	options := options.Find().SetProjection(bson.M{
		"_id": 1, "title": 1, "description": 1, "category": 1,
		"progress": 1, "favorite": 1, "metodology": 1,
		"columns_amount": 1, "created_at": 1, "updated_at": 1, "user_id": 1,
	})
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var board models.Board
		if err := cursor.Decode(&board); err != nil {
			return nil, err
		}
		boards = append(boards, &board)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return boards, nil
}

func (r *boardRepository) UpdateBoard(ctx context.Context, id uuid.UUID, updates *BoardUpdates) (*models.Board, error) {
	collection := r.db.Collection("Boards")

	updateFields := bson.M{}
	if updates.Title != nil {
		updateFields["title"] = *updates.Title
	}
	if updates.Description != nil {
		updateFields["description"] = *updates.Description
	}
	if updates.Progress != nil {
		updateFields["progress"] = *updates.Progress
	}
	if updates.Favorite != nil {
		updateFields["favorite"] = *updates.Favorite
	}
	if updates.Updated_at != nil {
		updateFields["updated_at"] = *updates.Updated_at
	}

	if len(updateFields) == 0 {
		return r.GetBoardInfo(ctx, id)
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, err
	}
	return r.GetBoardInfo(ctx, id)
}

func (r *boardRepository) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		columnsCollection := r.db.Collection("Columns")
		cursor, err := columnsCollection.Find(sc, bson.M{"desk_id": id})
		if err != nil {
			return err
		}
		defer cursor.Close(sc)

		tasksCollection := r.db.Collection("Tasks")
		for cursor.Next(sc) {
			var column models.Column
			if err := cursor.Decode(&column); err != nil {
				return err
			}
			_, err := tasksCollection.DeleteMany(sc, bson.M{"column_id": column.ID})
			if err != nil {
				return err
			}
		}
		_, err = columnsCollection.DeleteMany(sc, bson.M{"desk_id": id})
		if err != nil {
			return err
		}
		boardsCollection := r.db.Collection("Boards")
		_, err = boardsCollection.DeleteOne(sc, bson.M{"_id": id})
		if err != nil {
			return err
		}
		if err := session.CommitTransaction(sc); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return abortErr
		}
		return err
	}
	return nil
}

func (r *boardRepository) IncrementColumnsAmount(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	collection := r.db.Collection("Boards")
	update := bson.M{"$inc": bson.M{"columns_amount": 1}}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return uuid.Nil, err
	}
	var board models.Board
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&board)
	if err != nil {
		return uuid.Nil, err
	}
	return board.ID, nil
}

func (r *boardRepository) DecrementColumnsAmount(ctx context.Context, id uuid.UUID) error {
	collection := r.db.Collection("Boards")
	update := bson.M{"$inc": bson.M{"columns_amount": -1}}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

package repository

import (
	"context"

	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ColumnRepository interface {
	CreateColumn(ctx context.Context, column *models.Column) (*models.Column, error)
	GetColumnInfo(ctx context.Context, id uuid.UUID) (*models.Column, error)
	GetColumns(ctx context.Context, boardID uuid.UUID) ([]*models.Column, error)
	UpdateColumn(ctx context.Context, id uuid.UUID, updates *ColumnUpdates) (*models.Column, error)
	DeleteColumn(ctx context.Context, id uuid.UUID) error
}

type ColumnUpdates struct {
	Name *string `bson:"name,omitempty"`
}

type columnRepository struct {
	db *mongo.Database
}

func NewColumnRepository(db *mongo.Database) ColumnRepository {
	return &columnRepository{db: db}
}

func (r *columnRepository) CreateColumn(ctx context.Context, column *models.Column) (*models.Column, error) {
	collection := r.db.Collection("Columns")
	_, err := collection.InsertOne(ctx, column)
	if err != nil {
		return nil, err
	}
	return column, nil
}

func (r *columnRepository) GetColumnInfo(ctx context.Context, id uuid.UUID) (*models.Column, error) {
	collection := r.db.Collection("Columns")
	var column models.Column
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&column)
	if err != nil {
		return nil, err
	}
	return &column, nil
}

func (r *columnRepository) GetColumns(ctx context.Context, boardID uuid.UUID) ([]*models.Column, error) {
	collection := r.db.Collection("Columns")
	var columns []*models.Column
	options := options.Find().SetProjection(bson.M{"_id": 1, "name": 1, "order_number": 1})
	cursor, err := collection.Find(ctx, bson.M{"desk_id": boardID}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var column models.Column
		if err := cursor.Decode(&column); err != nil {
			return nil, err
		}
		columns = append(columns, &column)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return columns, nil
}

func (r *columnRepository) UpdateColumn(ctx context.Context, id uuid.UUID, updates *ColumnUpdates) (*models.Column, error) {
	collection := r.db.Collection("Columns")

	updateFields := bson.M{}
	if updates.Name != nil {
		updateFields["name"] = *updates.Name
	}

	if len(updateFields) == 0 {
		return r.GetColumnInfo(ctx, id)
	}

	_, err :=	collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, err
	}
	return r.GetColumnInfo(ctx, id)
}

func (r *columnRepository) DeleteColumn(ctx context.Context, id uuid.UUID) error {
	collection := r.db.Collection("Columns")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
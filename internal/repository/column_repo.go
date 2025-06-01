package repository

import (
	"context"

	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/SeiFlow-3P2/shared/telemetry"
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
	DecrementOrderNumbers(ctx context.Context, boardID uuid.UUID, orderNumber int) error
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
	ctx, span := telemetry.StartSpan(ctx, "ColumnRepository.CreateColumn")
	defer span.End()

	collection := r.db.Collection("Columns")
	_, err := collection.InsertOne(ctx, column)
	if err != nil {
		telemetry.RecordError(span, err)
		return nil, err
	}
	return column, nil
}

func (r *columnRepository) GetColumnInfo(ctx context.Context, id uuid.UUID) (*models.Column, error) {
	ctx, span := telemetry.StartSpan(ctx, "ColumnRepository.GetColumnInfo")
	defer span.End()

	collection := r.db.Collection("Columns")
	var column models.Column
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&column)
	if err != nil {
		telemetry.RecordError(span, err)
		return nil, err
	}
	return &column, nil
}

func (r *columnRepository) GetColumns(ctx context.Context, boardID uuid.UUID) ([]*models.Column, error) {
	ctx, span := telemetry.StartSpan(ctx, "ColumnRepository.GetColumns")
	defer span.End()

	collection := r.db.Collection("Columns")
	var columns []*models.Column
	options := options.Find().SetProjection(bson.M{"_id": 1, "name": 1, "order_number": 1})
	cursor, err := collection.Find(ctx, bson.M{"desk_id": boardID}, options)
	if err != nil {
		telemetry.RecordError(span, err)
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var column models.Column
		if err := cursor.Decode(&column); err != nil {
			telemetry.RecordError(span, err)
			return nil, err
		}
		columns = append(columns, &column)
	}
	if err := cursor.Err(); err != nil {
		telemetry.RecordError(span, err)
		return nil, err
	}
	return columns, nil
}

func (r *columnRepository) UpdateColumn(ctx context.Context, id uuid.UUID, updates *ColumnUpdates) (*models.Column, error) {
	ctx, span := telemetry.StartSpan(ctx, "ColumnRepository.UpdateColumn")
	defer span.End()

	collection := r.db.Collection("Columns")

	updateFields := bson.M{}
	if updates.Name != nil {
		updateFields["name"] = *updates.Name
	}

	if len(updateFields) == 0 {
		return r.GetColumnInfo(ctx, id)
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if err != nil {
		telemetry.RecordError(span, err)
		return nil, err
	}
	return r.GetColumnInfo(ctx, id)
}

func (r *columnRepository) DeleteColumn(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "ColumnRepository.DeleteColumn")
	defer span.End()

	collection := r.db.Collection("Columns")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		telemetry.RecordError(span, err)
		return err
	}
	return nil
}

func (r *columnRepository) DecrementOrderNumbers(ctx context.Context, boardID uuid.UUID, orderNumber int) error {
	ctx, span := telemetry.StartSpan(ctx, "ColumnRepository.DecrementOrderNumbers")
	defer span.End()

	collection := r.db.Collection("Columns")
	update := bson.M{"$inc": bson.M{"order_number": -1}}
	_, err := collection.UpdateMany(ctx, bson.M{
		"desk_id":      boardID,
		"order_number": bson.M{"$gt": orderNumber},
	}, update)
	if err != nil {
		telemetry.RecordError(span, err)
		return err
	}
	return err
}

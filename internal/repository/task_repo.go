package repository

import (
	"context"
	"time"

	"github.com/SeiFlow-3P2/board_service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *models.Task) (*models.Task, error)
	GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error)
	MoveTask(ctx context.Context, id uuid.UUID, newColumnID uuid.UUID) error
	UpdateTask(ctx context.Context, id uuid.UUID, updates *TaskUpdates) (*models.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

type TaskUpdates struct {
	Title       *string    `bson:"title,omitempty"`
	Description *string    `bson:"description,omitempty"`
	Deadline    *time.Time `bson:"deadline,omitempty"`
}

type taskRepository struct {
	db *mongo.Database
}

func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	collection := r.db.Collection("Tasks")
	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *taskRepository) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	collection := r.db.Collection("Tasks")
	var task models.Task
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) MoveTask(ctx context.Context, id uuid.UUID, newColumnID uuid.UUID) error {
	collection := r.db.Collection("Tasks")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"column_id": newColumnID}})
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, id uuid.UUID, updates *TaskUpdates) (*models.Task, error) {
	collection := r.db.Collection("Tasks")

	updateFields := bson.M{}
	if updates.Title != nil {
		updateFields["title"] = *updates.Title
	}
	if updates.Description != nil {
		updateFields["description"] = *updates.Description
	}
	if updates.Deadline != nil {
		updateFields["deadline"] = *updates.Deadline
	}

	if len(updateFields) == 0 {
		return r.GetTask(ctx, id)
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, err
	}
	return r.GetTask(ctx, id)
}

func (r *taskRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	collection := r.db.Collection("Tasks")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	ID             uuid.UUID `bson:"_id,omitempty"`
	Title          string    `bson:"title"`
	Description    string    `bson:"description"`
	Category       string    `bson:"category"`
	Progress       int       `bson:"progress"`
	Favorite       bool      `bson:"favorite"`
	Metodology     string    `bson:"metodology"`
	Columns_amount int       `bson:"columns_amount"`
	Created_at     time.Time `bson:"created_at"`
	Updated_at     time.Time `bson:"updated_at"`
	User_id        string    `bson:"user_id"`
	Columns        []Column  `bson:"columns,omitempty"`
}

type Column struct {
	ID           uuid.UUID `bson:"_id,omitempty"`
	Name         string    `bson:"name"`
	Order_number int       `bson:"order_number"`
	Desk_id      uuid.UUID `bson:"desk_id"`
	Tasks        []Task    `bson:"tasks,omitempty"`
}

type Task struct {
	ID          uuid.UUID `bson:"_id,omitempty"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Deadline    time.Time `bson:"deadline"`
	In_Calendar bool      `bson:"in_calendar"`
	Column_id   uuid.UUID `bson:"column_id"`
}

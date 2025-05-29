package models

import "time"

type BoardEvent struct {
	EventType   string    `json:"event_type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	UserID      string    `json:"user_id"`
}

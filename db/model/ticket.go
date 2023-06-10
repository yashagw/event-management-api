package model

import "time"

// Ticket represents a ticket in the database
type Ticket struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	EventID   int64     `json:"event_id"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

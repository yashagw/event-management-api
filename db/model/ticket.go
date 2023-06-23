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

type GetTicketParams struct {
	TicketID int64 `json:"ticket_id"`
	UserID   int64 `json:"user_id"`
}

type CreateTicketParams struct {
	UserID   int64 `json:"user_id"`
	EventID  int64 `json:"event_id"`
	Quantity int64 `json:"quantity"`
}

type DeleteTicketParams struct {
	UserID   int64 `json:"user_id"`
	TicketID int64 `json:"ticket_id"`
	EventID  int64 `json:"event_id"`
}

package model

import "time"

// Event represents an event in the database
type Event struct {
	ID           int64     `json:"id"`
	HostID       int64     `json:"host_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	TotalTickets int64     `json:"total_tickets"`
	LeftTickets  int64     `json:"left_tickets"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateEventParams struct {
	HostID       int64     `json:"host_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	TotalTickets int64     `json:"total_tickets"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

type GetEventParams struct {
	EventID int64 `json:"event_id"`
}

type ListEventsParams struct {
	HostID int64 `json:"host_id"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type ListEventsResponse struct {
	Records    []Event `json:"records"`
	NextOffset int     `json:"next_offset"`
}

package model

import "time"

// Event represents an event in the database
type Event struct {
	ID           int       `json:"id"`
	HostID       int       `json:"host_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	TotalTickets int       `json:"total_tickets"`
	TicketsLeft  int       `json:"tickets_left"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

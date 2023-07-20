package pgsql

import (
	"context"

	"github.com/yashagw/event-management-api/db/model"
)

func (provider *Provider) CreateEvent(context context.Context, request model.CreateEventParams) (*model.Event, error) {
	var event model.Event
	err := provider.conn.QueryRowContext(context, `
		INSERT INTO events (host_id, name, description, location, total_tickets, left_tickets, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, host_id, name, description, location, total_tickets, left_tickets, start_date, end_date, created_at
	`, request.HostID, request.Name, request.Description, request.Location, request.TotalTickets, request.TotalTickets, request.StartDate, request.EndDate).Scan(
		&event.ID,
		&event.HostID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.TotalTickets,
		&event.LeftTickets,
		&event.StartDate,
		&event.EndDate,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (provider *Provider) GetEvent(context context.Context, request model.GetEventParams) (*model.Event, error) {
	var event model.Event
	err := provider.conn.QueryRowContext(context, `
		SELECT id, host_id, name, description, location, total_tickets, left_tickets, start_date, end_date, created_at
		FROM events
		WHERE id = $1
	`, request.EventID).Scan(
		&event.ID,
		&event.HostID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.TotalTickets,
		&event.LeftTickets,
		&event.StartDate,
		&event.EndDate,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (provider *Provider) ListEvents(context context.Context, request model.ListEventsParams) (*model.ListEventsResponse, error) {
	if request.Limit <= 0 {
		request.Limit = 1
	}

	rows, err := provider.conn.QueryContext(context, `
		SELECT id, host_id, name, description, location, total_tickets, left_tickets, start_date, end_date, created_at
		FROM events
		WHERE host_id = $1
		LIMIT $2
		OFFSET $3
	`, request.HostID, request.Limit, request.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	nextOffset := request.Offset

	for rows.Next() {
		var event model.Event
		err := rows.Scan(
			&event.ID,
			&event.HostID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.TotalTickets,
			&event.LeftTickets,
			&event.StartDate,
			&event.EndDate,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
		nextOffset++
	}

	return &model.ListEventsResponse{
		Records:    events,
		NextOffset: nextOffset,
	}, nil
}

func (provider *Provider) DeleteEvent(context context.Context, id int64) error {
	_, err := provider.conn.ExecContext(context, `
		DELETE FROM events
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	return nil
}

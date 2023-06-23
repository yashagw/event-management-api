package pgsql

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/yashagw/event-management-api/db/model"
)

func (provider *Provider) GetTicket(context context.Context, request model.GetTicketParams) (*model.Ticket, error) {
	var ticket model.Ticket
	err := provider.conn.QueryRowContext(context, `
		SELECT id, user_id, event_id, quantity, created_at
		FROM tickets
		WHERE id = $1 AND user_id = $2
	`, request.TicketID, request.UserID).Scan(&ticket.ID, &ticket.UserID, &ticket.EventID, &ticket.Quantity, &ticket.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (p *Provider) CreateTicket(ctx context.Context, req model.CreateTicketParams) (*model.Ticket, error) {
	// Begin a transaction
	txProvider, err := p.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err != nil {
			txProvider.tx.Rollback()
		}
		txProvider.Close()
	}()

	// Check if there are enough tickets left for the event and lock the row
	var leftTickets int64
	err = txProvider.tx.QueryRowContext(ctx,
		"SELECT left_tickets FROM events WHERE id = $1 FOR UPDATE",
		req.EventID).Scan(&leftTickets)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if leftTickets < req.Quantity {
		return nil, errors.New("not enough tickets left for the event")
	}

	// Update the number of left tickets for the event
	_, err = txProvider.tx.ExecContext(ctx,
		"UPDATE events SET left_tickets = left_tickets - $1 WHERE id = $2",
		req.Quantity, req.EventID)
	if err != nil {
		return nil, err
	}

	// Create the ticket
	var ticketID int64
	var createdAt time.Time
	err = txProvider.tx.QueryRowContext(ctx,
		`INSERT INTO tickets (user_id, event_id, quantity) VALUES ($1, $2, $3) 
		RETURNING id, created_at`,
		req.UserID, req.EventID, req.Quantity).Scan(&ticketID, &createdAt)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = txProvider.tx.Commit()
	if err != nil {
		return nil, err
	}

	// Retrieve the created ticket
	createdTicket := &model.Ticket{
		ID:        ticketID,
		UserID:    req.UserID,
		EventID:   req.EventID,
		Quantity:  req.Quantity,
		CreatedAt: createdAt,
	}

	return createdTicket, nil
}

func (p *Provider) DeleteTicket(ctx context.Context, req model.DeleteTicketParams) error {
	// Begin a transaction
	txProvider, err := p.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if err != nil {
			txProvider.tx.Rollback()
		}
		txProvider.Close()
	}()

	// Check if the ticket exists and lock the row
	var ticketID int64
	err = txProvider.tx.QueryRowContext(ctx,
		"SELECT id FROM tickets WHERE id = $1 AND user_id = $2 FOR UPDATE",
		req.TicketID, req.UserID).Scan(&ticketID)
	if err != nil {
		return errors.WithStack(err)
	}

	// Delete the ticket
	var quantity int64
	err = txProvider.tx.QueryRowContext(ctx,
		"DELETE FROM tickets WHERE id = $1 RETURNING quantity",
		req.TicketID).Scan(&quantity)
	if err != nil {
		return errors.WithStack(err)
	}

	// Update the number of left tickets for the event
	_, err = txProvider.tx.ExecContext(ctx,
		"UPDATE events SET left_tickets = left_tickets + $1 WHERE id = $2",
		quantity, req.EventID)
	if err != nil {
		return errors.WithStack(err)
	}

	// Commit the transaction
	err = txProvider.tx.Commit()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

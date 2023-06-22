package pgsql

import (
	"context"

	"github.com/yashagw/event-management-api/db/model"
)

func (p *Provider) GetRequestToBecomeHost(context context.Context, userID int64) (*model.UserHostRequest, error) {
	var request model.UserHostRequest
	err := p.conn.QueryRowContext(context, `
		SELECT id, user_id, moderator_id, status, created_at, updated_at
		 FROM user_host_requests 
		 WHERE user_id = $1
		`, userID).Scan(
		&request.ID,
		&request.UserID,
		&request.ModeratorID,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (p *Provider) CreateRequestToBecomeHost(context context.Context, userID int64) error {
	_, err := p.conn.ExecContext(context, `
		INSERT INTO user_host_requests (user_id, status) VALUES ($1, $2)
		`, userID, model.UserHostRequestStatus_Pending)
	return err
}

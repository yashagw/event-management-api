package pgsql

import (
	"context"

	"github.com/yashagw/event-management-api/db/model"
)

func (p *Provider) ListPendingRequests(context context.Context, req model.ListPendingRequestsParams) (*model.ListPendingRequestsResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	rows, err := p.conn.QueryContext(context, `
		SELECT id, user_id, moderator_id, status, created_at, updated_at
		 FROM user_host_requests 
		 WHERE status = $1
		 LIMIT $2 OFFSET $3
		`, model.UserHostRequestStatus_Pending, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	var requests []*model.UserHostRequest
	nextoffset := req.Offset

	for rows.Next() {
		var request model.UserHostRequest
		err := rows.Scan(
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

		requests = append(requests, &request)
		nextoffset++
	}

	response := model.ListPendingRequestsResponse{
		Records:    requests,
		NextOffset: nextoffset,
	}

	return &response, nil
}

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

func (p *Provider) CreateRequestToBecomeHost(context context.Context, userID int64) (*model.UserHostRequest, error) {
	var request model.UserHostRequest
	err := p.conn.QueryRowContext(context, `
		INSERT INTO user_host_requests (user_id, status) VALUES ($1, $2)
		RETURNING id, user_id, moderator_id, status, created_at, updated_at
		`, userID, model.UserHostRequestStatus_Pending).Scan(
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

func (p *Provider) DeleteRequestToBecomeHost(context context.Context, id int64) error {
	_, err := p.conn.ExecContext(context, `
		DELETE FROM user_host_requests WHERE id = $1
		`, id)
	return err
}

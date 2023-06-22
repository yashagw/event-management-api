package pgsql

import (
	"context"
	"time"

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

func (p *Provider) ApproveDisapproveRequestToBecomeHost(context context.Context, request model.ApproveDisapproveRequestToBecomeHostParams) error {
	txProvider, err := p.BeginTx(context, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			txProvider.tx.Rollback()
		}
	}()

	if request.Approved {
		var userID int64

		err := txProvider.tx.QueryRowContext(context, `
			UPDATE user_host_requests SET status = $1, moderator_id = $2, updated_at = $3 WHERE id = $4
			RETURNING user_id
		`, model.UserHostRequestStatus_Approved, request.ModeratorID, time.Now(), request.RequestID).Scan(&userID)
		if err != nil {
			return err
		}

		_, err = txProvider.tx.ExecContext(context, `
			UPDATE users SET role = $1 WHERE id = $2
			`, model.UserRole_Host, userID)
		if err != nil {
			return err
		}

	} else {
		_, err = txProvider.tx.ExecContext(context, `
			UPDATE user_host_requests SET status = $1 WHERE id = $2
			`, model.UserHostRequestStatus_Rejected, request.RequestID)
	}

	if err := txProvider.tx.Commit(); err != nil {
		return err
	}

	txProvider.Close()

	return nil
}

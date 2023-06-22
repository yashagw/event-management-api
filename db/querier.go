package db

import (
	"context"

	"github.com/yashagw/event-management-api/db/model"
)

type UserQuerier interface {
	// CreateUser creates a new user in the database
	CreateUser(context context.Context, arg model.CreateUserParams) (*model.User, error)
	GetUserByEmail(context context.Context, email string) (*model.User, error)
	DeleteUser(context context.Context, id int64) error

	CreateRequestToBecomeHost(context context.Context, userID int64) (*model.UserHostRequest, error)
	GetRequestToBecomeHost(context context.Context, userID int64) (*model.UserHostRequest, error)
	DeleteRequestToBecomeHost(context context.Context, id int64) error
	ListPendingRequests(context context.Context, request model.ListPendingRequestsParams) (*model.ListPendingRequestsResponse, error)
}

type EventQuerier interface {
}

type DBQuerier interface {
	UserQuerier
	EventQuerier
}

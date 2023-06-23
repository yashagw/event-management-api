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
	ApproveDisapproveRequestToBecomeHost(context context.Context, request model.ApproveDisapproveRequestToBecomeHostParams) error
}

type EventQuerier interface {
	CreateEvent(context context.Context, request model.CreateEventParams) (*model.Event, error)
	GetEvent(context context.Context, request model.GetEventParams) (*model.Event, error)
	ListEvents(context context.Context, request model.ListEventsParams) (*model.ListEventsResponse, error)
	DeleteEvent(context context.Context, id int64) error
}

type TicketQuerier interface {
	CreateTicket(context context.Context, request model.CreateTicketParams) (*model.Ticket, error)
	GetTicket(context context.Context, request model.GetTicketParams) (*model.Ticket, error)
	DeleteTicket(context context.Context, request model.DeleteTicketParams) error
}

type DBQuerier interface {
	UserQuerier
	EventQuerier
	TicketQuerier
}

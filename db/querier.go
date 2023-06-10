package db

import (
	"context"

	"github.com/yashagw/event-management-api/db/model"
)

type UserQuerier interface {
	// CreateUser creates a new user in the database
	CreateUser(context context.Context, arg model.CreateUserParams) (*model.User, error)
	GetUserByEmail(context context.Context, email string) (*model.User, error)
}

type EventQuerier interface {
}

type DBQuerier interface {
	UserQuerier
	EventQuerier
}

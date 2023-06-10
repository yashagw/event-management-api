package pgsql

import (
	"context"

	"github.com/yashagw/event-management-api/db/model"
)

// CreateUser creates a new user in the database
func (p *Provider) CreateUser(context context.Context, arg model.CreateUserParams) (*model.User, error) {
	user := &model.User{}
	err := p.conn.QueryRowContext(context, `
		INSERT INTO users (name, email, hashed_password, role)
		VAlUES ($1, $2, $3, $4)
		RETURNING id, name, email, hashed_password, role, created_at, password_updated_at
	`, arg.Name, arg.Email, arg.HashedPassword, model.UserRole_User).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
		&user.CreatedAt,
		&user.PasswordUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail gets a user by email
func (p *Provider) GetUserByEmail(context context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := p.conn.QueryRowContext(context, `
		SELECT id, name, email, hashed_password, role, created_at, password_updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
		&user.CreatedAt,
		&user.PasswordUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

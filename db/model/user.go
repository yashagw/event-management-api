package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/yashagw/event-management-api/pb"
)

type UserRole int

const (
	UserRole_User UserRole = iota
	UserRole_Host
	UserRole_Moderator
	UserRole_Admin
)

// Implement the Scan method for UserRole
// It is used by the sql package to convert a value from the database into a UserRole
func (es *UserRole) Scan(value interface{}) error {
	if value == nil {
		*es = 0
		return nil
	}

	intValue, ok := value.(int64)
	if !ok {
		return fmt.Errorf("cannot scan value into UserRole")
	}

	*es = UserRole(intValue)
	return nil
}

// Implement the Value method for UserRole
// It is used by the sql package to convert a UserRole into a value that can be stored in the database
func (es UserRole) Value() (driver.Value, error) {
	return int64(es), nil
}

// Convert model.UserRole to pb.UserRole
func (role UserRole) ToProto() pb.UserRole {
	switch role {
	case UserRole_User:
		return pb.UserRole_UserRole_User
	case UserRole_Host:
		return pb.UserRole_UserRole_Host
	case UserRole_Moderator:
		return pb.UserRole_UserRole_Moderator
	case UserRole_Admin:
		return pb.UserRole_UserRole_Admin
	default:
		return pb.UserRole_UserRole_User
	}
}

// User represents a user in the database
type User struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	HashedPassword    string    `json:"hashed_password"`
	Role              UserRole  `json:"role"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordUpdatedAt time.Time `json:"password_updated_at"`
}

// CreateEventParams represents parameters to create an user
type CreateUserParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

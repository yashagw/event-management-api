package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type UserHostRequestStatus int

const (
	UserHostRequestStatus_Pending UserHostRequestStatus = iota
	UserHostRequestStatus_Rejected
	UserHostRequestStatus_Approved
)

// Implement the Scan method for UserRole
// It is used by the sql package to convert a value from the database into a UserRole
func (es *UserHostRequestStatus) Scan(value interface{}) error {
	if value == nil {
		*es = 0
		return nil
	}

	intValue, ok := value.(int64)
	if !ok {
		return fmt.Errorf("cannot scan value into UserHostRequestStatus")
	}

	*es = UserHostRequestStatus(intValue)
	return nil
}

// Implement the Value method for UserRole
// It is used by the sql package to convert a UserRole into a value that can be stored in the database
func (es UserHostRequestStatus) Value() (driver.Value, error) {
	return int64(es), nil
}

// UserHostRequest represents a request to become a host in the database
type UserHostRequest struct {
	ID          int64                 `json:"id"`
	UserID      int64                 `json:"user_id"`
	ModeratorID sql.NullInt64         `json:"moderator_id"`
	Status      UserHostRequestStatus `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type ListPendingRequestsParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListPendingRequestsResponse struct {
	Records    []*UserHostRequest `json:"records"`
	NextOffset int                `json:"next_offset"`
}

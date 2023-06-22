package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserHostRequestRole(t *testing.T) {
	var role UserHostRequestStatus

	// Test Scan method with valid value
	err := role.Scan(int64(UserHostRequestStatus_Pending))
	require.NoError(t, err)
	require.Equal(t, UserHostRequestStatus_Pending, role)

	// Test Scan method with nil value
	err = role.Scan(nil)
	require.NoError(t, err)
	require.Equal(t, UserHostRequestStatus_Pending, role) // Default value should be UserHostReqestStatus_Pending

	// Test Scan method with unsupported value type
	err = role.Scan("invalid")
	require.Error(t, err)
	require.Equal(t, UserHostRequestStatus_Pending, role) // Value should remain unchanged
}

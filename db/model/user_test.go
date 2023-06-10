package model

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRole(t *testing.T) {
	var ur UserRole

	// Test Scan method with valid value
	err := ur.Scan(int64(UserRole_Host))
	assert.NoError(t, err)
	assert.Equal(t, UserRole_Host, ur)

	// Test Scan method with nil value
	err = ur.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, UserRole_User, ur) // Default value should be UserRole_User

	// Test Scan method with unsupported value type
	err = ur.Scan("invalid")
	assert.Error(t, err)
	assert.Equal(t, UserRole_User, ur) // Value should remain unchanged

	// Test Value method
	val, err := ur.Value()
	assert.NoError(t, err)
	assert.Equal(t, driver.Value(int64(UserRole_User)), val)
}

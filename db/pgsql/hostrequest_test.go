package pgsql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yashagw/event-management-api/db/model"
)

func TestCreateRequestToBecomeHost(t *testing.T) {
	user := CreateRandomUser(t)
	err := provider.CreateRequestToBecomeHost(context.Background(), user.ID)
	require.NoError(t, err)

	// Unique Constraint on user_id so this should fail
	err = provider.CreateRequestToBecomeHost(context.Background(), user.ID)
	require.Error(t, err)
}

func TestGetRequestToBecomeHost(t *testing.T) {
	user := CreateRandomUser(t)
	err := provider.CreateRequestToBecomeHost(context.Background(), user.ID)
	require.NoError(t, err)

	request, err := provider.GetRequestToBecomeHost(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user.ID, request.UserID)
	require.Equal(t, model.UserHostRequestStatus_Pending, request.Status)
}

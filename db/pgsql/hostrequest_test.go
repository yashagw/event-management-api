package pgsql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yashagw/event-management-api/db/model"
)

func TestCreateRequestToBecomeHost(t *testing.T) {
	user := CreateRandomUser(t)
	request, err := provider.CreateRequestToBecomeHost(context.Background(), user.ID)
	require.NoError(t, err)
	defer func() {
		err = provider.DeleteRequestToBecomeHost(context.Background(), request.ID)
		require.NoError(t, err)

		err = provider.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)
	}()

	// Unique Constraint on user_id so this should fail
	_, err = provider.CreateRequestToBecomeHost(context.Background(), user.ID)
	require.Error(t, err)
}

func TestGetRequestToBecomeHost(t *testing.T) {
	user := CreateRandomUser(t)
	request, err := provider.CreateRequestToBecomeHost(context.Background(), user.ID)
	require.NoError(t, err)
	defer func() {
		err = provider.DeleteRequestToBecomeHost(context.Background(), request.ID)
		require.NoError(t, err)

		err = provider.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)
	}()

	request, err = provider.GetRequestToBecomeHost(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user.ID, request.UserID)
	require.Equal(t, model.UserHostRequestStatus_Pending, request.Status)
}

func TestListPendingRequests(t *testing.T) {
	for i := 0; i < 3; i++ {
		u := CreateRandomUser(t)
		request, err := provider.CreateRequestToBecomeHost(context.Background(), u.ID)
		require.NoError(t, err)
		defer func() {
			err = provider.DeleteRequestToBecomeHost(context.Background(), request.ID)
			require.NoError(t, err)

			err = provider.DeleteUser(context.Background(), u.ID)
			require.NoError(t, err)
		}()
	}

	req := model.ListPendingRequestsParams{
		Limit:  1,
		Offset: 0,
	}
	res1, err := provider.ListPendingRequests(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, 1, len(res1.Records))
	require.Equal(t, 1, res1.NextOffset)

	req = model.ListPendingRequestsParams{
		Limit:  1,
		Offset: 1,
	}
	res2, err := provider.ListPendingRequests(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, 1, len(res2.Records))
	require.Equal(t, 2, res2.NextOffset)

	req = model.ListPendingRequestsParams{
		Limit:  2,
		Offset: 0,
	}
	res3, err := provider.ListPendingRequests(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, 2, len(res3.Records))
	require.Equal(t, res1.Records[0].ID, res3.Records[0].ID)
	require.Equal(t, res2.Records[0].ID, res3.Records[1].ID)
	require.Equal(t, 2, res3.NextOffset)
}

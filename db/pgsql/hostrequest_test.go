package pgsql

import (
	"context"
	"database/sql"
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

func TestApproveDisapproveRequestToBecomeHost(t *testing.T) {
	user := CreateRandomUser(t)
	request, err := provider.CreateRequestToBecomeHost(context.Background(), user.ID)
	require.NoError(t, err)
	defer func() {
		err = provider.DeleteRequestToBecomeHost(context.Background(), request.ID)
		require.NoError(t, err)

		err = provider.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)
	}()

	moderatorId := CreateRandomUser(t).ID
	err = provider.ApproveDisapproveRequestToBecomeHost(context.Background(), model.ApproveDisapproveRequestToBecomeHostParams{
		Approved:    true,
		RequestID:   request.ID,
		ModeratorID: moderatorId,
	})
	require.NoError(t, err)

	r, err := provider.GetRequestToBecomeHost(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, model.UserHostRequestStatus_Approved, r.Status)
	require.Equal(t, sql.NullInt64{Int64: moderatorId, Valid: true}, r.ModeratorID)

	u1, err := provider.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.Equal(t, model.UserRole_Host, u1.Role)

	user2 := CreateRandomUser(t)
	request2, err := provider.CreateRequestToBecomeHost(context.Background(), user2.ID)
	require.NoError(t, err)
	defer func() {
		err = provider.DeleteRequestToBecomeHost(context.Background(), request2.ID)
		require.NoError(t, err)

		err = provider.DeleteUser(context.Background(), user2.ID)
		require.NoError(t, err)
	}()

	err = provider.ApproveDisapproveRequestToBecomeHost(context.Background(), model.ApproveDisapproveRequestToBecomeHostParams{
		Approved:    false,
		RequestID:   request2.ID,
		ModeratorID: moderatorId,
	})
	require.NoError(t, err)

	r, err = provider.GetRequestToBecomeHost(context.Background(), user2.ID)
	require.NoError(t, err)
	require.Equal(t, model.UserHostRequestStatus_Rejected, r.Status)
	require.Equal(t, sql.NullInt64{Int64: 0, Valid: false}, r.ModeratorID)

	u2, err := provider.GetUserByEmail(context.Background(), user2.Email)
	require.NoError(t, err)
	require.Equal(t, model.UserRole_User, u2.Role)
}

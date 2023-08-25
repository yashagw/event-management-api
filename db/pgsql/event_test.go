package pgsql

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/util"
)

func CreateRandomEvent(t *testing.T, host *model.User) *model.Event {
	arg := model.CreateEventParams{
		HostID:       host.ID,
		Name:         util.RandomName(),
		Description:  util.RandomString(10),
		Location:     util.RandomString(10),
		TotalTickets: util.RandomInt(1, 100),
		StartDate:    time.Now().UTC(),
		EndDate:      time.Now().Add(time.Hour * 24).UTC(),
	}

	event, err := provider.CreateEvent(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, event.ID)
	require.Equal(t, arg.HostID, event.HostID)
	require.Equal(t, arg.Name, event.Name)
	require.Equal(t, arg.Description, event.Description)
	require.Equal(t, arg.Location, event.Location)
	require.Equal(t, arg.TotalTickets, event.TotalTickets)
	require.Equal(t, arg.TotalTickets, event.LeftTickets)
	require.Equal(t, arg.StartDate.Format(time.RFC3339), event.StartDate.Format(time.RFC3339)) // Compare formatted time strings
	require.Equal(t, arg.EndDate.Format(time.RFC3339), event.EndDate.Format(time.RFC3339))     // Compare formatted time strings
	require.NotEmpty(t, event.CreatedAt)

	return event
}

func TestCreateEvent(t *testing.T) {
	host := CreateRandomUser(t)
	event := CreateRandomEvent(t, host)
	defer func() {
		err := provider.DeleteEvent(context.Background(), event.ID)
		require.NoError(t, err)

		err = provider.DeleteUser(context.Background(), host.ID)
		require.NoError(t, err)
	}()

	fetchedEvent, err := provider.GetEvent(context.Background(), model.GetEventParams{
		EventID: event.ID,
	})
	require.NoError(t, err)
	require.Equal(t, event, fetchedEvent)
}

func TestDeleteEvent(t *testing.T) {
	host := CreateRandomUser(t)
	event := CreateRandomEvent(t, host)
	defer func() {
		err := provider.DeleteEvent(context.Background(), event.ID)
		require.NoError(t, err)

		err = provider.DeleteUser(context.Background(), host.ID)
		require.NoError(t, err)
	}()

	// Delete the event
	err := provider.DeleteEvent(context.Background(), event.ID)
	require.NoError(t, err)

	// Fetch the event from the database
	_, err = provider.GetEvent(context.Background(), model.GetEventParams{
		EventID: event.ID,
	})
	require.Error(t, err)
}

func TestListEvents(t *testing.T) {
	user := CreateRandomUser(t)
	defer func() {
		err := provider.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)
	}()

	user2 := CreateRandomUser(t)
	defer func() {
		err := provider.DeleteUser(context.Background(), user2.ID)
		require.NoError(t, err)
	}()

	// Create 10 random events for user
	for i := 0; i < 10; i++ {
		event := CreateRandomEvent(t, user)
		defer func() {
			err := provider.DeleteEvent(context.Background(), event.ID)
			require.NoError(t, err)
		}()
	}

	// Create 10 random events for user2
	for i := 0; i < 10; i++ {
		event := CreateRandomEvent(t, user2)
		defer func() {
			err := provider.DeleteEvent(context.Background(), event.ID)
			require.NoError(t, err)
		}()
	}

	// List the events filter by host_id
	response, err := provider.ListEvents(context.Background(), model.ListEventsParams{
		HostID: user.ID,
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Equal(t, 10, len(response.Records))

	// List all events
	response, err = provider.ListEvents(context.Background(), model.ListEventsParams{})
	require.NoError(t, err)
	require.Equal(t, 20, len(response.Records))
}

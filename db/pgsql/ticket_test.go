package pgsql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/util"
)

func CreateRandomTicket(t *testing.T, user *model.User, event *model.Event) *model.Ticket {
	arg := model.CreateTicketParams{
		UserID:   user.ID,
		EventID:  event.ID,
		Quantity: util.RandomInt(1, event.LeftTickets),
	}

	ticket, err := provider.CreateTicket(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ticket.ID)
	require.Equal(t, arg.UserID, ticket.UserID)
	require.Equal(t, arg.EventID, ticket.EventID)
	require.Equal(t, arg.Quantity, ticket.Quantity)
	require.NotEmpty(t, ticket.CreatedAt)

	return ticket
}

func TestCreateTicket(t *testing.T) {
	user := CreateRandomUser(t)
	event := CreateRandomEvent(t, user)
	defer func() {
		err := provider.DeleteEvent(context.Background(), event.ID)
		require.NoError(t, err)

		err = provider.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)
	}()

	ticket := CreateRandomTicket(t, user, event)
	defer func() {
		arg := model.DeleteTicketParams{
			UserID:   user.ID,
			TicketID: ticket.ID,
			EventID:  event.ID,
		}
		err := provider.DeleteTicket(context.Background(), arg)
		require.NoError(t, err)
	}()

	fetchedTicket, err := provider.GetTicket(context.Background(), model.GetTicketParams{
		UserID:   user.ID,
		TicketID: ticket.ID,
	})
	require.NoError(t, err)
	require.Equal(t, ticket, fetchedTicket)
}

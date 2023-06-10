package pgsql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/util"
)

func CreateRandomUser(t *testing.T) *model.User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := model.CreateUserParams{
		Name:           util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := provider.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, model.UserRole_User, user.Role)
	require.NotEmpty(t, user.CreatedAt)
	require.NotEmpty(t, user.PasswordUpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUserByEmail(t *testing.T) {
	userResponse := CreateRandomUser(t)

	userResponse2, err := provider.GetUserByEmail(context.Background(), userResponse.Email)
	require.NoError(t, err)
	require.Equal(t, userResponse, userResponse2)

	_, err = provider.GetUserByEmail(context.Background(), util.RandomEmail())
	require.Error(t, err)
}

package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yashagw/event-management-api/db/mock"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/util"
	mockwk "github.com/yashagw/event-management-api/worker/mock"
)

type eqCreateUserParamsMatcher struct {
	arg      model.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(model.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg model.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (user model.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = model.User{
		Name:           util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		Role:           model.UserRole_User,
		CreatedAt:      time.Now(),
	}

	return user, password
}

func TestLoginUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testcases := []struct {
		name          string
		body          gin.H
		buildStubs    func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Not Found Email",
			body: gin.H{
				"email":    "invalid@gmail.com",
				"password": password,
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), "invalid@gmail.com").Times(1).Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Password",
			body: gin.H{
				"email":    user.Email,
				"password": "invalid1234",
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"email":    "invalid",
				"password": password,
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), "invalid").Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			providerCtrl := gomock.NewController(t)
			defer providerCtrl.Finish()
			provider := mockdb.NewMockProvider(providerCtrl)

			redisCtrl := gomock.NewController(t)
			defer redisCtrl.Finish()
			distributor := mockwk.NewMockTaskDistributor(redisCtrl)

			tc.buildStubs(provider, distributor)

			server := newTestServer(t, provider, distributor)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testcases := []struct {
		name          string
		body          gin.H
		buildStubs    func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				arg := model.CreateUserParams{
					Name:  user.Name,
					Email: user.Email,
				}

				provider.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).Return(&user, nil)
				worker.EXPECT().DistributeTaskSendEmailVerify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateEmail",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).
					Return(nil, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"name":     user.Name,
				"email":    "invalid",
				"password": password,
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"password": "short",
			},
			buildStubs: func(provider *mockdb.MockProvider, worker *mockwk.MockTaskDistributor) {
				provider.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			providerCtrl := gomock.NewController(t)
			defer providerCtrl.Finish()
			provider := mockdb.NewMockProvider(providerCtrl)

			redisCtrl := gomock.NewController(t)
			defer redisCtrl.Finish()
			distributor := mockwk.NewMockTaskDistributor(redisCtrl)

			tc.buildStubs(provider, distributor)

			server := newTestServer(t, provider, distributor)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user model.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser model.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}

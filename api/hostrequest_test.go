package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yashagw/event-management-api/db/mock"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/token"
)

func randomUserHostRequest(t *testing.T, user *model.User, status model.UserHostRequestStatus) model.UserHostRequest {
	res := model.UserHostRequest{
		UserID: user.ID,
		Status: status,
	}

	return res
}

func TestListPendingUserHostRequests(t *testing.T) {
	user, _ := randomUser(t)

	moderator, _ := randomUser(t)
	moderator.Role = model.UserRole_Moderator

	type Query struct {
		Limit  int
		Offset int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(provider *mockdb.MockProvider)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			query: Query{
				Limit:  5,
				Offset: 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, moderator.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				var requests []*model.UserHostRequest
				for i := 0; i < 5; i++ {
					user, _ := randomUser(t)
					request := randomUserHostRequest(t, &user, model.UserHostRequestStatus_Pending)
					requests = append(requests, &request)
				}

				arg := model.ListPendingRequestsParams{
					Limit:  5,
					Offset: 0,
				}
				res := model.ListPendingRequestsResponse{
					Records:    requests,
					NextOffset: 5,
				}

				provider.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(moderator.Email)).
					Times(1).Return(&moderator, nil)
				provider.EXPECT().ListPendingRequests(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(&res, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Not Authorized",
			query: Query{
				Limit:  5,
				Offset: 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			query: Query{
				Limit:  5,
				Offset: 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, moderator.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				arg := model.ListPendingRequestsParams{
					Limit:  5,
					Offset: 0,
				}

				provider.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(moderator.Email)).
					Times(1).Return(&moderator, nil)
				provider.EXPECT().ListPendingRequests(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			provider := mockdb.NewMockProvider(ctrl)
			tc.buildStubs(provider)

			server := newTestServer(t, provider)
			recorder := httptest.NewRecorder()

			url := "/moderator/requests"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("limit", fmt.Sprintf("%d", tc.query.Limit))
			q.Add("offset", fmt.Sprintf("%d", tc.query.Offset))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestBecomeHostRequest(t *testing.T) {
	user, _ := randomUser(t)

	host, _ := randomUser(t)
	host.Role = model.UserRole_Host

	testcases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(provider *mockdb.MockProvider)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Okay",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), user.Email).
					Times(1).Return(&user, nil)
				provider.EXPECT().CreateRequestToBecomeHost(gomock.Any(), user.ID).
					Times(1).Return(nil, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, "{\"message\":\"request to become host created\"}", recorder.Body.String())
			},
		},
		{
			name: "Not Authorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "wrongemail", time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), "wrongemail").
					Times(1).Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Already Host",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, host.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), host.Email).
					Times(1).Return(&host, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), user.Email).
					Times(1).Return(&user, nil)
				provider.EXPECT().CreateRequestToBecomeHost(gomock.Any(), user.ID).
					Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unique Constraint Violation",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), user.Email).
					Times(1).Return(&user, nil)
				provider.EXPECT().CreateRequestToBecomeHost(gomock.Any(), user.ID).
					Times(1).Return(nil, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			provider := mockdb.NewMockProvider(ctrl)
			tc.buildStubs(provider)

			server := newTestServer(t, provider)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodPost, "/users/become-host", nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

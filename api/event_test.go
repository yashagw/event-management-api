package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yashagw/event-management-api/db/mock"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/token"
	mockwk "github.com/yashagw/event-management-api/worker/mock"
)

func TestCreateEvent(t *testing.T) {
	host, _ := randomUser(t)
	host.Role = model.UserRole_Host

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(provider *mockdb.MockProvider)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":          "Test Event",
				"description":   "Test Description",
				"location":      "Test Location",
				"total_tickets": 100,
				"start_date":    "2021-01-01T00:00:00Z",
				"end_date":      "2021-01-02T00:00:00Z",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, host.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				arg := model.CreateEventParams{
					HostID:       host.ID,
					Name:         "Test Event",
					Description:  "Test Description",
					Location:     "Test Location",
					TotalTickets: 100,
					StartDate:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:      time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				}

				provider.EXPECT().GetUserByEmail(gomock.Any(), host.Email).Times(1).Return(&host, nil)
				provider.EXPECT().CreateEvent(gomock.Any(), gomock.Eq(arg)).Times(1).Return(&model.Event{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			providerCtrl := gomock.NewController(t)
			defer providerCtrl.Finish()
			provider := mockdb.NewMockProvider(providerCtrl)
			tc.buildStubs(provider)

			redisCtrl := gomock.NewController(t)
			defer redisCtrl.Finish()
			distributor := mockwk.NewMockTaskDistributor(redisCtrl)

			server := newTestServer(t, provider, distributor)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/hosts/events", bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListEvents(t *testing.T) {
	host, _ := randomUser(t)
	host.Role = model.UserRole_Host

	user, _ := randomUser(t)

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
			name: "OK",
			query: Query{
				Limit:  10,
				Offset: 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, host.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				arg := model.ListEventsParams{
					Limit:  10,
					Offset: 0,
				}

				provider.EXPECT().GetUserByEmail(gomock.Any(), host.Email).Times(1).Return(&host, nil)
				provider.EXPECT().ListEvents(gomock.Any(), gomock.Eq(arg)).Times(1).Return(&model.ListEventsResponse{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			query: Query{
				Limit:  10,
				Offset: 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(&user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			providerCtrl := gomock.NewController(t)
			defer providerCtrl.Finish()
			provider := mockdb.NewMockProvider(providerCtrl)
			tc.buildStubs(provider)

			redisCtrl := gomock.NewController(t)
			defer redisCtrl.Finish()
			distributor := mockwk.NewMockTaskDistributor(redisCtrl)

			server := newTestServer(t, provider, distributor)
			recorder := httptest.NewRecorder()

			url := "/hosts/events"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

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

func TestGetEvent(t *testing.T) {
	testCases := []struct {
		name          string
		eventID       int
		buildStubs    func(provider *mockdb.MockProvider)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			eventID: 1,
			buildStubs: func(provider *mockdb.MockProvider) {
				arg := model.GetEventParams{
					EventID: 1,
				}
				provider.EXPECT().GetEvent(gomock.Any(), gomock.Eq(arg)).Times(1).Return(&model.Event{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			providerCtrl := gomock.NewController(t)
			defer providerCtrl.Finish()
			provider := mockdb.NewMockProvider(providerCtrl)
			tc.buildStubs(provider)

			redisCtrl := gomock.NewController(t)
			defer redisCtrl.Finish()
			distributor := mockwk.NewMockTaskDistributor(redisCtrl)

			server := newTestServer(t, provider, distributor)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/events/%d", tc.eventID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

package api

import (
	"bytes"
	"encoding/json"
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

func TestCreateTicket(t *testing.T) {
	user, _ := randomUser(t)

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
				"event_id": 1,
				"quantity": 1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				arg := model.CreateTicketParams{
					EventID:  1,
					UserID:   user.ID,
					Quantity: 1,
				}

				provider.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(&user, nil)
				provider.EXPECT().CreateTicket(gomock.Any(), gomock.Eq(arg)).Times(1).Return(&model.Ticket{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "No Authorization",
			body: gin.H{
				"event_id": 1,
				"quantity": 1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, host.Email, time.Minute)
			},
			buildStubs: func(provider *mockdb.MockProvider) {
				provider.EXPECT().GetUserByEmail(gomock.Any(), host.Email).Times(1).Return(&host, nil)
				provider.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users/ticket", bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

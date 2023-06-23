package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/token"
)

type CreateTicketParams struct {
	EventID  int64 `json:"event_id"`
	Quantity int64 `json:"quantity"`
}

func (server *Server) CreateTicket(context *gin.Context) {
	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	userEmail := payload.Username

	user, err := server.provider.GetUserByEmail(context, userEmail)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.Role != model.UserRole_User {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	var params model.CreateTicketParams
	if err := context.ShouldBindJSON(&params); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ticket, err := server.provider.CreateTicket(context, model.CreateTicketParams{
		EventID:  params.EventID,
		UserID:   user.ID,
		Quantity: params.Quantity,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, ticket)
}

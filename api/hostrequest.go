package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/token"
)

func (server *Server) BecomeHost(context *gin.Context) {
	// TODO: Delete old request for user if not approved after 30 days
	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	userEmail := payload.Username

	user, err := server.provider.GetUserByEmail(context, userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Authorization to make sure payload have user role
	if user.Role != model.UserRole_User {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	err = server.provider.CreateRequestToBecomeHost(context, user.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				context.JSON(http.StatusConflict, gin.H{"message": "request to become host already exists"})
				return
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "request to become host created"})
}

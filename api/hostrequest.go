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

	// Authorization to make sure payload have user role
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

	if user.Role != model.UserRole_User {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	// Create request to become host
	_, err = server.provider.CreateRequestToBecomeHost(context, user.ID)
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

type ListPendingUserHostRequestsParams struct {
	Limit  int `form:"limit" binding:"required,min=1,max=1000"`
	Offset int `form:"offset"`
}

func (server *Server) ListPendingUserHostRequests(context *gin.Context) {
	// Authorization to make sure payload have user role
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

	// Only admin and moderator can list pending requests
	if user.Role == model.UserRole_User || user.Role == model.UserRole_Host {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	var req ListPendingUserHostRequestsParams
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.provider.ListPendingRequests(context, model.ListPendingRequestsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, response)
}

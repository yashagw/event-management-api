package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/token"
)

// BecomeHost godoc
// @Summary      Creates a new request to become host.
// @Description  Creates a new request to become host.
// @Tags         user
// @Produce      json
// @Success      200 {object} ResponseMessage "request to become host created"
// @Failure      401 {object} ResponseMessage "Not Authorized"
// @Router       /users/host [post]
// @Security     Bearer
func (server *Server) BecomeHost(context *gin.Context) {
	// TODO: Delete old request for user if not approved after 30 days

	// Authorization to make sure payload have user role
	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	userEmail := payload.Username

	user, err := server.provider.GetUserByEmail(context, userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusUnauthorized, ResponseMessage{Message: "Not Authorized"})
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.Role != model.UserRole_User {
		context.JSON(http.StatusUnauthorized, ResponseMessage{Message: "Not Authorized"})
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

	context.JSON(http.StatusOK, ResponseMessage{Message: "request to become host created"})
}

type ListPendingUserHostRequestsParams struct {
	Limit  int `form:"limit" binding:"required,min=1,max=1000"`
	Offset int `form:"offset"`
}

// ListPendingUserHostRequests godoc
// @Summary      Lists pending requests to become host.
// @Description  Lists pending requests to become host.
// @Tags         moderator
// @Produce      json
// @Param        limit query int true "Limit"
// @Param        offset query int false "Offset"
// @Success      200 {object} model.ListPendingRequestsResponse
// @Failure      401 {object} ResponseMessage "Not Authorized"
// @Router       /moderators/requests [get]
// @Security     Bearer
func (server *Server) ListPendingUserHostRequests(context *gin.Context) {
	// Authorization to make sure payload have user role
	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	userEmail := payload.Username

	user, err := server.provider.GetUserByEmail(context, userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusUnauthorized, ResponseMessage{Message: "Not Authorized"})
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Only admin and moderator can list pending requests
	if user.Role == model.UserRole_User || user.Role == model.UserRole_Host {
		context.JSON(http.StatusUnauthorized, ResponseMessage{Message: "Not Authorized"})
		return
	}

	// Listing pending requests
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

type ApproveDisapproveUserHostRequestParams struct {
	Approved  bool  `json:"approved"`
	RequestID int64 `json:"request_id"`
}

// ApproveDisapproveUserHostRequest godoc
// @Summary      Approves or disapproves a request to become host.
// @Description  Approves or disapproves a request to become host.
// @Tags         moderator
// @Produce      json
// @Param        request body ApproveDisapproveUserHostRequestParams true "Request"
// @Success      200 {object} ResponseMessage "request approved/disapproved"
// @Failure      401 {object} ResponseMessage "Not Authorized"
// @Router       /moderators/requests [post]
// @Security     Bearer
func (server *Server) ApproveDisapproveUserHostRequest(context *gin.Context) {
	// Authorization to make sure payload have user role
	payload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	userEmail := payload.Username

	user, err := server.provider.GetUserByEmail(context, userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusUnauthorized, ResponseMessage{Message: "Not Authorized"})
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Only admin and moderator can approve requests
	if user.Role == model.UserRole_User || user.Role == model.UserRole_Host {
		context.JSON(http.StatusUnauthorized, ResponseMessage{Message: "Not Authorized"})
		return
	}

	var req ApproveDisapproveUserHostRequestParams
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dbReq := model.ApproveDisapproveRequestToBecomeHostParams{
		RequestID:   req.RequestID,
		Approved:    req.Approved,
		ModeratorID: user.ID,
	}
	err = server.provider.ApproveDisapproveRequestToBecomeHost(context, dbReq)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, ResponseMessage{"request approved/disapproved"})
}

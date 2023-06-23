package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/token"
)

type CreateEventParams struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	TotalTickets int64     `json:"total_tickets"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

func (server *Server) CreateEvent(context *gin.Context) {
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

	// Only host can create event
	if user.Role != model.UserRole_Host {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	var params CreateEventParams
	if err := context.ShouldBindJSON(&params); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: Validate params start date and end date are in the future

	event, err := server.provider.CreateEvent(context, model.CreateEventParams{
		HostID:       user.ID,
		Name:         params.Name,
		Description:  params.Description,
		Location:     params.Location,
		TotalTickets: params.TotalTickets,
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, event)
}

type ListEventsParams struct {
	Limit  int `form:"limit" binding:"required,min=1,max=1000"`
	Offset int `form:"offset"`
}

func (server *Server) ListEvents(context *gin.Context) {
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

	// Only host can list events
	if user.Role != model.UserRole_Host {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	var params ListEventsParams
	if err := context.ShouldBindQuery(&params); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	events, err := server.provider.ListEvents(context, model.ListEventsParams{
		HostID: user.ID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, events)
}

type GetEventParams struct {
	EventID int64 `uri:"event_id" binding:"required"`
}

func (server *Server) GetEvent(context *gin.Context) {
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

	// Only host can get event
	if user.Role != model.UserRole_Host {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
		return
	}

	var params GetEventParams
	if err := context.ShouldBindUri(&params); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	event, err := server.provider.GetEvent(context, model.GetEventParams{
		HostID:  user.ID,
		EventID: params.EventID,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, event)
}

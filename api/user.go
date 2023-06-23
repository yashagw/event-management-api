package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/util"
	"github.com/yashagw/event-management-api/worker"
)

type UserResponse struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordUpdatedAt time.Time `json:"password_updated_at"`
}

// CreateUserParams represents the parameters used to create a user.
type CreateUserParams struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// CreateUser             godoc
// @Summary      Creates a new user.
// @Description  Creates a new user.
// @Tags         user
// @Produce      json
// @Param        user body CreateUserParams true "User"
// @Success      201 {object} UserResponse
// @Router       /users [post]
func (server *Server) CreateUser(context *gin.Context) {
	var req CreateUserParams
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reqParams := model.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}
	user, err := server.provider.CreateUser(context, reqParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	taskPayload := worker.PayloadSendEmailVerify{
		Email: user.Email,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = server.distributor.DistributeTaskSendEmailVerify(context, &taskPayload, opts...)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	res := UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordUpdatedAt: user.PasswordUpdatedAt,
	}

	context.JSON(http.StatusCreated, res)
}

// LoginUserParams represents the parameters used to login a user.
type LoginUserParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginUserResponse represents the response to a login user request
type LoginUserResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// LoginUser              godoc
// @Summary      Logs in a user.
// @Description  Logs in a user.
// @Tags         user
// @Produce      json
// @Param        user body LoginUserParams true "User"
// @Success      200 {object} LoginUserResponse
// @Router       /users/login [post]
func (server *Server) LoginUser(context *gin.Context) {
	var req LoginUserParams
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.provider.GetUserByEmail(context, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.Email, server.config.AccessTokenDuration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := LoginUserResponse{
		Token: accessToken,
		User: UserResponse{
			ID:                user.ID,
			Name:              user.Name,
			Email:             user.Email,
			CreatedAt:         user.CreatedAt,
			PasswordUpdatedAt: user.PasswordUpdatedAt,
		},
	}

	context.JSON(http.StatusOK, res)
}

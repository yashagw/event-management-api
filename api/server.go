package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yashagw/event-management-api/db"
	"github.com/yashagw/event-management-api/token"
	"github.com/yashagw/event-management-api/util"
)

// Server will serve HTTP requests for our event service.
type Server struct {
	provider   db.Provider
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config util.Config, provider db.Provider) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		provider:   provider,
		config:     config,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.LoginUser)

	userAuthRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	userAuthRoutes.POST("/users/host", server.BecomeHost)
	userAuthRoutes.POST("/users/ticket", server.CreateTicket)

	moderatorAuthRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	moderatorAuthRoutes.GET("/moderator/requests", server.ListPendingUserHostRequests)
	moderatorAuthRoutes.POST("/moderator/requests/", server.ApproveDisapproveUserHostRequest)

	hostAuthRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	hostAuthRoutes.POST("/host/events", server.CreateEvent)
	hostAuthRoutes.GET("/host/events", server.ListEvents)
	hostAuthRoutes.GET("/host/events/:event_id", server.GetEvent)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

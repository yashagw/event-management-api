package api

import (
	"fmt"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yashagw/event-management-api/docs"
	"github.com/yashagw/event-management-api/worker"

	"github.com/gin-gonic/gin"
	"github.com/yashagw/event-management-api/db"
	"github.com/yashagw/event-management-api/token"
	"github.com/yashagw/event-management-api/util"
)

// Server will serve HTTP requests for our event service.
type Server struct {
	provider    db.Provider
	config      util.Config
	tokenMaker  token.Maker
	router      *gin.Engine
	distributor worker.TaskDistributor
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config util.Config, provider db.Provider, distributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		provider:    provider,
		config:      config,
		tokenMaker:  tokenMaker,
		distributor: distributor,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.LoginUser)

	userAuthRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	userAuthRoutes.POST("/users/host", server.BecomeHost)
	userAuthRoutes.POST("/users/ticket", server.CreateTicket)

	moderatorAuthRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	moderatorAuthRoutes.GET("/moderator/requests", server.ListPendingUserHostRequests)
	moderatorAuthRoutes.POST("/moderator/requests/", server.ApproveDisapproveUserHostRequest)

	hostAuthRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	hostAuthRoutes.POST("/hosts/events", server.CreateEvent)
	hostAuthRoutes.GET("/hosts/events", server.ListEvents)
	hostAuthRoutes.GET("/hosts/events/:event_id", server.GetEvent)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

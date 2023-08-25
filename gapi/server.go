package gapi

import (
	"fmt"

	_ "github.com/yashagw/event-management-api/docs"
	"github.com/yashagw/event-management-api/pb"
	"github.com/yashagw/event-management-api/worker"

	"github.com/yashagw/event-management-api/db"
	"github.com/yashagw/event-management-api/token"
	"github.com/yashagw/event-management-api/util"
)

// Server will serve HTTP requests for our event service.
type Server struct {
	pb.UnimplementedEventManagementServer
	provider    db.Provider
	config      util.Config
	tokenMaker  token.Maker
	distributor worker.TaskDistributor
}

// NewServer creates a new gRPC server
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

	return server, nil
}

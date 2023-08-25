package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/hibiken/asynq"
	"github.com/yashagw/event-management-api/api"
	"github.com/yashagw/event-management-api/db"
	"github.com/yashagw/event-management-api/gapi"
	"github.com/yashagw/event-management-api/pb"
	"github.com/yashagw/event-management-api/util"
	"github.com/yashagw/event-management-api/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runTaskProcessor(redisOpt asynq.RedisClientOpt, provider db.Provider) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, provider)
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("cannot start task processor:", err)
	}
}

// @title     Event Mangement API
// @version	  1.0
// @description	API server for event management system.
// @contact.name	Yash Agarwal
// @contact.email	yash.ag@outlook.com
// @securityDefinitions.apikey Bearer
// @in header
// @name authorization
// @description Type "bearer" followed by a space and JWT token.
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	provider, err := db.New(conn)
	if err != nil {
		log.Fatal("cannot create db provider:", err)
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	// go runTaskProcessor(redisOpt, provider)

	runGrpcServer(config, provider, taskDistributor)
}

func runGrpcServer(config util.Config, provider db.Provider, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, provider, taskDistributor)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEventManagementServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listerner:", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runGinServer(config util.Config, provider db.Provider, taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(config, provider, taskDistributor)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

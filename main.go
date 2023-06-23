package main

import (
	"database/sql"
	"log"

	"github.com/hibiken/asynq"
	"github.com/yashagw/event-management-api/api"
	"github.com/yashagw/event-management-api/db"
	"github.com/yashagw/event-management-api/util"
	"github.com/yashagw/event-management-api/worker"
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
	go runTaskProcessor(redisOpt, provider)

	server, err := api.NewServer(config, provider, taskDistributor)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

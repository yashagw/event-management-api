package main

import (
	"database/sql"
	"log"

	"github.com/yashagw/event-management-api/api"
	"github.com/yashagw/event-management-api/db"
	"github.com/yashagw/event-management-api/util"
)

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

	server, err := api.NewServer(config, provider)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

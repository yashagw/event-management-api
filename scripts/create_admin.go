package scripts

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/util"
)

func CreateAdmin() {
	name := "admin"
	email := "admin@gmai.com"
	password := "secret1234"
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		log.Fatal("cannot hash password:", err)
	}

	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	_, err = conn.Exec("INSERT INTO users (name, email, hashed_password, role) VALUES ($1, $2, $3, $4)",
		name, email, hashedPassword, model.UserRole_Admin)
	if err != nil {
		log.Fatal("cannot create admin user:", err)
	}

	fmt.Println("admin user created successfully with email:", email, "and password:", password, "")
}

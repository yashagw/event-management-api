package pgsql

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/yashagw/event-management-api/util"
)

var provider *Provider

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	provider, err = NewProvider(testDB)
	if err != nil {
		log.Fatal("cannot create provider:", err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

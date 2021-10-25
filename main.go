package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/xiahezzz/simplebank/api"
	db "github.com/xiahezzz/simplebank/db/sqlc"
	"github.com/xiahezzz/simplebank/db/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config in main:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("server cannot start:", err)
	}
}

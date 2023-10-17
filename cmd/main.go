package main

import (
	"database/sql"
	"github.com/aalug/cv-backend-go/internal/api"
	"github.com/aalug/cv-backend-go/internal/config"
	db "github.com/aalug/cv-backend-go/internal/db/sqlc"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env file: ", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}

	store := db.NewStore(conn)

	// @BasePath /api/v1
	// @contact.name aalug
	// @contact.url https://github.com/aalug
	// @contact.email a.a.gulczynski@gmail.com

	server := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server:", err)
	}
}

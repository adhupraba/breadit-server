package lib

import (
	"database/sql"
	"log"

	"github.com/adhupraba/breadit-server/internal/database"
)

var DB *database.Queries

func ConnectDb() {
	log.Println("db url", EnvConfig.DbUrl)
	conn, err := sql.Open("postgres", EnvConfig.DbUrl)

	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	DB = database.New(conn)
}

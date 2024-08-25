package database

import (
	"context"
	database_connector "github.com/devyk100/gengou-db/pkg/database"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

var Conn *pgx.Conn
var err error
var Queries *database_connector.Queries
var ctx = context.Background()

func DbInit() {

	dsn := os.Getenv("DATABASE_URL")
	Conn, err = pgx.Connect(ctx, dsn)
	if err != nil {
		panic(err.Error())
	}
	Queries = database_connector.New(Conn)
	log.Println("Connected to database")
}

func DbClose() {
	err := Conn.Close(context.Background())
	if err != nil {
		return
	}
}

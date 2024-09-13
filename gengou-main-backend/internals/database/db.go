package database

import (
	"context"
	database_connector "github.com/devyk100/gengou-db/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var Conn *pgx.Conn
var err error
var Queries *database_connector.Queries
var ctx = context.Background()

func DbInit() {

	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	// Adjust pool settings if needed
	poolConfig.MaxConns = 20 // Example: Set maximum number of connections
	poolConfig.MinConns = 5  // Example: Set minimum number of connections

	// Create connection pool
	Conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	//defer Conn.Close() // Ensure the pool is closed when you're done

	Queries = database_connector.New(Conn)
	log.Println("Connected to database")
}

func DbClose() {
	err := Conn.Close(context.Background())
	if err != nil {
		return
	}
}

package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPostgres() (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(context.Background(), "postgresql://root:1234567890@localhost:5432/autom?sslmode=disable")
	return
}

func InitPgxPool() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databaseUrl := "postgresql://root:1234567890@localhost:5432/autom?sslmode=disable"

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatalf("Unable to parse DB config: %v", err)
	}

	config.MaxConns = 10 // limit number of connections in pool

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Optional: Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Connected to Postgres with pgxpool")
	return pool
}

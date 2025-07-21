package postgres

import (
	"context"
	"log"
	"time"

	"github.com/b0nbon1/stratal/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPostgres(cfg *config.Config) (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(context.Background(), cfg.Database.DSN())
	return
}

func InitPgxPool(cfg *config.Config) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(cfg.Database.DSN())
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

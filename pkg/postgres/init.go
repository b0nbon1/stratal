package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func InitPostgres() (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	return
}

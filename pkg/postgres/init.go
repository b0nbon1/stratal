package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func InitPostgres() (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(context.Background(), "postgresql://root:1234567890@localhost:5432/autom?sslmode=disable")
	return
}

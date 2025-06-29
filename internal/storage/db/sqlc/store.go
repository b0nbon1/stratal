package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	CreateJobWithTasksTx(ctx context.Context, jobParams CreateJobParams, taskInputs []CreateTaskParams) (JobWithTaskResult, error)
	CreateJobRunTx(ctx context.Context, jobID pgtype.UUID, triggeredBy string) (JobRunResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

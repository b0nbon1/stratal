package main

import (
	"context"
	"log"

	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/b0nbon1/stratal/internal/postgres"
	"github.com/b0nbon1/stratal/internal/queue"
	"github.com/b0nbon1/stratal/internal/scheduler"
	"github.com/b0nbon1/stratal/internal/worker"
	"github.com/b0nbon1/stratal/cmd/server"
)

func main() {
	conn, err := postgres.InitPostgres()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer conn.Close(context.Background())
	store := db.New(conn)
	q := queue.NewRedisQueue("localhost:6379", "", 0, "jobs")
	scheduler.StartScheduler(q, *store)
	worker.StartWorker(q)

	server := server.NewServer(store)
	err = server.Start(":8040")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

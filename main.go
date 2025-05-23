package main

import (
	"context"
	"log"

	"github.com/b0nbon1/stratal/api"
	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/postgres"
	"github.com/b0nbon1/stratal/pkg/queue"
	"github.com/b0nbon1/stratal/pkg/scheduler"
	"github.com/b0nbon1/stratal/pkg/worker"
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

	server := api.NewServer(*store)
	err = server.Start(":8040")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

package main

import (
	"log"

	"github.com/b0nbon1/temporal-lite/cmd"
	db "github.com/b0nbon1/temporal-lite/db/sqlc"
	"github.com/b0nbon1/temporal-lite/pkg/postgres"
	"github.com/b0nbon1/temporal-lite/pkg/queue"
	"github.com/b0nbon1/temporal-lite/pkg/scheduler"
	"github.com/b0nbon1/temporal-lite/pkg/worker"
)

func main() {
	conn, err := postgres.InitPostgres()
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.New(conn)
	q := queue.NewRedisQueue("localhost:6379", "", 0, "jobs")
	scheduler.StartScheduler(q, store)
	worker.StartWorker(q)

	
	server := cmd.NewServer()
	err = server.Start(":8080")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}


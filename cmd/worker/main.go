package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/b0nbon1/stratal/internal/queue"
	psql "github.com/b0nbon1/stratal/internal/storage/db"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/internal/worker"
)

func main() {

    conn, err := psql.InitPostgres()
    if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

    ctx := context.Background()
    defer conn.Close(ctx)
    fmt.Println("Connected to database successfully")

    store := db.New(conn)
	q := queue.NewRedisQueue("localhost:6379", "", 0, "job_runs")


    go worker.StartWorker(q, *store)
    fmt.Println("Worker started successfully")

    quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
    <-quitChannel

    fmt.Println("Stopped by signal, exiting gracefully...")
}

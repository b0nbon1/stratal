package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/b0nbon1/stratal/internal/queue"
	"github.com/b0nbon1/stratal/internal/scheduler"
	psql "github.com/b0nbon1/stratal/internal/storage/db"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/internal/worker"
)

func main() {
    ctx := context.Background()

    fmt.Println("Connected to database successfully")

	pool := psql.InitPgxPool()
	defer pool.Close()

    store := db.NewStore(pool)
	q := queue.NewRedisQueue("localhost:6379", "", 0, "job_runs")

	scheduler.StartScheduler(q, store.(*db.SQLStore), ctx)

    go worker.StartWorker(ctx, q, store.(*db.SQLStore))
    fmt.Println("Worker started successfully")

    quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
    <-quitChannel

    fmt.Println("Stopped by signal, exiting gracefully...")
}

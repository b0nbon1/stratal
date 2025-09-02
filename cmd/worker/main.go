package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/b0nbon1/stratal/internal/config"
	"github.com/b0nbon1/stratal/internal/queue"
	"github.com/b0nbon1/stratal/internal/scheduler"
	"github.com/b0nbon1/stratal/internal/security"
	psql "github.com/b0nbon1/stratal/internal/storage/db"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/internal/worker"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg := config.Load()

	// Initialize secret manager
	secretManager, err := security.NewSecretManager(cfg.Security.EncryptionKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize secret manager: %v", err))
	}

	fmt.Println("Connected to database successfully")

	pool := psql.InitPgxPool(cfg)
	defer pool.Close()

	store := db.NewStore(pool)
	q := queue.NewRedisQueue(cfg, "job_runs")

	go scheduler.StartScheduler(q, store.(*db.SQLStore), ctx)

	go worker.StartWorker(ctx, q, store.(*db.SQLStore), secretManager)
	fmt.Println("Worker started successfully")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	fmt.Println("Stopped by signal, exiting gracefully...")
}

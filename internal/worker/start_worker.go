package worker

import (
	"context"
	"fmt"

	"github.com/b0nbon1/stratal/internal/logger"
	"github.com/b0nbon1/stratal/internal/queue"
	"github.com/b0nbon1/stratal/internal/security"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

type Worker struct {
	ctx           context.Context
	q             queue.TaskQueue
	store         *db.SQLStore
	secretManager *security.SecretManager
	logSystem     *logger.Logger
}

func StartWorker(ctx context.Context, q queue.TaskQueue, store *db.SQLStore, secretManager *security.SecretManager) {
	logSystem := logger.NewLogger(store, "internal/storage/files/logs")
	defer logSystem.Close()

	worker := &Worker{
		ctx:           ctx,
		q:             q,
		store:         store,
		secretManager: secretManager,
		logSystem:     logSystem,
	}

	go worker.Start()
}

func (w *Worker) Start() {
	fmt.Println("Starting worker...")
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println("Worker context cancelled, shutting down...")
			return
		default:
			fmt.Println("Processing next job...")
			w.ProcessNextJob()
			fmt.Println("Job processed")
		}
	}
}

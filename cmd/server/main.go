package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/b0nbon1/stratal/internal/api"
	"github.com/b0nbon1/stratal/internal/config"
	"github.com/b0nbon1/stratal/internal/queue"
	"github.com/b0nbon1/stratal/internal/security"
	postgres "github.com/b0nbon1/stratal/internal/storage/db"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	_ "net/http/pprof"
)

func main() {
	cfg := config.Load()

	secretManager, err := security.NewSecretManager(cfg.Security.EncryptionKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize secret manager: %v", err))
	}

	pool := postgres.InitPgxPool(cfg)
	defer pool.Close()
	store := db.NewStore(pool)

	q := queue.NewRedisQueue(cfg, "job_runs", "workers", 3)

	hs := api.NewHTTPServer(cfg.Server.Address(), store.(*db.SQLStore), q, secretManager)

	if err := hs.Start(); err != nil {
		panic(err)
	}
	defer hs.Stop()

	fmt.Printf("Server running at http://%s\n", cfg.Server.Address())
	if err := hs.Server.ListenAndServe(); err != nil {
		panic(err)
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Stopped by signal, exiting gracefully...")
}

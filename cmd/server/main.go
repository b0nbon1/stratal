package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/b0nbon1/stratal/internal/api"
	"github.com/b0nbon1/stratal/internal/queue"
	postgres "github.com/b0nbon1/stratal/internal/storage/db"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

func main() {

	pool := postgres.InitPgxPool()
	defer pool.Close()

	q := queue.NewRedisQueue("localhost:6379", "", 0, "job_runs")

	// Now use pool with sqlc or raw queries
	store := db.NewStore(pool)

	hs := api.NewHTTPServer(":8080", store.(*db.SQLStore), q)

	if err := hs.Start(); err != nil {
		panic(err)
	}
	defer hs.Stop()

	fmt.Println("Server running at http://localhost:8080")
	if err := hs.Server.ListenAndServe(); err != nil {
		panic(err)
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Stopped by signal, exiting gracefully...")


}
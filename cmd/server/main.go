package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/b0nbon1/stratal/internal/api"
)

func main() {

	hs := api.NewHTTPServer(":8080", nil)

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
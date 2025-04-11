package main

import (
	"log"

	"github.com/b0nbon1/temporal-lite/cmd"
)

func main() {
	
	server := cmd.NewServer()
	err := server.Start(":8080")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}


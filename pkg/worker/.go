package worker

import (
	"fmt"
	"log"

	"github.com/b0nbon1/temporal-lite/pkg/queue"
	"github.com/b0nbon1/temporal-lite/utils"
)

func StartWorker(q queue.TaskQueue) {
	go func() {
		for {
			tasks, err := q.XReadGeneric("0", 5000, utils.TaskMapper)
			if err != nil {
				log.Println("Error reading from queue:", err)
				continue
			}
			for job := range tasks {
				fmt.Println(job)
			}
		}
	}()
	
}

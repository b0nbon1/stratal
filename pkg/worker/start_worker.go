package worker

import (
	"fmt"
	"time"

	// "log"
	// "time"

	"github.com/b0nbon1/stratal/internal/processor"
	"github.com/b0nbon1/stratal/pkg/queue"
	// "github.com/b0nbon1/stratal/utils"
)

func StartWorker(q queue.TaskQueue) {
	fmt.Println("Starting worker...")
	go func() {
		for {
			fmt.Println("Worker started ==============")
			task, err := q.Dequeue()
			if err != nil {
				fmt.Println("Error dequeuing task:", err)
				continue
			}
			// fmt.Println("Dequeued task:", task)
			// Process the task here
			// For example, you can print the task or perform some operations on it
			processor.ProcessJob(task)
			fmt.Println("Processing task:", task)
			// Simulate some processing time
			time.Sleep(1 * time.Second)
		}
	}()

}

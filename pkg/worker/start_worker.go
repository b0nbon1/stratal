package worker

import (
	"encoding/json"
	"fmt"

	db "github.com/b0nbon1/stratal/db/sqlc"
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
			var job db.Job
			err = json.Unmarshal([]byte(task), &job)
			if err != nil {
				fmt.Println("Failed to unmarshal:", err)
				continue
			}

			processor.ProcessJob(job)
		}
	}()

}

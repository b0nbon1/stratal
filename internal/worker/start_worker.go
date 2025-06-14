package worker

import (
	"fmt"

	"github.com/b0nbon1/stratal/internal/queue"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

func StartWorker(q queue.TaskQueue, store db.Queries) {
	fmt.Println("Starting worker...")
	for {
		fmt.Println("Worker started ==============")
		job_run_id, err := q.Dequeue()
		if err != nil {
			fmt.Println("Error dequeuing task:", err)
			continue
		}

		fmt.Println("Dequeued job run ID:", job_run_id)

		

		fmt.Println(job_run_id)
	}
}

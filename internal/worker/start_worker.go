package worker

import (
	"context"
	"fmt"

	"github.com/b0nbon1/stratal/internal/queue"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
)

func StartWorker(ctx context.Context, q queue.TaskQueue, store db.Queries) {
	fmt.Println("Starting worker...")
	for {
		fmt.Println("Worker started ==============")
		jobRunId, err := q.Dequeue()
		if err != nil {
			fmt.Println("Error dequeuing job_run:", err)
			continue
		}

		jobRunIdUUID, err := utils.ParseUUID(jobRunId)
		if err != nil {
			fmt.Println("Error parsing job_run_id:", err)
			continue
		}
		

		jobRun, err := store.GetJobRun(ctx, jobRunIdUUID)
		if err != nil {
			fmt.Println("Error getting job_run:", err)
			continue
		}


		fmt.Println(jobRun)
	}
}

package worker

import (
	"context"
	"fmt"

	// "github.com/b0nbon1/stratal/internal/processor"
	"github.com/b0nbon1/stratal/internal/queue"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
)

func StartWorker(ctx context.Context, q queue.TaskQueue, store *db.SQLStore) {
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

		if jobRun.Status.String != "pending" {
			fmt.Println("Job run is not pending, skipping:", jobRunId)
			continue
		}

		// update the status of job and queue
		// updateParams := db.UpdateJobRunStatusParams{
		// 	ID:     jobRunIdUUID,
		// 	Status: utils.ParseText("running"),
		// }
		// err = store.UpdateJobRunStatus(ctx, updateParams)
		// if err != nil {
		// 	fmt.Println("Error updating job_run status to running:", err)
		// 	continue
		// }
		// fmt.Println("Job run status updated to running:", jobRunId)

		// Fetch the job associated with the job run with tasks
		job, err := store.GetJobWithTasks(ctx, jobRun.JobID)
		if err != nil {
			fmt.Println("Error getting job with tasks:", err)
			continue
		}

		fmt.Println("Fetched job with tasks:", job)

		// // Process the job run
		// err = processor.ProcessJob(job.Tasks)
		// if err != nil {	
		// 	fmt.Println("Error processing job_run:", err)

		// 	continue
		// }
	}
}

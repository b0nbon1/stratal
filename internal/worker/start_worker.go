package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/b0nbon1/stratal/internal/processor"
	"github.com/b0nbon1/stratal/internal/queue"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func StartWorker(ctx context.Context, q queue.TaskQueue, store *db.SQLStore) {
	fmt.Println("Starting worker...")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker context cancelled, shutting down...")
			return
		default:
			fmt.Println("Processing next job...")
			processNextJob(ctx, q, store)
			fmt.Println("Job processed")
		}
	}
}

func processNextJob(ctx context.Context, q queue.TaskQueue, store *db.SQLStore) {
	fmt.Println("Worker polling for jobs...")

	jobRunId, err := q.Dequeue()
	if err != nil {
		fmt.Printf("Error dequeuing job_run: %v\n", err)
		time.Sleep(2 * time.Second) // Back off on error
		return
	}

	if jobRunId == "" {
		fmt.Println("No job_run to process")
		time.Sleep(2 * time.Second) // Back off on error
		return
	}

	jobRunIdUUID, err := utils.ParseUUID(jobRunId)
	if err != nil {
		fmt.Printf("Error parsing job_run_id %s: %v\n", jobRunId, err)
		return
	}

	// Process the job
	if err := processJobRun(ctx, store, jobRunIdUUID); err != nil {
		fmt.Printf("Error processing job_run %s: %v\n", jobRunId, err)
		// Job stays in failed state, could implement retry logic here
		updateJobRunError(ctx, store, jobRunIdUUID, "Failed to process job", err)
	}
}

func processJobRun(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID) error {
	// Fetch job run
	jobRun, err := store.GetJobRun(ctx, jobRunID)
	if err != nil {
		return fmt.Errorf("failed to get job_run: %w", err)
	}

	// Verify status
	if jobRun.Status.String != "queued" && jobRun.Status.String != "pending" {
		return fmt.Errorf("job run %s is not in queued/pending state, current state: %s",
			jobRunID.String(), jobRun.Status.String)
	}

	// Update status to running
	startTime := pgtype.Timestamp{Time: time.Now(), Valid: true}
	err = store.UpdateJobRun(ctx, db.UpdateJobRunParams{
		ID:          jobRunID,
		Status:      pgtype.Text{String: "running", Valid: true},
		StartedAt:   startTime,
		TriggeredBy: jobRun.TriggeredBy,
		Metadata:    jobRun.Metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to update job_run status to running: %w", err)
	}

	fmt.Printf("Job run %s status updated to running\n", jobRunID.String())

	// Fetch the job with tasks
	job, err := store.GetJobWithTasks(ctx, jobRun.JobID)
	if err != nil {
		updateJobRunError(ctx, store, jobRunID, "Failed to fetch job details", err)
		return fmt.Errorf("failed to get job with tasks: %w", err)
	}

	// Process the job
	err = processor.ProcessJob(ctx, store, jobRunID, job)

	// Update final status
	finishTime := pgtype.Timestamp{Time: time.Now(), Valid: true}
	finalStatus := "completed"
	var errorMsg pgtype.Text

	if err != nil {
		finalStatus = "failed"
		errorMsg = pgtype.Text{String: err.Error(), Valid: true}
		fmt.Printf("Job run %s failed: %v\n", jobRunID.String(), err)
	} else {
		fmt.Printf("Job run %s completed successfully\n", jobRunID.String())
	}

	updateErr := store.UpdateJobRun(ctx, db.UpdateJobRunParams{
		ID:           jobRunID,
		Status:       pgtype.Text{String: finalStatus, Valid: true},
		StartedAt:    startTime,
		FinishedAt:   finishTime,
		ErrorMessage: errorMsg,
		TriggeredBy:  jobRun.TriggeredBy,
		Metadata:     jobRun.Metadata,
	})

	if updateErr != nil {
		fmt.Printf("Failed to update final job_run status: %v\n", updateErr)
	}

	return err
}

func updateJobRunError(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID, message string, err error) {
	fullError := fmt.Sprintf("%s: %v", message, err)
	updateErr := store.UpdateJobRunError(ctx, db.UpdateJobRunErrorParams{
		ID:           jobRunID,
		ErrorMessage: pgtype.Text{String: fullError, Valid: true},
	})
	if updateErr != nil {
		fmt.Printf("Failed to update job_run error: %v\n", updateErr)
	}
}

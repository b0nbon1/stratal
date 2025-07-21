package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/b0nbon1/stratal/internal/logger"
	"github.com/b0nbon1/stratal/internal/processor"
	"github.com/b0nbon1/stratal/internal/queue"
	"github.com/b0nbon1/stratal/internal/security"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func StartWorker(ctx context.Context, q queue.TaskQueue, store *db.SQLStore) {
	StartWorkerWithSecrets(ctx, q, store, nil)
}

func StartWorkerWithSecrets(ctx context.Context, q queue.TaskQueue, store *db.SQLStore, secretManager *security.SecretManager) {
	// Initialize the logger system
	logSystem := logger.NewLogger(store, "internal/storage/files/logs")
	defer logSystem.Close()

	fmt.Println("Starting worker...")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker context cancelled, shutting down...")
			return
		default:
			fmt.Println("Processing next job...")
			processNextJobWithSecrets(ctx, q, store, secretManager, logSystem)
			fmt.Println("Job processed")
		}
	}
}

func processNextJob(ctx context.Context, q queue.TaskQueue, store *db.SQLStore) {
	processNextJobWithSecrets(ctx, q, store, nil, nil)
}

func processNextJobWithSecrets(ctx context.Context, q queue.TaskQueue, store *db.SQLStore, secretManager *security.SecretManager, logSystem *logger.Logger) {
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

	// Get or create a logger for this job run
	var jobLogger *logger.JobRunLogger
	if logSystem != nil {
		jobLogger, err = logSystem.GetJobRunLogger(jobRunId)
		if err != nil {
			fmt.Printf("Error creating logger for job_run %s: %v\n", jobRunId, err)
			// Continue processing even if logger creation fails
		}
	}

	// Log job run start
	if jobLogger != nil {
		jobLogger.Info(fmt.Sprintf("Starting job run %s", jobRunId))
	}

	// Process the job
	if err := processJobRunWithSecrets(ctx, store, jobRunIdUUID, secretManager, jobLogger); err != nil {
		fmt.Printf("Error processing job_run %s: %v\n", jobRunId, err)
		if jobLogger != nil {
			jobLogger.Error(fmt.Sprintf("Job run failed: %v", err))
		}
		// Job stays in failed state, could implement retry logic here
		updateJobRunError(ctx, store, jobRunIdUUID, "Failed to process job", err)
	} else {
		if jobLogger != nil {
			jobLogger.Info(fmt.Sprintf("Job run %s completed successfully", jobRunId))
		}
	}

	// Close the logger for this job run when done
	if logSystem != nil {
		logSystem.CloseJobRunLogger(jobRunId)
	}
}

func processJobRun(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID) error {
	return processJobRunWithSecrets(ctx, store, jobRunID, nil, nil)
}

func processJobRunWithSecrets(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID, secretManager *security.SecretManager, jobLogger *logger.JobRunLogger) error {
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

	// Process the job using the processor
	if secretManager != nil {
		return processor.ProcessJobWithSecrets(ctx, store, secretManager, jobRunID, job, jobLogger)
	} else {
		return processor.ProcessJob(ctx, store, jobRunID, job, jobLogger)
	}
}

func updateJobRunError(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID, message string, err error) {
	errorMessage := fmt.Sprintf("%s: %v", message, err)
	updateErr := store.UpdateJobRunError(ctx, db.UpdateJobRunErrorParams{
		ID:           jobRunID,
		ErrorMessage: utils.ParseText(errorMessage),
	})
	if updateErr != nil {
		fmt.Printf("Failed to update job run error: %v\n", updateErr)
	}
}

package worker

import (
	"fmt"
	"time"

	"github.com/b0nbon1/stratal/internal/logger"
	"github.com/b0nbon1/stratal/internal/processor"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

 func (w *Worker) ProcessNextJob() {
	fmt.Println("Worker polling for jobs...")

	msgID, values, err := w.q.Dequeue(5 * time.Second)
    if err != nil {
        fmt.Println("Error dequeue:", err)
        return
    }
    if msgID == "" {
        return
    }

    jobRunId := values["job_run_id"].(string)
    fmt.Println("Processing job:", jobRunId)

	if jobRunId == "" {
		fmt.Println("No job_run to process")
		time.Sleep(2 * time.Second)
		return
	}

	jobRunIdUUID, err := utils.ParseUUID(jobRunId)
	if err != nil {
		fmt.Printf("Error parsing job_run_id %s: %v\n", jobRunId, err)
		return
	}

	var jobLogger *logger.JobRunLogger
	if w.logSystem != nil {
		jobLogger, err = w.logSystem.GetJobRunLogger(jobRunId)
		if err != nil {
			fmt.Printf("Error creating logger for job_run %s: %v\n", jobRunId, err)
		}
	}

	if jobLogger != nil {
		jobLogger.Info(fmt.Sprintf("Starting job run %s", jobRunId))
	}

	if err := w.ProcessJobRun(jobRunIdUUID, jobLogger); err != nil {
		fmt.Printf("Error processing job_run %s: %v\n", jobRunId, err)
		if jobLogger != nil {
			jobLogger.Error(fmt.Sprintf("Job run failed: %v", err))
		}
		// Job stays in failed state, could implement retry logic here
		w.UpdateJobRunError(jobRunIdUUID, "Failed to process job", err)
	} else {
		if jobLogger != nil {
			jobLogger.Info(fmt.Sprintf("Job run %s completed successfully", jobRunId))
		}
	}

	if w.logSystem != nil {
		w.logSystem.CloseJobRunLogger(jobRunId)
	}
}

func (w *Worker) ProcessJobRun(jobRunID pgtype.UUID, jobLogger *logger.JobRunLogger) error {
	jobRun, err := w.store.GetJobRun(w.ctx, jobRunID)
	if err != nil {
		return fmt.Errorf("failed to get job_run: %w", err)
	}

	if jobRun.Status.String != "queued" && jobRun.Status.String != "pending" {
		if jobRun.Status.String == "paused" {
			fmt.Printf("Job run %s is paused, skipping processing\n", jobRunID.String())
			return nil
		}
		return fmt.Errorf("job run %s is not in queued/pending/paused state, current state: %s",
			jobRunID.String(), jobRun.Status.String)
	}

	startTime := pgtype.Timestamp{Time: time.Now(), Valid: true}
	err = w.store.UpdateJobRun(w.ctx, db.UpdateJobRunParams{
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

	job, err := w.store.GetJobWithTasks(w.ctx, jobRun.JobID)
	if err != nil {
		w.UpdateJobRunError(jobRunID, "Failed to fetch job details", err)
		return fmt.Errorf("failed to get job with tasks: %w", err)
	}

	if w.secretManager != nil {
		return processor.ProcessJob(w.ctx, w.store, w.secretManager, jobRunID, job, jobLogger)
	} else {
		return processor.ProcessJob(w.ctx, w.store, nil, jobRunID, job, jobLogger)
	}
}

func (w *Worker) UpdateJobRunError(jobRunID pgtype.UUID, message string, err error) {
	errorMessage := fmt.Sprintf("%s: %v", message, err)
	updateErr := w.store.UpdateJobRunError(w.ctx, db.UpdateJobRunErrorParams{
		ID:           jobRunID,
		ErrorMessage: utils.ParseText(errorMessage),
	})
	if updateErr != nil {
		fmt.Printf("Failed to update job run error: %v\n", updateErr)
	}
}
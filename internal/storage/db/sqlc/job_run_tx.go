package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type JobRunResult struct {
	JobRunId   string
	TaskRunIds []string
}

func (store *SQLStore) CreateJobRunTx(ctx context.Context, jobID pgtype.UUID, triggeredBy string) (JobRunResult, error) {
	var result JobRunResult

	err := store.execTx(ctx, func(q *Queries) error {
		jobRun, err := q.CreateJobRun(ctx, CreateJobRunParams{
			JobID:       jobID,
			TriggeredBy: pgtype.Text{String: triggeredBy},
			Status:      pgtype.Text{String: "pending"},
		})
		if err != nil {
			return fmt.Errorf("unable to create job_run %w", err)
		}

		// Step 2: Fetch all tasks for the job
		tasks, err := q.GetTasksByJobID(ctx, jobID)
		if err != nil {
			return fmt.Errorf("failed to get tasks: %w", err)
		}

		// Step 3: Create task_runs
		var taskRunIDs []string
		for _, task := range tasks {
			taskRun, err := q.CreateTaskRun(ctx, CreateTaskRunParams{
				JobRunID: jobRun.ID,
				TaskID:   task.ID,
				Status:   pgtype.Text{String: "pending"},
			})
			if err != nil {
				return fmt.Errorf("failed to create task run: %w", err)
			}
			taskRunIDs = append(taskRunIDs, taskRun.ID.String())
		}

		result.JobRunId = jobID.String()
		result.TaskRunIds = taskRunIDs
		return nil
	})

	return result, err
}

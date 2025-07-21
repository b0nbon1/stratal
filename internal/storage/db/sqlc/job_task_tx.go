package db

import (
	"context"
	"fmt"
)

type JobWithTaskResult struct {
	Job   *CreateJobRow
	Tasks []CreateTaskRow
}

func (store *SQLStore) CreateJobWithTasksTx(
	ctx context.Context,
	jobParams CreateJobParams,
	taskInputs []CreateTaskParams,
) (JobWithTaskResult, error) {

	var result JobWithTaskResult

	err := store.execTx(ctx, func(q *Queries) error {
		job, err := q.CreateJob(ctx, jobParams)
		if err != nil {
			fmt.Println("create job failed: %w", err)
			return fmt.Errorf("create job failed: %w", err)
		}

		result.Job = &job

		for _, t := range taskInputs {
			t.JobID = job.ID
			task, err := q.CreateTask(ctx, t)
			if err != nil {
				return fmt.Errorf("task '%s' failed: %w", t.Name, err)
			}
			result.Tasks = append(result.Tasks, task)
		}

		return nil

	})

	return result, err
}

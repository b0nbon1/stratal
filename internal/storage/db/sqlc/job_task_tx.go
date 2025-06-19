package db

import (
	"context"
	"fmt"
	"sync"
)

type JobWithTaskResult struct {
	Job *CreateJobRow
	Tasks []CreateTaskRow
}

func (store *SQLStore) CreateJobWithTasksTx(
	ctx context.Context,
	jobParams CreateJobParams,
	taskInputs []CreateTaskParams,
) (JobWithTaskResult, error) {

	var result JobWithTaskResult

	err := store.execTx(ctx, func(q *Queries) error {
	// Step 1: Create the job
	job, err := q.CreateJob(ctx, jobParams)
	if err != nil {
		return fmt.Errorf("create job failed: %w", err)
	}

	// Step 2: Prepare to create tasks concurrently
	const maxConcurrent = 5
	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	taskChan := make(chan CreateTaskRow, len(taskInputs))
	errChan := make(chan error, len(taskInputs))

	for _, taskInput := range taskInputs {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(t CreateTaskParams) {
			defer wg.Done()
			defer func() { <-semaphore }()

			t.JobID = job.ID

			task, err := q.CreateTask(ctx, t)
			if err != nil {
				errChan <- fmt.Errorf("task '%s' failed: %w", t.Name, err)
				return
			}

			taskChan <- task
		}(taskInput)
	}

	wg.Wait()
	close(taskChan)
	close(errChan)

	if len(errChan) > 0 {
		return fmt.Errorf("at least one task failed: %v", <-errChan)
	}

	// Collect successful tasks
	var createdTasks []CreateTaskRow
	for t := range taskChan {
		createdTasks = append(createdTasks, t)
	}

	result.Tasks = createdTasks

	return nil

	})


	return result, err
}

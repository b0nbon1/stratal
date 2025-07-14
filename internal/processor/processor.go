package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func ProcessJob(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID, job db.GetJobWithTasksRow) error {
	fmt.Printf("Processing job run: %s for job: %s\n", jobRunID.String(), job.ID.String())

	// Parse tasks from JSON
	var tasks []db.Task
	if err := json.Unmarshal(job.Tasks, &tasks); err != nil {
		return fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks to process")
		return nil
	}

	// Get all task runs for this job run
	taskRuns, err := store.ListTaskRuns(ctx, jobRunID)
	if err != nil {
		return fmt.Errorf("failed to list task runs: %w", err)
	}

	// Create map of task ID to task run for easy lookup
	taskRunMap := make(map[string]db.ListTaskRunsRow)
	for _, tr := range taskRuns {
		taskRunMap[tr.TaskID.String()] = tr
	}

	// Sort tasks by dependencies
	sortedTasks, err := utils.TopoSort(tasks)
	if err != nil {
		return fmt.Errorf("failed to sort tasks: %w", err)
	}

	fmt.Printf("Executing %d tasks in dependency order\n", len(sortedTasks))

	// Track completed tasks for dependency checking
	completedTasks := make(map[string]bool)

	// Execute tasks in order
	for _, task := range sortedTasks {
		// Check if all dependencies are completed
		if !areDependenciesMet(task, completedTasks) {
			return fmt.Errorf("task %s dependencies not met", task.Name)
		}

		taskRun, exists := taskRunMap[task.ID.String()]
		if !exists {
			return fmt.Errorf("task run not found for task %s", task.ID.String())
		}

		// Update task run status to running
		startTime := pgtype.Timestamp{Time: time.Now(), Valid: true}
		err = store.UpdateTaskRun(ctx, db.UpdateTaskRunParams{
			ID:        taskRun.ID,
			Status:    pgtype.Text{String: "running", Valid: true},
			StartedAt: startTime,
		})
		if err != nil {
			fmt.Printf("Failed to update task run status to running: %v\n", err)
		}

		fmt.Printf("Executing task: %s (type: %s)\n", task.Name, task.Type)

		// Execute the task
		output, taskErr := ExecuteTask(ctx, task)

		// Update task run with results
		finishTime := pgtype.Timestamp{Time: time.Now(), Valid: true}
		var status string
		var errorMsg pgtype.Text
		var exitCode pgtype.Int4

		if taskErr != nil {
			status = "failed"
			errorMsg = pgtype.Text{String: taskErr.Error(), Valid: true}
			exitCode = pgtype.Int4{Int32: 1, Valid: true}
			fmt.Printf("Task %s failed: %v\n", task.Name, taskErr)
		} else {
			status = "completed"
			exitCode = pgtype.Int4{Int32: 0, Valid: true}
			completedTasks[task.ID.String()] = true
			fmt.Printf("Task %s completed successfully\n", task.Name)
		}

		updateErr := store.UpdateTaskRun(ctx, db.UpdateTaskRunParams{
			ID:           taskRun.ID,
			Status:       pgtype.Text{String: status, Valid: true},
			StartedAt:    startTime,
			FinishedAt:   finishTime,
			ExitCode:     exitCode,
			Output:       pgtype.Text{String: output, Valid: true},
			ErrorMessage: errorMsg,
		})

		if updateErr != nil {
			fmt.Printf("Failed to update task run result: %v\n", updateErr)
		}

		// If task failed, stop processing (you can make this configurable)
		if taskErr != nil {
			return fmt.Errorf("task %s failed, stopping job execution: %w", task.Name, taskErr)
		}
	}

	return nil
}

func areDependenciesMet(task db.Task, completedTasks map[string]bool) bool {
	if task.Config.DependsOn == nil || len(task.Config.DependsOn) == 0 {
		return true
	}

	for _, depID := range task.Config.DependsOn {
		if !completedTasks[depID] {
			return false
		}
	}
	return true
}

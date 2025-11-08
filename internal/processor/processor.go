package processor

import (
	"context"
	"encoding/json"
	"fmt"
	// "sync"
	// "time"

	"github.com/b0nbon1/stratal/internal/logger"
	"github.com/b0nbon1/stratal/internal/security"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

// TaskOutput stores the output from a task execution
type TaskOutput struct {
	TaskID   string
	TaskName string
	Output   string
	Error    error
}

// TaskLevel groups tasks that can be executed in parallel
type TaskLevel struct {
	Level int
	Tasks []db.Task
}

func ProcessJob(ctx context.Context, store *db.SQLStore, secretManager *security.SecretManager, jobRunID pgtype.UUID, job db.GetJobWithTasksRow, jobLogger *logger.JobRunLogger) error {
	fmt.Printf("Processing job run: %s for job: %s\n", jobRunID.String(), job.ID.String())

	if jobLogger != nil {
		jobLogger.Info(fmt.Sprintf("Processing job run: %s ", jobRunID.String()))
	}

	// Update job run status to running
	err := store.UpdateJobRunStatus(ctx, db.UpdateJobRunStatusParams{
		ID:     jobRunID,
		Status: utils.ParseText("running"),
	})
	if err != nil {
		if jobLogger != nil {
			jobLogger.Error(fmt.Sprintf("Failed to update job run status to running: %v", err))
		}
		return fmt.Errorf("failed to update job run status to running: %w", err)
	}

	// Parse tasks from JSON
	var tasks []db.Task
	if err := json.Unmarshal(job.Tasks, &tasks); err != nil {
		if jobLogger != nil {
			jobLogger.Error(fmt.Sprintf("Failed to unmarshal tasks: %v", err))
		}
		return fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks to execute")
		if jobLogger != nil {
			jobLogger.Info("No tasks to execute")
		}
		return completeJobRun(ctx, store, jobRunID, "completed", jobLogger)
	}

	// Sort tasks based on dependencies
	sortedTasks, err := utils.TopoSort(tasks)
	if err != nil {
		if jobLogger != nil {
			jobLogger.Error(fmt.Sprintf("Failed to sort tasks: %v", err))
		}
		return fmt.Errorf("failed to sort tasks: %w", err)
	}

	// Execute tasks in dependency order
	taskOutputs := make(map[string]string)
	taskNameToID := make(map[string]string)

	// Build task name to ID mapping
	for _, task := range sortedTasks {
		taskNameToID[task.Name] = task.ID.String()
	}

	// For now, use a dummy user ID for secret resolution
	userID := pgtype.UUID{}
	userID.Scan("00000000-0000-0000-0000-000000000001")

	for _, task := range sortedTasks {
		// Check if job run has been paused before executing next task
		currentJobRun, checkErr := store.GetJobRun(ctx, jobRunID)
		if checkErr != nil {
			if jobLogger != nil {
				jobLogger.Error(fmt.Sprintf("Failed to check job run status: %v", checkErr))
			}
			return fmt.Errorf("failed to check job run status: %w", checkErr)
		}

		if currentJobRun.Status.String == "paused" {
			if jobLogger != nil {
				jobLogger.Info("Job run has been paused, stopping execution")
			}
			return nil // Exit gracefully without error
		}

		fmt.Printf("Executing task: %s (type: %s)\n", task.Name, task.Type)
		if jobLogger != nil {
			jobLogger.Info(fmt.Sprintf("Executing task: %s (type: %s)", task.Name, task.Type))
		}

		var output string
		var err error

		if secretManager != nil {
			output, err = ExecuteTaskWithSecrets(ctx, task, store, secretManager, userID, taskOutputs, jobRunID, jobLogger)
		} else {
			output, err = ExecuteTaskWithOutputs(ctx, task, taskOutputs, taskNameToID, jobRunID, store, jobLogger)
		}

		if err != nil {
			fmt.Printf("Task %s failed: %v\n", task.Name, err)
			if jobLogger != nil {
				jobLogger.Error(fmt.Sprintf("Task %s failed: %v", task.Name, err))
			}
			// Mark job run as failed
			failErr := store.UpdateJobRunError(ctx, db.UpdateJobRunErrorParams{
				ID:           jobRunID,
				ErrorMessage: utils.ParseText(fmt.Sprintf("Task %s failed: %v", task.Name, err)),
			})
			if failErr != nil {
				fmt.Printf("Failed to update job run error: %v\n", failErr)
				if jobLogger != nil {
					jobLogger.Error(fmt.Sprintf("Failed to update job run error: %v", failErr))
				}
			}
			return err
		}

		// Store task output for future tasks
		taskOutputs[task.Name] = output
		fmt.Printf("Task %s completed successfully\n", task.Name)
		if jobLogger != nil {
			jobLogger.Info(fmt.Sprintf("Task %s completed successfully", task.Name))
		}
	}

	fmt.Println("All tasks completed successfully")
	if jobLogger != nil {
		jobLogger.Info("All tasks completed successfully")
	}
	return completeJobRun(ctx, store, jobRunID, "completed", jobLogger)
}

func completeJobRun(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID, status string, jobLogger *logger.JobRunLogger) error {
	err := store.UpdateJobRunStatus(ctx, db.UpdateJobRunStatusParams{
		ID:     jobRunID,
		Status: utils.ParseText(status),
	})
	if err != nil {
		if jobLogger != nil {
			jobLogger.Error(fmt.Sprintf("Failed to update job run status to %s: %v", status, err))
		}
		return fmt.Errorf("failed to update job run status to %s: %w", status, err)
	}
	fmt.Printf("Job run %s marked as %s\n", jobRunID.String(), status)
	if jobLogger != nil {
		jobLogger.Info(fmt.Sprintf("Job run %s marked as %s", jobRunID.String(), status))
	}
	return nil
}

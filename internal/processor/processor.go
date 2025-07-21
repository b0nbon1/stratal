package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

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

func ProcessJob(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID, job db.GetJobWithTasksRow) error {
	return ProcessJobWithSecrets(ctx, store, nil, jobRunID, job)
}

func ProcessJobWithSecrets(ctx context.Context, store *db.SQLStore, secretManager *security.SecretManager, jobRunID pgtype.UUID, job db.GetJobWithTasksRow) error {
	fmt.Printf("Processing job run: %s for job: %s\n", jobRunID.String(), job.ID.String())

	// Update job run status to running
	err := store.UpdateJobRunStatus(ctx, db.UpdateJobRunStatusParams{
		ID:     jobRunID,
		Status: utils.ParseText("running"),
	})
	if err != nil {
		return fmt.Errorf("failed to update job run status to running: %w", err)
	}

	// Parse tasks from JSON
	var tasks []db.Task
	if err := json.Unmarshal(job.Tasks, &tasks); err != nil {
		return fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks to execute")
		return completeJobRun(ctx, store, jobRunID, "completed")
	}

	// Sort tasks based on dependencies
	sortedTasks, err := utils.TopoSort(tasks)
	if err != nil {
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
		fmt.Printf("Executing task: %s (type: %s)\n", task.Name, task.Type)

		var output string
		var err error

		if secretManager != nil {
			output, err = ExecuteTaskWithSecrets(ctx, task, store, secretManager, userID, taskOutputs, jobRunID)
		} else {
			output, err = ExecuteTaskWithOutputs(ctx, task, taskOutputs, taskNameToID, jobRunID, store)
		}

		if err != nil {
			fmt.Printf("Task %s failed: %v\n", task.Name, err)
			// Mark job run as failed
			failErr := store.UpdateJobRunError(ctx, db.UpdateJobRunErrorParams{
				ID:           jobRunID,
				ErrorMessage: utils.ParseText(fmt.Sprintf("Task %s failed: %v", task.Name, err)),
			})
			if failErr != nil {
				fmt.Printf("Failed to update job run error: %v\n", failErr)
			}
			return err
		}

		// Find the task run ID using job_run_id and task_id
		taskRun, err := store.GetTaskRunByJobRunAndTaskID(ctx, db.GetTaskRunByJobRunAndTaskIDParams{
			JobRunID: jobRunID,
			TaskID:   task.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to find task run for job_run_id %s and task_id %s: %w", jobRunID.String(), task.ID.String(), err)
		}

		err = store.UpdateTaskRunOutput(ctx, db.UpdateTaskRunOutputParams{
			ID:     taskRun.ID,
			Output: utils.ParseText(output),
		})
		if err != nil {
			return fmt.Errorf("failed to update task run output: %w", err)
		}

		// Store task output for subsequent tasks
		taskOutputs[task.Name] = strings.TrimSpace(output)
		fmt.Printf("Task %s completed with output: %s\n", task.Name, taskOutputs[task.Name])
	}

	// Mark job run as completed
	return completeJobRun(ctx, store, jobRunID, "completed")
}

func completeJobRun(ctx context.Context, store *db.SQLStore, jobRunID pgtype.UUID, status string) error {
	err := store.UpdateJobRunStatus(ctx, db.UpdateJobRunStatusParams{
		ID:     jobRunID,
		Status: utils.ParseText(status),
	})
	if err != nil {
		return fmt.Errorf("failed to update job run status to %s: %w", status, err)
	}
	fmt.Printf("Job run %s marked as %s\n", jobRunID.String(), status)
	return nil
}

// groupTasksByDependencyLevel groups tasks into levels based on dependencies
func groupTasksByDependencyLevel(tasks []db.Task) ([]TaskLevel, error) {
	// Create maps for quick task lookup
	taskMap := make(map[string]db.Task)
	taskNameToID := make(map[string]string)

	for _, task := range tasks {
		taskMap[task.ID.String()] = task
		taskNameToID[task.Name] = task.ID.String()
	}

	// Normalize dependencies: convert task names to IDs
	for i, task := range tasks {
		if task.Config.DependsOn != nil {
			normalizedDeps := make([]string, 0, len(task.Config.DependsOn))
			for _, dep := range task.Config.DependsOn {
				// Check if it's a task ID (UUID format) or task name
				if _, exists := taskMap[dep]; exists {
					// It's already a task ID
					normalizedDeps = append(normalizedDeps, dep)
				} else if taskID, exists := taskNameToID[dep]; exists {
					// It's a task name, convert to ID
					normalizedDeps = append(normalizedDeps, taskID)
				} else {
					return nil, fmt.Errorf("task %s has unknown dependency: %s", task.Name, dep)
				}
			}
			tasks[i].Config.DependsOn = normalizedDeps
		}
	}

	// Calculate dependency levels
	levels := make(map[string]int)
	visited := make(map[string]bool)

	var calculateLevel func(taskID string) (int, error)
	calculateLevel = func(taskID string) (int, error) {
		if level, exists := levels[taskID]; exists {
			return level, nil
		}

		if visited[taskID] {
			return 0, fmt.Errorf("circular dependency detected involving task %s", taskID)
		}
		visited[taskID] = true

		task, exists := taskMap[taskID]
		if !exists {
			return 0, fmt.Errorf("task %s not found", taskID)
		}

		maxLevel := 0
		if task.Config.DependsOn != nil {
			for _, depID := range task.Config.DependsOn {
				depLevel, err := calculateLevel(depID)
				if err != nil {
					return 0, err
				}
				if depLevel >= maxLevel {
					maxLevel = depLevel + 1
				}
			}
		}

		levels[taskID] = maxLevel
		visited[taskID] = false
		return maxLevel, nil
	}

	// Calculate level for each task
	for _, task := range tasks {
		if _, err := calculateLevel(task.ID.String()); err != nil {
			return nil, err
		}
	}

	// Group tasks by level
	levelMap := make(map[int][]db.Task)
	for _, task := range tasks {
		level := levels[task.ID.String()]
		levelMap[level] = append(levelMap[level], task)
	}

	// Convert to sorted slice
	var taskLevels []TaskLevel
	maxLevel := 0
	for level := range levelMap {
		if level > maxLevel {
			maxLevel = level
		}
	}

	for i := 0; i <= maxLevel; i++ {
		if tasks, exists := levelMap[i]; exists {
			taskLevels = append(taskLevels, TaskLevel{
				Level: i,
				Tasks: tasks,
			})
		}
	}

	return taskLevels, nil
}

// executeTasksInParallel executes multiple tasks concurrently
func executeTasksInParallel(
	ctx context.Context,
	store *db.SQLStore,
	tasks []db.Task,
	taskRunMap map[string]db.ListTaskRunsRow,
	previousOutputs map[string]string,
	taskNameToID map[string]string,
) ([]TaskOutput, error) {
	var wg sync.WaitGroup
	results := make([]TaskOutput, len(tasks))

	// Use a semaphore to limit concurrent executions
	maxConcurrent := 5
	if len(tasks) < maxConcurrent {
		maxConcurrent = len(tasks)
	}
	semaphore := make(chan struct{}, maxConcurrent)

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t db.Task) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := TaskOutput{
				TaskID:   t.ID.String(),
				TaskName: t.Name,
			}

			// Get task run
			taskRun, exists := taskRunMap[t.ID.String()]
			if !exists {
				result.Error = fmt.Errorf("task run not found for task %s", t.ID.String())
				results[index] = result
				return
			}

			// Update task run status to running
			startTime := pgtype.Timestamp{Time: time.Now(), Valid: true}
			err := store.UpdateTaskRun(ctx, db.UpdateTaskRunParams{
				ID:        taskRun.ID,
				Status:    pgtype.Text{String: "running", Valid: true},
				StartedAt: startTime,
			})
			if err != nil {
				fmt.Printf("Failed to update task run status to running: %v\n", err)
			}

			fmt.Printf("Executing task: %s (type: %s)\n", t.Name, t.Type)

			// Interpolate parameters with outputs from previous tasks
			interpolatedTask := t
			if t.Type == "builtin" && t.Config.Parameters != nil {
				interpolatedTask.Config.Parameters = interpolateParameters(t.Config.Parameters, previousOutputs, taskNameToID)
			}

			// Execute the task
			output, taskErr := ExecuteTaskWithOutputs(ctx, interpolatedTask, previousOutputs, taskNameToID, taskRun.JobRunID, store)
			result.Output = output
			result.Error = taskErr

			// Update task run with results
			finishTime := pgtype.Timestamp{Time: time.Now(), Valid: true}
			var status string
			var errorMsg pgtype.Text
			var exitCode pgtype.Int4

			if taskErr != nil {
				status = "failed"
				errorMsg = pgtype.Text{String: taskErr.Error(), Valid: true}
				exitCode = pgtype.Int4{Int32: 1, Valid: true}
				fmt.Printf("Task %s failed: %v\n", t.Name, taskErr)
			} else {
				status = "completed"
				exitCode = pgtype.Int4{Int32: 0, Valid: true}
				fmt.Printf("Task %s completed successfully\n", t.Name)
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

			results[index] = result
		}(i, task)
	}

	wg.Wait()
	return results, nil
}

// interpolateParameters replaces ${task_name.output} references with actual outputs
func interpolateParameters(params map[string]string, outputs map[string]string, taskNameToID map[string]string) map[string]string {
	interpolated := make(map[string]string)

	for key, value := range params {
		// Look for ${task_name.output} pattern
		interpolated[key] = interpolateValue(value, outputs, taskNameToID)
	}

	return interpolated
}

// interpolateValue replaces output references in a single value
func interpolateValue(value string, outputs map[string]string, taskNameToID map[string]string) string {
	// Simple implementation - can be enhanced with regex for more complex cases
	for taskName, taskID := range taskNameToID {
		placeholder := fmt.Sprintf("${%s.output}", taskName)
		if output, exists := outputs[taskID]; exists {
			value = replaceAll(value, placeholder, output)
		}
	}
	return value
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	// Simple implementation - in production, use strings.ReplaceAll
	result := s
	for {
		index := findSubstring(result, old)
		if index == -1 {
			break
		}
		result = result[:index] + new + result[index+len(old):]
	}
	return result
}

// findSubstring finds the index of the first occurrence of substr in s
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// areDependenciesMet checks if all dependencies of a task are completed
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

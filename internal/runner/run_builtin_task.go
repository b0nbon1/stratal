package runner

import (
	"context"
	"fmt"
	"strings"

	"github.com/b0nbon1/stratal/internal/runner/tasks"
)

// TaskFunc is the signature for builtin task functions
type TaskFunc func(ctx context.Context, params map[string]string) (string, error)

// taskRegistry holds all registered builtin tasks
var taskRegistry = map[string]TaskFunc{
	"send_email":   wrapLegacyTask(tasks.SendEmailTask),
	"http_request": tasks.HTTPRequestTask,
	// Add more builtin tasks here as they are implemented
	// "database_query": tasks.DatabaseQueryTask,
	// "file_operation": tasks.FileOperationTask,
	// "slack_notification": tasks.SlackNotificationTask,
}

// RunBuiltinTask executes a builtin task by name with given parameters
func RunBuiltinTask(ctx context.Context, name string, params map[string]string) (string, error) {
	// Normalize task name
	taskName := strings.ToLower(strings.TrimSpace(name))

	// Look up task in registry
	taskFunc, exists := taskRegistry[taskName]
	if !exists {
		return "", fmt.Errorf("unknown builtin task: %s", name)
	}

	// Log task execution
	fmt.Printf("Executing builtin task: %s with %d parameters\n", taskName, len(params))

	// Execute task with context
	output, err := taskFunc(ctx, params)

	if err != nil {
		return output, fmt.Errorf("task %s failed: %w", taskName, err)
	}

	return output, nil
}

// RegisterBuiltinTask allows registering new builtin tasks at runtime
func RegisterBuiltinTask(name string, fn TaskFunc) error {
	taskName := strings.ToLower(strings.TrimSpace(name))

	if _, exists := taskRegistry[taskName]; exists {
		return fmt.Errorf("task %s is already registered", taskName)
	}

	taskRegistry[taskName] = fn
	return nil
}

// GetAvailableTasks returns a list of all registered builtin tasks
func GetAvailableTasks() []string {
	tasks := make([]string, 0, len(taskRegistry))
	for name := range taskRegistry {
		tasks = append(tasks, name)
	}
	return tasks
}

// wrapLegacyTask wraps old-style task functions that don't accept context
func wrapLegacyTask(fn func(map[string]string) error) TaskFunc {
	return func(ctx context.Context, params map[string]string) (string, error) {
		// Check context cancellation before executing
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		// Execute the legacy task
		err := fn(params)
		if err != nil {
			return "", err
		}

		return "Task completed successfully", nil
	}
}

// Example of a simple builtin task that accepts context
func echoTask(ctx context.Context, params map[string]string) (string, error) {
	message, exists := params["message"]
	if !exists {
		return "", fmt.Errorf("missing required parameter: message")
	}

	return fmt.Sprintf("Echo: %s", message), nil
}

// Initialize with some basic tasks
func init() {
	// Register the echo task as an example
	taskRegistry["echo"] = echoTask
}

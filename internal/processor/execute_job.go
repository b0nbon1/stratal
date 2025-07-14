package processor

import (
	"context"
	"fmt"

	"github.com/b0nbon1/stratal/internal/runner"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

// ExecuteTaskWithOutputs executes a task with access to outputs from previous tasks
func ExecuteTaskWithOutputs(ctx context.Context, task db.Task, outputs map[string]string, taskNameToID map[string]string) (string, error) {
	switch task.Type {
	case "builtin":
		// Builtin tasks already have interpolated parameters from the processor
		output, err := runner.RunBuiltinTask(ctx, task.Name, task.Config.Parameters)
		return output, err
	case "custom":
		// For custom scripts, create a map of task names to outputs
		taskNameOutputs := make(map[string]string)
		for taskName, taskID := range taskNameToID {
			if output, exists := outputs[taskID]; exists {
				taskNameOutputs[taskName] = output
			}
		}
		output, err := runner.RunCustomScriptWithOutputs(ctx, task.Config.Script, taskNameOutputs)
		return output, err
	default:
		return "", fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

// ExecuteTask maintains backward compatibility
func ExecuteTask(ctx context.Context, task db.Task) (string, error) {
	return ExecuteTaskWithOutputs(ctx, task, nil, nil)
}

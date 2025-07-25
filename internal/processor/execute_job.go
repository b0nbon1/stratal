package processor

import (
	"context"
	"fmt"

	"github.com/b0nbon1/stratal/internal/logger"
	"github.com/b0nbon1/stratal/internal/runner"
	"github.com/b0nbon1/stratal/internal/security"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// ExecuteTaskWithOutputs executes a task with access to outputs from previous tasks
func ExecuteTaskWithOutputs(ctx context.Context, task db.Task, outputs map[string]string, taskNameToID map[string]string, jobRunID pgtype.UUID, store *db.SQLStore, jobLogger *logger.JobRunLogger) (string, error) {
	// Get task run ID for logging
	taskRun, err := store.GetTaskRunByJobRunAndTaskID(ctx, db.GetTaskRunByJobRunAndTaskIDParams{
		JobRunID: jobRunID,
		TaskID:   task.ID,
	})
	if err != nil {
		if jobLogger != nil {
			jobLogger.Error(fmt.Sprintf("Failed to find task run for task %s: %v", task.Name, err))
		}
		return "", fmt.Errorf("failed to find task run for task %s: %w", task.Name, err)
	}

	taskRunID := taskRun.ID.String()

	if jobLogger != nil {
		jobLogger.InfoWithTaskRun(taskRunID, fmt.Sprintf("Starting execution of task %s (type: %s)", task.Name, task.Type))
	}

	switch task.Type {
	case "builtin":
		// Builtin tasks already have interpolated parameters from the processor
		output, err := runner.RunBuiltinTask(ctx, task.Name, task.Config.Parameters, outputs)
		if err != nil && jobLogger != nil {
			jobLogger.ErrorWithTaskRun(taskRunID, fmt.Sprintf("Builtin task %s failed: %v", task.Name, err))
		} else if jobLogger != nil {
			jobLogger.InfoWithTaskRun(taskRunID, fmt.Sprintf("Builtin task %s completed successfully", task.Name))
		}
		return output, err
	case "custom":
		if task.Config.Script == nil {
			err := fmt.Errorf("custom task %s has no script configuration", task.Name)
			if jobLogger != nil {
				jobLogger.ErrorWithTaskRun(taskRunID, err.Error())
			}
			return "", err
		}
		output, err := runner.RunCustomScriptWithOutputs(ctx, task.Config.Script, outputs)
		if err != nil && jobLogger != nil {
			jobLogger.ErrorWithTaskRun(taskRunID, fmt.Sprintf("Custom script task %s failed: %v", task.Name, err))
		} else if jobLogger != nil {
			jobLogger.InfoWithTaskRun(taskRunID, fmt.Sprintf("Custom script task %s completed successfully", task.Name))
		}
		return output, err
	default:
		err := fmt.Errorf("unsupported task type: %s", task.Type)
		if jobLogger != nil {
			jobLogger.ErrorWithTaskRun(taskRunID, err.Error())
		}
		return "", err
	}
}

// ExecuteTaskWithSecrets executes a task with access to outputs, parameters, and secrets
func ExecuteTaskWithSecrets(ctx context.Context, task db.Task, store *db.SQLStore, secretManager *security.SecretManager, userID pgtype.UUID, outputs map[string]string, jobRunID pgtype.UUID, jobLogger *logger.JobRunLogger) (string, error) {
	// Get task run ID for logging
	taskRun, err := store.GetTaskRunByJobRunAndTaskID(ctx, db.GetTaskRunByJobRunAndTaskIDParams{
		JobRunID: jobRunID,
		TaskID:   task.ID,
	})
	if err != nil {
		if jobLogger != nil {
			jobLogger.Error(fmt.Sprintf("Failed to find task run for task %s: %v", task.Name, err))
		}
		return "", fmt.Errorf("failed to find task run for task %s: %w", task.Name, err)
	}

	taskRunID := taskRun.ID.String()

	if jobLogger != nil {
		jobLogger.InfoWithTaskRun(taskRunID, fmt.Sprintf("Starting execution of task %s with secrets (type: %s)", task.Name, task.Type))
	}

	// Use parameter resolver to get resolved parameters and secrets
	resolver := NewParameterResolver(store, secretManager)
	resolvedParams, secretEnvVars, err := resolver.ResolveParameters(ctx, task, userID, outputs)
	if err != nil {
		if jobLogger != nil {
			jobLogger.ErrorWithTaskRun(taskRunID, fmt.Sprintf("Failed to resolve parameters for task %s: %v", task.Name, err))
		}
		return "", fmt.Errorf("failed to resolve parameters for task %s: %w", task.Name, err)
	}

	switch task.Type {
	case "builtin":
		// For builtin tasks, merge resolved parameters and pass to runner
		allParams := make(map[string]string)
		for k, v := range resolvedParams {
			allParams[k] = v
		}
		// Add secrets to parameters for builtin tasks
		for k, v := range secretEnvVars {
			allParams[k] = v
		}
		output, err := runner.RunBuiltinTask(ctx, task.Name, allParams, outputs)
		if err != nil && jobLogger != nil {
			jobLogger.ErrorWithTaskRun(taskRunID, fmt.Sprintf("Builtin task %s failed: %v", task.Name, err))
		} else if jobLogger != nil {
			jobLogger.InfoWithTaskRun(taskRunID, fmt.Sprintf("Builtin task %s completed successfully", task.Name))
		}
		return output, err
	case "custom":
		if task.Config.Script == nil {
			err := fmt.Errorf("custom task %s has no script configuration", task.Name)
			if jobLogger != nil {
				jobLogger.ErrorWithTaskRun(taskRunID, err.Error())
			}
			return "", err
		}
		output, err := runner.RunCustomScriptWithSecrets(ctx, task.Config.Script, resolvedParams, secretEnvVars, outputs)
		if err != nil && jobLogger != nil {
			jobLogger.ErrorWithTaskRun(taskRunID, fmt.Sprintf("Custom script task %s failed: %v", task.Name, err))
		} else if jobLogger != nil {
			jobLogger.InfoWithTaskRun(taskRunID, fmt.Sprintf("Custom script task %s completed successfully", task.Name))
		}
		return output, err
	default:
		err := fmt.Errorf("unsupported task type: %s", task.Type)
		if jobLogger != nil {
			jobLogger.ErrorWithTaskRun(taskRunID, err.Error())
		}
		return "", err
	}
}

// ExecuteTask maintains backward compatibility
func ExecuteTask(ctx context.Context, task db.Task, store *db.SQLStore) (string, error) {
	// For backward compatibility, we need a dummy jobRunID
	// This should not be used in production - callers should use ExecuteTaskWithOutputs instead
	dummyJobRunID := pgtype.UUID{}
	return ExecuteTaskWithOutputs(ctx, task, nil, nil, dummyJobRunID, store, nil)
}

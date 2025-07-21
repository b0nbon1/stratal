package processor

import (
	"context"
	"fmt"

	"github.com/b0nbon1/stratal/internal/runner"
	"github.com/b0nbon1/stratal/internal/security"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// ExecuteTaskWithOutputs executes a task with access to outputs from previous tasks
func ExecuteTaskWithOutputs(ctx context.Context, task db.Task, outputs map[string]string, taskNameToID map[string]string, jobRunID pgtype.UUID, store *db.SQLStore) (string, error) {
	switch task.Type {
	case "builtin":
		// Builtin tasks already have interpolated parameters from the processor
		output, err := runner.RunBuiltinTask(ctx, task.Name, task.Config.Parameters, outputs)
		return output, err
	case "custom":
		if task.Config.Script == nil {
			return "", fmt.Errorf("custom task %s has no script configuration", task.Name)
		}
		return runner.RunCustomScriptWithOutputs(ctx, task.Config.Script, outputs)
	default:
		return "", fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

// ExecuteTaskWithSecrets executes a task with access to outputs, parameters, and secrets
func ExecuteTaskWithSecrets(
	ctx context.Context,
	task db.Task,
	store *db.SQLStore,
	secretManager *security.SecretManager,
	userID pgtype.UUID,
	taskOutputs map[string]string,
	jobRunID pgtype.UUID,
) (string, error) {
	// Use parameter resolver to get resolved parameters and secrets
	resolver := NewParameterResolver(store, secretManager)
	resolvedParams, secretEnvVars, err := resolver.ResolveParameters(ctx, task, userID, taskOutputs)
	if err != nil {
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
		return runner.RunBuiltinTask(ctx, task.Name, allParams, taskOutputs)
	case "custom":
		if task.Config.Script == nil {
			return "", fmt.Errorf("custom task %s has no script configuration", task.Name)
		}
		return runner.RunCustomScriptWithSecrets(ctx, task.Config.Script, resolvedParams, secretEnvVars, taskOutputs)
	default:
		return "", fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

// ExecuteTask maintains backward compatibility
func ExecuteTask(ctx context.Context, task db.Task, store *db.SQLStore) (string, error) {
	// For backward compatibility, we need a dummy jobRunID
	// This should not be used in production - callers should use ExecuteTaskWithOutputs instead
	dummyJobRunID := pgtype.UUID{}
	return ExecuteTaskWithOutputs(ctx, task, nil, nil, dummyJobRunID, store)
}

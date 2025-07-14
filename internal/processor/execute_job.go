package processor

import (
	"context"
	"fmt"

	"github.com/b0nbon1/stratal/internal/runner"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

func ExecuteTask(ctx context.Context, task db.Task) (string, error) {
	switch task.Type {
	case "builtin":
		output, err := runner.RunBuiltinTask(ctx, task.Name, task.Config.Parameters)
		return output, err
	case "custom":
		output, err := runner.RunCustomScript(ctx, task.Config.Script)
		return output, err
	default:
		return "", fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

package processor

import (
	"fmt"

	"github.com/b0nbon1/stratal/internal/runner"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

func ExecuteTask(task db.Task) error {
	switch task.Type {
	case "builtin":
		return runner.RunBuiltinTask(task.Name, task.Config.Parameters)
	case "custom":
		return runner.RunCustomScript(task.Config.Script)
	default:
		return fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

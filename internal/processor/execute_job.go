package processor

import (
	"fmt"

	"github.com/b0nbon1/stratal/db/dto"
	"github.com/b0nbon1/stratal/internal/builtin"
)

func ExecuteTask(task dto.TaskConfig) error {
	switch task.Type {
	case "builtin":
		return builtin.RunBuiltinTask(task.Name, task.Parameters)
	case "custom":
		return builtin.RunCustomScript(task.Script)
	default:
		return fmt.Errorf("unsupported task type: %s", task.Type)
	}
}

package runner

import (
	"fmt"

	"github.com/b0nbon1/stratal/internal/storage/db/dto"
)

func RunCustomScript(script *dto.ScriptConfig) error {
	fmt.Println(script)
	return nil
}
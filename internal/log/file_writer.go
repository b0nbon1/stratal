package logger

import (
	"encoding/json"
	"fmt"
	"os"

	db "github.com/b0nbon1/stratal/db/sqlc"
)

type FileWriter struct {
    Path string
}

func (fw *FileWriter) WriteLogs(logs []db.CreateJobLogParams) error {
    file, err := os.OpenFile(fw.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("open file: %w", err)
    }
    defer file.Close()

    for _, logEntry := range logs {
        jsonLog, _ := json.Marshal(logEntry)
        file.Write(jsonLog)
        file.Write([]byte("\n"))
    }

    return nil
}

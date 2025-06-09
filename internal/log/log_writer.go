package logger

import (
	"time"

	db "github.com/b0nbon1/stratal/db/sqlc"
)

type TaskLog struct {
    TaskID    string
    LogLevel  string
    Message   string
    CreatedAt time.Time
}

type LogWriter interface {
    WriteLogs([]db.CreateJobLogParams) error
}

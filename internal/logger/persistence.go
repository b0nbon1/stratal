package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// writeSystemLogToDatabase writes a system log entry to the database
func (l *Logger) writeSystemLogToDatabase(entry LogEntry) {
	if l.store == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var metadataJSON []byte
	if entry.Metadata != nil {
		metadataJSON, _ = json.Marshal(entry.Metadata)
	}

	err := l.store.CreateSystemLog(ctx, db.CreateSystemLogParams{
		Level:    string(entry.Level),
		Stream:   entry.Stream,
		Message:  entry.Message,
		Metadata: metadataJSON,
	})

	if err != nil {
		log.Printf("Failed to write system log to database: %v", err)
	}
}

// writeToFile writes a log entry to the job run's log file
func (jrl *JobRunLogger) writeToFile(entry LogEntry) {
	jrl.fileMux.Lock()
	defer jrl.fileMux.Unlock()

	if jrl.logFile == nil {
		return
	}

	logLine := fmt.Sprintf("[%s] [%s] [%s] [%s] %s\n",
		entry.Timestamp.Format("2006-01-02 15:04:05.000"),
		entry.Type,
		entry.Level,
		entry.Stream,
		entry.Message)

	if _, err := jrl.logFile.WriteString(logLine); err != nil {
		log.Printf("Failed to write to log file: %v", err)
	}
}

// writeToDatabase writes a log entry to the database
func (jrl *JobRunLogger) writeToDatabase(entry LogEntry) {
	jrl.dbMux.Lock()
	defer jrl.dbMux.Unlock()

	if jrl.store == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var metadataJSON []byte
	if entry.Metadata != nil {
		metadataJSON, _ = json.Marshal(entry.Metadata)
	}

	var jobRunUUID, taskRunUUID pgtype.UUID

	if entry.JobRunID != "" {
		if err := jobRunUUID.Scan(entry.JobRunID); err != nil {
			log.Printf("Failed to parse job run ID %s: %v", entry.JobRunID, err)
			return
		}
	}

	if entry.TaskRunID != "" {
		if err := taskRunUUID.Scan(entry.TaskRunID); err != nil {
			log.Printf("Failed to parse task run ID %s: %v", entry.TaskRunID, err)
			return
		}
	}

	switch entry.Type {
	case JobLogType:
		err := jrl.store.CreateJobLog(ctx, db.CreateJobLogParams{
			JobRunID: jobRunUUID,
			Level:    string(entry.Level),
			Stream:   entry.Stream,
			Message:  entry.Message,
			Metadata: metadataJSON,
		})
		if err != nil {
			log.Printf("Failed to write job log to database: %v", err)
		}

	case TaskLogType:
		err := jrl.store.CreateTaskLog(ctx, db.CreateTaskLogParams{
			JobRunID:  jobRunUUID,
			TaskRunID: taskRunUUID,
			Level:     string(entry.Level),
			Stream:    entry.Stream,
			Message:   entry.Message,
			Metadata:  metadataJSON,
		})
		if err != nil {
			log.Printf("Failed to write task log to database: %v", err)
		}

	default:
		err := jrl.store.CreateLog(ctx, db.CreateLogParams{
			Type:      string(entry.Type),
			JobRunID:  jobRunUUID,
			TaskRunID: taskRunUUID,
			Timestamp: pgtype.Timestamp{Time: entry.Timestamp, Valid: true},
			Level:     string(entry.Level),
			Stream:    entry.Stream,
			Message:   entry.Message,
			Metadata:  metadataJSON,
		})
		if err != nil {
			log.Printf("Failed to write log to database: %v", err)
		}
	}
}

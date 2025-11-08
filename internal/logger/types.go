package logger

import (
	"os"
	"sync"
	"time"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

type LogType string

const (
	SystemLogType LogType = "system"
	JobLogType    LogType = "job"
	TaskLogType   LogType = "task"
)

type LogLevel string

const (
	InfoLevel  LogLevel = "info"
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
	DebugLevel LogLevel = "debug"
)

type LogEntry struct {
	Type      LogType
	JobRunID  string
	TaskRunID string
	Level     LogLevel
	Message   string
	Timestamp time.Time
	Stream    string // "stdout", "stderr", or "system"
	Metadata  map[string]interface{}
}

type JobRunLogger struct {
	jobRunID string
	store    *db.SQLStore
	logFile  *os.File
	fileMux  sync.Mutex
	dbMux    sync.Mutex
	logger   *Logger // Reference to the main Logger
}

// Logger manages job run loggers
type Logger struct {
	store       *db.SQLStore
	loggers     map[string]*JobRunLogger
	loggersMux  sync.RWMutex
	baseLogPath string
	streamer    *LogStreamer
}

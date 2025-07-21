package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// LogType represents the type of log entry
type LogType string

const (
	SystemLogType LogType = "system" // System-level logs (server startup, shutdown, etc.)
	JobLogType    LogType = "job"    // Job-level logs (job start, completion, errors)
	TaskLogType   LogType = "task"   // Task-level logs (individual task execution)
)

// LogLevel represents the logging level
type LogLevel string

const (
	InfoLevel  LogLevel = "info"
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
	DebugLevel LogLevel = "debug"
)

// LogEntry represents a log entry with enhanced fields for the new logs table
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

// JobRunLogger handles logging for a specific job run
type JobRunLogger struct {
	jobRunID string
	store    *db.SQLStore
	logFile  *os.File
	fileMux  sync.Mutex
	dbMux    sync.Mutex
	logger   *Logger // Add a reference to the main Logger
}

// Logger manages job run loggers
type Logger struct {
	store       *db.SQLStore
	loggers     map[string]*JobRunLogger
	loggersMux  sync.RWMutex
	baseLogPath string
	streamer    *LogStreamer // Add streaming support
}

// NewLogger creates a new logger instance
func NewLogger(store *db.SQLStore, baseLogPath string) *Logger {
	if baseLogPath == "" {
		baseLogPath = "internal/storage/files/logs"
	}

	// Ensure log directory exists
	if err := os.MkdirAll(baseLogPath, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
	}

	return &Logger{
		store:       store,
		loggers:     make(map[string]*JobRunLogger),
		baseLogPath: baseLogPath,
		streamer:    NewLogStreamer(), // Initialize streaming
	}
}

// GetStreamer returns the log streamer for setting up HTTP handlers
func (l *Logger) GetStreamer() *LogStreamer {
	return l.streamer
}

// LogSystem logs a system-level message
func (l *Logger) LogSystem(level LogLevel, message string, metadata map[string]interface{}) {
	entry := LogEntry{
		Type:      SystemLogType,
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Stream:    "system",
		Metadata:  metadata,
	}

	// Write to database
	l.writeSystemLogToDatabase(entry)

	// Broadcast to streaming clients
	if l.streamer != nil {
		streamMsg := LogMessage{
			Type:      string(SystemLogType),
			JobRunID:  "",
			Timestamp: entry.Timestamp,
			Level:     level,
			Stream:    entry.Stream,
			Message:   message,
			Metadata:  metadata,
		}
		l.streamer.BroadcastLog(streamMsg)
	}

	// Console output
	log.Printf("[SYSTEM] %s: %s", level, message)
}

// writeSystemLogToDatabase writes a system log to the database
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

// GetJobRunLogger gets or creates a logger for a specific job run
func (l *Logger) GetJobRunLogger(jobRunID string) (*JobRunLogger, error) {
	l.loggersMux.Lock()
	defer l.loggersMux.Unlock()

	if logger, exists := l.loggers[jobRunID]; exists {
		return logger, nil
	}

	// Create new job run logger
	logger, err := l.createJobRunLogger(jobRunID)
	if err != nil {
		return nil, fmt.Errorf("failed to create job run logger: %w", err)
	}

	l.loggers[jobRunID] = logger
	return logger, nil
}

// createJobRunLogger creates a new JobRunLogger
func (l *Logger) createJobRunLogger(jobRunID string) (*JobRunLogger, error) {
	// Create log file path: {job_run_id}-{date}.txt
	now := time.Now()
	filename := fmt.Sprintf("%s-%s.txt", jobRunID, now.Format("2006-01-02"))
	logFilePath := filepath.Join(l.baseLogPath, filename)

	// Create or open log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %s: %w", logFilePath, err)
	}

	return &JobRunLogger{
		jobRunID: jobRunID,
		store:    l.store,
		logFile:  logFile,
		logger:   l, // Set the reference to the main Logger
	}, nil
}

// CloseJobRunLogger closes and removes a job run logger
func (l *Logger) CloseJobRunLogger(jobRunID string) {
	l.loggersMux.Lock()
	defer l.loggersMux.Unlock()

	if logger, exists := l.loggers[jobRunID]; exists {
		logger.Close()
		delete(l.loggers, jobRunID)
	}
}

// Close closes all job run loggers
func (l *Logger) Close() {
	l.loggersMux.Lock()
	defer l.loggersMux.Unlock()

	for jobRunID, logger := range l.loggers {
		logger.Close()
		delete(l.loggers, jobRunID)
	}
}

// Job-level logging methods
func (jrl *JobRunLogger) LogJob(level LogLevel, message string, metadata map[string]interface{}) {
	jrl.logEntry(JobLogType, "", level, message, "system", metadata)
}

// Task-level logging methods
func (jrl *JobRunLogger) LogTask(taskRunID string, level LogLevel, message string, stream string, metadata map[string]interface{}) {
	jrl.logEntry(TaskLogType, taskRunID, level, message, stream, metadata)
}

// Convenience methods for different log levels
func (jrl *JobRunLogger) Info(message string) {
	jrl.LogJob(InfoLevel, message, nil)
}

func (jrl *JobRunLogger) Error(message string) {
	jrl.LogJob(ErrorLevel, message, nil)
}

func (jrl *JobRunLogger) Warn(message string) {
	jrl.LogJob(WarnLevel, message, nil)
}

func (jrl *JobRunLogger) Debug(message string) {
	jrl.LogJob(DebugLevel, message, nil)
}

// Task-specific logging methods
func (jrl *JobRunLogger) InfoWithTaskRun(taskRunID string, message string) {
	jrl.LogTask(taskRunID, InfoLevel, message, "stdout", nil)
}

func (jrl *JobRunLogger) ErrorWithTaskRun(taskRunID string, message string) {
	jrl.LogTask(taskRunID, ErrorLevel, message, "stderr", nil)
}

// logEntry is the core logging method
func (jrl *JobRunLogger) logEntry(logType LogType, taskRunID string, level LogLevel, message string, stream string, metadata map[string]interface{}) {
	entry := LogEntry{
		Type:      logType,
		JobRunID:  jrl.jobRunID,
		TaskRunID: taskRunID,
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Stream:    stream,
		Metadata:  metadata,
	}

	// Write to file
	jrl.writeToFile(entry)

	// Write to database
	jrl.writeToDatabase(entry)

	// Broadcast to streaming clients if streamer is available
	if jrl.logger != nil && jrl.logger.streamer != nil {
		streamMsg := LogMessage{
			Type:      string(logType),
			JobRunID:  jrl.jobRunID,
			TaskRunID: taskRunID,
			Timestamp: entry.Timestamp,
			Level:     level,
			Stream:    stream,
			Message:   message,
			Metadata:  metadata,
		}
		jrl.logger.streamer.BroadcastLog(streamMsg)
	}
}

// writeToFile writes log entry to file
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

// writeToDatabase writes log entry to database
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

	// Parse job run ID
	if entry.JobRunID != "" {
		if err := jobRunUUID.Scan(entry.JobRunID); err != nil {
			log.Printf("Failed to parse job run ID %s: %v", entry.JobRunID, err)
			return
		}
	}

	// Parse task run ID if provided
	if entry.TaskRunID != "" {
		if err := taskRunUUID.Scan(entry.TaskRunID); err != nil {
			log.Printf("Failed to parse task run ID %s: %v", entry.TaskRunID, err)
			return
		}
	}

	// Choose the appropriate database method based on log type
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
		// Generic log creation
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

// GetWriter returns an io.Writer for the specified stream
func (jrl *JobRunLogger) GetWriter(stream string) io.Writer {
	return &logWriter{
		logger: jrl,
		stream: stream,
	}
}

// GetWriterForTaskRun returns an io.Writer for a specific task run and stream
func (jrl *JobRunLogger) GetWriterForTaskRun(taskRunID, stream string) io.Writer {
	return &taskLogWriter{
		logger:    jrl,
		stream:    stream,
		taskRunID: taskRunID,
	}
}

// Close closes the job run logger
func (jrl *JobRunLogger) Close() {
	jrl.fileMux.Lock()
	defer jrl.fileMux.Unlock()

	if jrl.logFile != nil {
		jrl.logFile.Close()
		jrl.logFile = nil
	}
}

// logWriter implements io.Writer for general logging
type logWriter struct {
	logger *JobRunLogger
	stream string
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	if w.stream == "stderr" {
		w.logger.LogJob(ErrorLevel, message, nil)
	} else {
		w.logger.LogJob(InfoLevel, message, nil)
	}
	return len(p), nil
}

// taskLogWriter implements io.Writer for task-specific logging
type taskLogWriter struct {
	logger    *JobRunLogger
	stream    string
	taskRunID string
}

func (w *taskLogWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	if w.stream == "stderr" {
		w.logger.LogTask(w.taskRunID, ErrorLevel, message, w.stream, nil)
	} else {
		w.logger.LogTask(w.taskRunID, InfoLevel, message, w.stream, nil)
	}
	return len(p), nil
}

// parseUUID parses a string UUID into pgtype.UUID
func parseUUID(uuidStr string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(uuidStr)
	return uuid, err
}

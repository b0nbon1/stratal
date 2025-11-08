package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
)

// NewLogger creates a new Logger instance
func NewLogger(store *db.SQLStore, baseLogPath string) *Logger {
	if baseLogPath == "" {
		baseLogPath = "internal/storage/files/logs"
	}

	if err := os.MkdirAll(baseLogPath, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
	}

	return &Logger{
		store:       store,
		loggers:     make(map[string]*JobRunLogger),
		baseLogPath: baseLogPath,
		streamer:    NewLogStreamer(),
	}
}

// GetStreamer returns the log streamer for WebSocket/SSE connections
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

	l.writeSystemLogToDatabase(entry)
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
	log.Printf("[SYSTEM] %s: %s", level, message)
}

// GetJobRunLogger retrieves or creates a logger for a specific job run
func (l *Logger) GetJobRunLogger(jobRunID string) (*JobRunLogger, error) {
	l.loggersMux.Lock()
	defer l.loggersMux.Unlock()

	if logger, exists := l.loggers[jobRunID]; exists {
		return logger, nil
	}

	logger, err := l.createJobRunLogger(jobRunID)
	if err != nil {
		return nil, fmt.Errorf("failed to create job run logger: %w", err)
	}

	l.loggers[jobRunID] = logger
	return logger, nil
}

// createJobRunLogger creates a new JobRunLogger with its associated log file
func (l *Logger) createJobRunLogger(jobRunID string) (*JobRunLogger, error) {
	now := time.Now()
	filename := fmt.Sprintf("%s-%s.txt", jobRunID, now.Format("2006-01-02"))
	logFilePath := filepath.Join(l.baseLogPath, filename)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %s: %w", logFilePath, err)
	}

	return &JobRunLogger{
		jobRunID: jobRunID,
		store:    l.store,
		logFile:  logFile,
		logger:   l,
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

// LogJob logs a job-level message
func (jrl *JobRunLogger) LogJob(level LogLevel, message string, metadata map[string]interface{}) {
	jrl.logEntry(JobLogType, "", level, message, "system", metadata)
}

// LogTask logs a task-level message
func (jrl *JobRunLogger) LogTask(taskRunID string, level LogLevel, message string, stream string, metadata map[string]interface{}) {
	jrl.logEntry(TaskLogType, taskRunID, level, message, stream, metadata)
}

// Info logs an info-level job message
func (jrl *JobRunLogger) Info(message string) {
	jrl.LogJob(InfoLevel, message, nil)
}

// Error logs an error-level job message
func (jrl *JobRunLogger) Error(message string) {
	jrl.LogJob(ErrorLevel, message, nil)
}

// Warn logs a warning-level job message
func (jrl *JobRunLogger) Warn(message string) {
	jrl.LogJob(WarnLevel, message, nil)
}

// Debug logs a debug-level job message
func (jrl *JobRunLogger) Debug(message string) {
	jrl.LogJob(DebugLevel, message, nil)
}

// InfoWithTaskRun logs an info-level task message
func (jrl *JobRunLogger) InfoWithTaskRun(taskRunID string, message string) {
	jrl.LogTask(taskRunID, InfoLevel, message, "stdout", nil)
}

// ErrorWithTaskRun logs an error-level task message
func (jrl *JobRunLogger) ErrorWithTaskRun(taskRunID string, message string) {
	jrl.LogTask(taskRunID, ErrorLevel, message, "stderr", nil)
}

// logEntry is the core logging method that handles file, database, and streaming
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

	jrl.writeToFile(entry)
	jrl.writeToDatabase(entry)

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

// Close closes the job run logger and its log file
func (jrl *JobRunLogger) Close() {
	jrl.fileMux.Lock()
	defer jrl.fileMux.Unlock()

	if jrl.logFile != nil {
		jrl.logFile.Close()
		jrl.logFile = nil
	}
}

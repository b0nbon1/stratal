# Logging System

The Stratal logging system provides comprehensive logging capabilities for job runs, storing logs both in the database and in files for easy access and debugging.

## Features

- **Dual Logging**: Logs are written to both the database (`task_logs` table) and log files
- **Per-Job Log Files**: Each job run gets its own log file with the format `{job_run_id}-{date}.txt`
- **Structured Logging**: Logs include timestamps, log levels, and stream information
- **Task-Level Logging**: Individual tasks within a job run can be logged separately
- **Automatic Cleanup**: Log files are managed and closed automatically when job runs complete

## Log File Location

Log files are stored in: `storage/files/logs/{job_run_id}-{date}.txt`

Example: `storage/files/logs/a1b2c3d4-e5f6-7890-abcd-ef1234567890-2024-01-15.txt`

## Log Levels

- `info`: General information about job/task execution
- `error`: Error messages and failures
- `warn`: Warning messages
- `debug`: Detailed debugging information

## Database Storage

Logs are stored in the `task_logs` table with the following structure:

```sql
CREATE TABLE task_logs (
    id BIGSERIAL PRIMARY KEY,
    task_run_id UUID REFERENCES task_runs (id) ON DELETE CASCADE,
    timestamp TIMESTAMP DEFAULT now (),
    stream TEXT CHECK (stream IN ('stdout', 'stderr')) NOT NULL,
    message TEXT NOT NULL
);
```

## Log Format

### File Log Format
```
[2024-01-15 10:30:45.123] [info] [stdout] Starting job run a1b2c3d4-e5f6-7890-abcd-ef1234567890
[2024-01-15 10:30:45.234] [info] [stdout] Processing job run: a1b2c3d4-e5f6-7890-abcd-ef1234567890 for job: b2c3d4e5-f6g7-8901-bcde-f23456789012
[2024-01-15 10:30:45.345] [info] [stdout] Executing task: info_task (type: builtin)
[2024-01-15 10:30:46.456] [info] [stdout] Builtin task info_task completed successfully
[2024-01-15 10:30:46.567] [info] [stdout] All tasks completed successfully
[2024-01-15 10:30:46.678] [info] [stdout] Job run a1b2c3d4-e5f6-7890-abcd-ef1234567890 completed successfully
```

## Usage

### Logger Initialization

The logger is automatically initialized when the worker starts:

```go
// Initialize the logger system
logSystem := logger.NewLogger(store, "storage/files/logs")
defer logSystem.Close()
```

### Job Run Logging

Each job run automatically gets its own logger instance:

```go
// Get or create a logger for this job run
jobLogger, err := logSystem.GetJobRunLogger(jobRunId)
if err != nil {
    fmt.Printf("Error creating logger for job_run %s: %v\n", jobRunId, err)
}

// Log job run start
jobLogger.Info(fmt.Sprintf("Starting job run %s", jobRunId))

// Log errors
jobLogger.Error(fmt.Sprintf("Job run failed: %v", err))
```

### Task-Level Logging

Individual tasks can log to both the file and database:

```go
// Log with task run ID for database storage
jobLogger.InfoWithTaskRun(taskRunID, "Starting execution of task")
jobLogger.ErrorWithTaskRun(taskRunID, "Task execution failed")
```

## API Endpoints

### List Task Logs

Get logs for a specific task run:
```sql
-- SQL Query
SELECT id, task_run_id, timestamp, stream, message
FROM task_logs
WHERE task_run_id = $1
ORDER BY timestamp ASC;
```

### List Job Run Logs

Get all logs for a job run:
```sql
-- SQL Query
SELECT tl.id, tl.task_run_id, tl.timestamp, tl.stream, tl.message
FROM task_logs tl
JOIN task_runs tr ON tl.task_run_id = tr.id
WHERE tr.job_run_id = $1
ORDER BY tl.timestamp ASC;
```

## Integration

The logging system is integrated into:

1. **Worker**: Initializes the logger system and manages job run loggers
2. **Processor**: Logs job processing steps and task execution
3. **Task Execution**: Logs individual task starts, completions, and failures
4. **Error Handling**: Logs all errors and failures with context

## Benefits

- **Debugging**: Easy access to detailed logs for troubleshooting failed jobs
- **Monitoring**: Real-time visibility into job execution progress
- **Audit Trail**: Complete history of job execution stored in the database
- **File Access**: Log files can be downloaded or viewed directly from the file system
- **Performance**: Asynchronous logging doesn't block job execution

## Configuration

The log directory can be configured when creating the logger:

```go
// Default: "storage/files/logs"
logger.NewLogger(store, "custom/log/path")

// Use default path
logger.NewLogger(store, "")
```

## Cleanup

- Log files are automatically closed when job runs complete
- The logger system properly closes all active loggers on shutdown
- Database logs are cleaned up automatically via foreign key constraints 
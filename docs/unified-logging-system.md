# Unified Logging System

The Stratal unified logging system consolidates all logs into a single `logs` table, supporting three different log types: **System**, **Job**, and **Task** logs. This provides a comprehensive view of all system activity while maintaining the ability to filter and query specific log types.

## Overview

The new unified logging system replaces the previous `task_logs` table with a more flexible `logs` table that can handle different types of logs with optional relationships to job runs and task runs.

## Database Schema

### Logs Table Structure

```sql
CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    type TEXT CHECK (type IN ('system', 'job', 'task')) NOT NULL DEFAULT 'system',
    job_run_id UUID REFERENCES job_runs (id) ON DELETE CASCADE NULL,
    task_run_id UUID REFERENCES task_runs (id) ON DELETE CASCADE NULL,
    timestamp TIMESTAMP DEFAULT now(),
    level TEXT CHECK (level IN ('info', 'error', 'warn', 'debug')) NOT NULL DEFAULT 'info',
    stream TEXT CHECK (stream IN ('stdout', 'stderr', 'system')) NOT NULL DEFAULT 'system',
    message TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Indexes

```sql
-- Performance indexes
CREATE INDEX idx_logs_type ON logs (type);
CREATE INDEX idx_logs_job_run_id ON logs (job_run_id);
CREATE INDEX idx_logs_task_run_id ON logs (task_run_id);
CREATE INDEX idx_logs_timestamp ON logs (timestamp);
CREATE INDEX idx_logs_level ON logs (level);

-- Composite indexes for common queries
CREATE INDEX idx_logs_job_run_timestamp ON logs (job_run_id, timestamp);
CREATE INDEX idx_logs_task_run_timestamp ON logs (task_run_id, timestamp);
```

## Log Types

### 1. System Logs (`type = 'system'`)
- **Purpose**: General system-level events (server startup, shutdown, configuration changes)
- **Relationships**: No direct relationship to jobs or tasks
- **Fields**: `job_run_id` and `task_run_id` are NULL
- **Examples**: Server startup, configuration changes, health checks

```go
logger.LogSystem(logger.InfoLevel, "Server started successfully", map[string]interface{}{
    "port": 8080,
    "version": "1.0.0",
})
```

### 2. Job Logs (`type = 'job'`)
- **Purpose**: Job-level events (job start, completion, errors)
- **Relationships**: Connected to a specific job run via `job_run_id`
- **Fields**: `job_run_id` is set, `task_run_id` is NULL
- **Examples**: Job started, job completed, job failed

```go
jobLogger.LogJob(logger.InfoLevel, "Job execution started", map[string]interface{}{
    "total_tasks": 5,
    "estimated_duration": "2m",
})
```

### 3. Task Logs (`type = 'task'`)
- **Purpose**: Individual task execution logs
- **Relationships**: Connected to both job run and task run via `job_run_id` and `task_run_id`
- **Fields**: Both `job_run_id` and `task_run_id` are set
- **Examples**: Task started, task output, task errors

```go
jobLogger.LogTask(taskRunID, logger.InfoLevel, "HTTP request completed", "stdout", map[string]interface{}{
    "status_code": 200,
    "response_time": "150ms",
})
```

## API Endpoints

### Get Logs by Job Run
```http
GET /api/v1/logs/job-runs/{job_run_id}?limit=100&offset=0
```
Returns all logs (job and task) associated with a specific job run.

### Get Logs by Task Run
```http
GET /api/v1/logs/task-runs/{task_run_id}
```
Returns all logs associated with a specific task run.

### Get System Logs
```http
GET /api/v1/logs/system?limit=100&offset=0
```
Returns system-level logs.

### Get Logs by Type
```http
GET /api/v1/logs/type/{type}?limit=100&offset=0
```
Returns logs of a specific type (`system`, `job`, or `task`).

### Streaming Endpoints
```http
GET /api/v1/logs/stream/ws?job_run_id={uuid}    # WebSocket
GET /api/v1/logs/stream/sse?job_run_id={uuid}   # Server-Sent Events
GET /api/v1/logs/stream/status                   # Connection status
```

## Usage Examples

### Backend Logging

#### System Logging
```go
// Initialize logger
logger := logger.NewLogger(store, "internal/storage/files/logs")

// Log system events
logger.LogSystem(logger.InfoLevel, "Database connection established", map[string]interface{}{
    "host": "localhost",
    "database": "stratal",
})

logger.LogSystem(logger.ErrorLevel, "Configuration validation failed", map[string]interface{}{
    "error": "invalid port number",
    "config_file": "/etc/stratal/config.yaml",
})
```

#### Job Logging
```go
// Get job run logger
jobLogger, err := logger.GetJobRunLogger(jobRunID)
if err != nil {
    return err
}
defer logger.CloseJobRunLogger(jobRunID)

// Log job events
jobLogger.LogJob(logger.InfoLevel, "Job started", map[string]interface{}{
    "job_name": "data_processing",
    "parameters": params,
})

jobLogger.LogJob(logger.ErrorLevel, "Job failed", map[string]interface{}{
    "error": err.Error(),
    "duration": time.Since(startTime).String(),
})
```

#### Task Logging
```go
// Log task events
jobLogger.LogTask(taskRunID, logger.InfoLevel, "Task started", "stdout", map[string]interface{}{
    "task_name": "http_request",
    "url": "https://api.example.com",
})

jobLogger.LogTask(taskRunID, logger.ErrorLevel, "Task failed", "stderr", map[string]interface{}{
    "error": "connection timeout",
    "retry_count": 3,
})
```

### Frontend Integration

#### TypeScript Interface
```typescript
interface LogMessage {
  type: 'system' | 'job' | 'task';
  job_run_id: string;
  task_run_id?: string;
  timestamp: string;
  level: 'info' | 'error' | 'warn' | 'debug';
  stream: 'stdout' | 'stderr' | 'system';
  message: string;
  metadata?: Record<string, any>;
}
```

#### Fetching Logs
```typescript
import { LogStreamingService } from './utils/LogStreamer';

const logService = new LogStreamingService();

// Get job run logs
const jobLogs = await logService.getJobRunLogs('job-123', 100, 0);

// Get system logs
const systemLogs = await logService.getSystemLogs(50, 0);

// Get logs by type
const taskLogs = await logService.getLogsByType('task', 100, 0);
```

#### Real-time Streaming
```typescript
// Stream all log types
const streamer = logService.streamAllLogs({
  onMessage: (log: LogMessage) => {
    switch (log.type) {
      case 'system':
        console.log(`[SYSTEM] ${log.message}`);
        break;
      case 'job':
        console.log(`[JOB ${log.job_run_id}] ${log.message}`);
        break;
      case 'task':
        console.log(`[TASK ${log.task_run_id}] ${log.message}`);
        break;
    }
  }
});
```

#### React Component Usage
```tsx
import LogViewer from './components/LogViewer';

function SystemDashboard() {
  return (
    <div>
      <h2>System Logs</h2>
      <LogViewer filterType="system" height="400px" />
      
      <h2>Job Execution</h2>
      <LogViewer jobRunId="job-123" height="600px" />
    </div>
  );
}
```

## Migration from task_logs

The migration automatically:
1. Drops the existing `task_logs` table
2. Creates the new unified `logs` table
3. Sets up appropriate indexes for performance

### Migration Files
- **Up**: `000002_column-job-run-logs.up.sql`
- **Down**: `000002_column-job-run-logs.down.sql`

To apply the migration:
```bash
migrate -path internal/storage/db/migration -database "postgresql://..." up
```

To rollback:
```bash
migrate -path internal/storage/db/migration -database "postgresql://..." down 1
```

## Performance Considerations

### Indexes
The new table includes optimized indexes for common query patterns:
- Type-based filtering
- Job run associations
- Task run associations
- Time-based queries
- Log level filtering

### Query Optimization
- Use type-specific endpoints when possible
- Implement pagination for large result sets
- Consider log retention policies for long-running systems

### Storage Management
- Monitor table size growth
- Implement log rotation/archival strategies
- Use appropriate log levels to control volume

## Monitoring and Alerting

### Log Volume Metrics
```sql
-- Count logs by type
SELECT type, COUNT(*) 
FROM logs 
GROUP BY type;

-- Count logs by level in last hour
SELECT level, COUNT(*) 
FROM logs 
WHERE timestamp > NOW() - INTERVAL '1 hour'
GROUP BY level;
```

### Error Rate Monitoring
```sql
-- Error rate by job runs
SELECT 
  job_run_id,
  COUNT(CASE WHEN level = 'error' THEN 1 END) as errors,
  COUNT(*) as total_logs,
  (COUNT(CASE WHEN level = 'error' THEN 1 END)::float / COUNT(*)) * 100 as error_rate
FROM logs 
WHERE type IN ('job', 'task')
  AND timestamp > NOW() - INTERVAL '1 hour'
GROUP BY job_run_id
HAVING COUNT(CASE WHEN level = 'error' THEN 1 END) > 0;
```

## Best Practices

### Logging Strategy
1. **System Logs**: Use for infrastructure events, not application logic
2. **Job Logs**: Log job lifecycle events and high-level status
3. **Task Logs**: Detailed execution information and task-specific data

### Metadata Usage
- Include relevant context in metadata JSON field
- Keep metadata structured and searchable
- Avoid sensitive information in metadata

### Log Levels
- **Debug**: Detailed diagnostic information
- **Info**: General informational messages
- **Warn**: Warning conditions that might need attention
- **Error**: Error conditions that need immediate attention

### Performance Tips
- Use appropriate log levels to control volume
- Implement client-side filtering to reduce bandwidth
- Consider log sampling for high-frequency events
- Monitor database performance and tune indexes as needed

This unified logging system provides comprehensive visibility into all aspects of the Stratal system while maintaining performance and flexibility for different use cases. 
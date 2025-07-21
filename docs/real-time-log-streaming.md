# Real-Time Log Streaming

The Stratal log streaming system provides real-time log delivery to client applications using both WebSockets and Server-Sent Events (SSE). This enables live monitoring of job execution with immediate feedback.

## Features

- **Dual Protocol Support**: WebSockets and Server-Sent Events (SSE)
- **Job-Specific Streaming**: Stream logs for specific job runs
- **Global Streaming**: Stream logs from all job runs
- **Automatic Reconnection**: Robust reconnection with exponential backoff
- **Real-Time Filtering**: Filter logs by level, stream, and search terms
- **Historical Logs**: Access to stored logs via REST API
- **File Downloads**: Download complete log files

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │   HTTP Server   │    │     Worker      │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │ LogViewer   │◄┼────┼►│ LogStreamer │◄┼────┼►│   Logger    │ │
│ │ Component   │ │    │ │   Handler   │ │    │ │   System    │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │LogStreaming │ │    │ │    API      │ │    │ │   Database  │ │
│ │  Service    │◄┼────┼►│  Endpoints  │◄┼────┼►│    Logs     │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## API Endpoints

### Streaming Endpoints

#### WebSocket Streaming
```
GET /api/v1/logs/stream/ws?job_run_id={uuid}
```

**Parameters:**
- `job_run_id` (optional): Stream logs for a specific job run. Omit for global streaming.

**Example:**
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/logs/stream/ws?job_run_id=abc-123');

ws.onmessage = (event) => {
  const logMessage = JSON.parse(event.data);
  console.log(logMessage);
};
```

#### Server-Sent Events
```
GET /api/v1/logs/stream/sse?job_run_id={uuid}
```

**Parameters:**
- `job_run_id` (optional): Stream logs for a specific job run. Omit for global streaming.

**Example:**
```javascript
const eventSource = new EventSource('/api/v1/logs/stream/sse?job_run_id=abc-123');

eventSource.addEventListener('log', (event) => {
  const logMessage = JSON.parse(event.data);
  console.log(logMessage);
});
```

### REST API Endpoints

#### Get Job Run Logs
```
GET /api/v1/logs/job-runs/{job_run_id}?limit=100&offset=0
```

**Response:**
```json
{
  "logs": [
    {
      "id": 1,
      "task_run_id": "def-456",
      "timestamp": "2024-01-15T10:30:45.123Z",
      "stream": "stdout",
      "message": "Task started successfully"
    }
  ],
  "pagination": {
    "total": 250,
    "limit": 100,
    "offset": 0,
    "count": 100
  }
}
```

#### Get Task Run Logs
```
GET /api/v1/logs/task-runs/{task_run_id}
```

#### Download Log File
```
GET /api/v1/logs/job-runs/{job_run_id}/download?date=2024-01-15
```

#### Streaming Status
```
GET /api/v1/logs/stream/status
```

**Response:**
```json
{
  "status": "active",
  "connections": {
    "global": {
      "websocket": 3,
      "sse": 2
    },
    "jobs": {
      "abc-123": {
        "websocket": 1,
        "sse": 0
      }
    }
  }
}
```

## Log Message Format

```typescript
interface LogMessage {
  job_run_id: string;
  task_run_id?: string;
  timestamp: string;
  level: 'info' | 'error' | 'warn' | 'debug';
  stream: 'stdout' | 'stderr' | 'system';
  message: string;
  task_name?: string;
}
```

**Example:**
```json
{
  "job_run_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "task_run_id": "b2c3d4e5-f6g7-8901-bcde-f23456789012",
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "info",
  "stream": "stdout",
  "message": "Task execution completed successfully",
  "task_name": "http_request_task"
}
```

## Client-Side Usage

### TypeScript/JavaScript

```typescript
import { LogStreamingService } from './utils/LogStreamer';

const logService = new LogStreamingService();

// Stream logs for a specific job run
const streamer = logService.streamJobRunLogs('abc-123', {
  onMessage: (log) => {
    console.log(`[${log.level}] ${log.message}`);
  },
  onConnect: () => {
    console.log('Connected to log stream');
  },
  onDisconnect: () => {
    console.log('Disconnected from log stream');
  },
  onError: (error) => {
    console.error('Streaming error:', error);
  }
});

// Stop streaming when done
streamer.disconnect();
```

### React Component

```typescript
import React from 'react';
import LogViewer from './components/LogViewer';

function JobDetailsPage({ jobRunId }: { jobRunId: string }) {
  return (
    <div>
      <h1>Job Details</h1>
      <LogViewer 
        jobRunId={jobRunId}
        height="600px"
        autoScroll={true}
      />
    </div>
  );
}
```

### Plain JavaScript

```javascript
// Using WebSocket directly
const ws = new WebSocket('ws://localhost:8080/api/v1/logs/stream/ws');

ws.onopen = () => {
  console.log('Connected to log stream');
};

ws.onmessage = (event) => {
  const log = JSON.parse(event.data);
  const logElement = document.createElement('div');
  logElement.textContent = `[${log.timestamp}] [${log.level}] ${log.message}`;
  document.getElementById('logs').appendChild(logElement);
};

// Using Server-Sent Events
const eventSource = new EventSource('/api/v1/logs/stream/sse');

eventSource.addEventListener('log', (event) => {
  const log = JSON.parse(event.data);
  console.log(log);
});
```

## Configuration

### Server Configuration

```go
// In your HTTP server setup
if httpServer.logSystem != nil {
    logsHandler := NewLogsHandler(httpServer.store, httpServer.logSystem.GetStreamer())
    logsHandler.RegisterLogRoutes(v1)
}
```

### Client Configuration

```typescript
const logService = new LogStreamingService('http://localhost:8080', {
  autoReconnect: true,
  maxReconnectAttempts: 5,
  onError: (error) => {
    console.error('Log streaming error:', error);
  }
});
```

## Performance Considerations

### Server-Side
- **Connection Limits**: Monitor active connections to prevent resource exhaustion
- **Memory Management**: Implement log rotation and cleanup for old connections
- **Rate Limiting**: Consider rate limiting for high-frequency log producers

### Client-Side
- **Buffer Management**: Limit the number of logs kept in memory (default: 1000)
- **Auto-Scroll**: Disable auto-scroll for better performance with high log volumes
- **Filtering**: Use client-side filtering to reduce rendering overhead

## Security

### Authentication
```typescript
// Add authentication headers to WebSocket connections
const ws = new WebSocket('ws://localhost:8080/api/v1/logs/stream/ws', [], {
  headers: {
    'Authorization': 'Bearer ' + token
  }
});
```

### CORS Configuration
```go
// Configure CORS for SSE endpoints
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Headers", "Authorization")
```

## Monitoring and Debugging

### Connection Status
```typescript
// Check streaming status
const status = await logService.getStreamingStatus();
console.log('Active connections:', status.connections);
```

### Debug Mode
```typescript
// Enable debug logging
const streamer = logService.streamJobRunLogs(jobRunId, {
  onMessage: (log) => {
    if (log.level === 'debug') {
      console.debug('Debug log:', log.message);
    }
  }
});
```

## Error Handling

### Common Issues

1. **Connection Drops**: Automatic reconnection with exponential backoff
2. **Message Parse Errors**: Graceful handling of malformed messages
3. **Browser Limits**: WebSocket/SSE connection limits per domain

### Error Recovery

```typescript
const streamer = logService.streamJobRunLogs(jobRunId, {
  onError: (error) => {
    console.error('Stream error:', error);
    
    // Implement custom retry logic if needed
    setTimeout(() => {
      streamer.connect();
    }, 5000);
  },
  autoReconnect: true,
  maxReconnectAttempts: 10
});
```

## Best Practices

1. **Use WebSockets** for real-time applications requiring bidirectional communication
2. **Use SSE** for simple real-time updates and better browser compatibility
3. **Implement proper cleanup** to avoid memory leaks
4. **Filter logs client-side** to reduce bandwidth usage
5. **Use job-specific streaming** instead of global streaming when possible
6. **Implement exponential backoff** for reconnection attempts
7. **Monitor connection counts** to prevent resource exhaustion

## Integration Examples

### Job Execution Dashboard
```typescript
import { LogStreamingService } from './utils/LogStreamer';
import LogViewer from './components/LogViewer';

function Dashboard() {
  const [activeJobs, setActiveJobs] = useState([]);
  const logService = new LogStreamingService();

  // Stream all logs for monitoring
  useEffect(() => {
    const streamer = logService.streamAllLogs({
      onMessage: (log) => {
        // Update job status based on log messages
        updateJobStatus(log.job_run_id, log);
      }
    });

    return () => streamer.disconnect();
  }, []);

  return (
    <div>
      {activeJobs.map(job => (
        <LogViewer key={job.id} jobRunId={job.id} />
      ))}
    </div>
  );
}
```

This real-time log streaming system provides comprehensive monitoring capabilities for job execution, enabling immediate visibility into system behavior and quick debugging of issues. 
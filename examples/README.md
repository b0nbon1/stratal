# Stratal Examples

This directory contains examples and test files for the Stratal job processing system.

## Files

- `simple_parallel_example.json` - Example job with parallel task execution and dependencies
- `job_with_secrets_example.json` - Example job that uses encrypted secrets
- `insomnia_test_example.json` - Complete Insomnia test collection for testing the system
- `README.md` - This file

## Getting Started with Insomnia Testing

### Prerequisites

1. **Start the Stratal system:**
   ```bash
   # Start PostgreSQL (if using Docker)
   make postgres
   make createdb
   make migrateup

   # Start Redis (if using Docker)
   docker run --name redis -p 6379:6379 -d redis

   # Set required environment variables
   export ENCRYPTION_KEY="your-32-character-encryption-key-here"
   export DB_PASSWORD="1234567890"  # or your actual DB password

   # Start the server
   make run_server_dev

   # In another terminal, start the worker
   make run_worker_dev
   ```

2. **Install Insomnia:** Download from [insomnia.rest](https://insomnia.rest/)

### Import Test Collection

1. Open Insomnia
2. Click "Import/Export" or use `Ctrl+O` (Windows/Linux) or `Cmd+O` (Mac)
3. Select "Import Data" and choose the `insomnia_test_example.json` file
4. The collection "Stratal Job System Test" will be imported with all test requests

### Test Workflow

The test collection is organized in a logical sequence. Follow these steps:

#### 1. Health Checks
- **1. Health Check** - Verify server is running
- **2. API Health Check** - Verify API endpoints are working

#### 2. Create Jobs
- **3. Create Simple Parallel Job** - Creates a job with multiple parallel tasks that demonstrate:
  - Parallel execution (3 random number generators run simultaneously)
  - Task dependencies (sum calculation waits for all generators)
  - Output passing between tasks
- **4. Create Job with Builtin Tasks** - Creates a job using builtin tasks:
  - HTTP GET request to httpbin.org
  - HTTP POST request with JSON payload
  - Echo task demonstration

#### 3. Job Management
- **5. List All Jobs** - View all created jobs
- **6. Get Job Details** - Get specific job information (update `job_id` variable)

#### 4. Job Execution
- **7. Create Job Run** - Execute a previously created job (update `job_id` variable)
- **8. Get Job Run Status** - Monitor job execution status (update `job_run_id` variable)

#### 5. Immediate Execution
- **9. Create and Run Job Immediately** - Create and execute a job in one step

### Environment Variables

The collection includes environment variables that you can customize:

- `base_url`: Default is `http://localhost:8080` (change if server runs on different port)
- `job_id`: Replace with actual job ID from job creation responses
- `job_run_id`: Replace with actual job run ID from job run creation responses

### Example Test Flow

1. **Start by running requests 1-2** to verify your system is running
2. **Run request 3** to create a test job. Copy the returned `id` value
3. **Update the `job_id` environment variable** with the copied ID
4. **Run request 6** to see job details with tasks
5. **Run request 7** to create a job run. Copy the returned `JobRunID`
6. **Update the `job_run_id` environment variable** with the copied ID  
7. **Run request 8** repeatedly to monitor job execution status
8. **Try request 9** for immediate job execution

### Understanding Task Types

The system supports different task types:

#### Custom Tasks (`type: "custom"`)
Execute custom scripts in various languages:
```json
{
  "name": "my_script",
  "type": "custom", 
  "order": 1,
  "config": {
    "script": {
      "language": "bash",  // or "python", "node", etc.
      "code": "echo 'Hello World'"
    }
  }
}
```

#### Builtin Tasks (`type: "builtin"`)
Use predefined system tasks:
```json
{
  "name": "api_call",
  "type": "builtin",
  "order": 1, 
  "config": {
    "parameters": {
      "task_name": "http_request",
      "url": "https://api.example.com",
      "method": "GET"
    }
  }
}
```

Available builtin tasks:
- `http_request` - Make HTTP requests
- `send_email` - Send emails via SMTP
- `echo` - Simple echo task for testing

### Task Dependencies and Ordering

Tasks are executed based on:
1. **Order**: Tasks with the same order run in parallel
2. **Dependencies**: `depends_on` ensures prerequisite tasks complete first

Example:
```json
{
  "name": "dependent_task",
  "order": 2,
  "config": {
    "depends_on": ["task1", "task2"],  // Wait for these to complete
    "script": {
      "language": "bash",
      "code": "echo 'Previous tasks completed'"
    }
  }
}
```

### Output Passing Between Tasks

Task outputs are available to subsequent tasks via environment variables:
- Format: `TASK_OUTPUT_<TASK_NAME_UPPERCASE>`
- Example: Output from `generate_random_1` is available as `TASK_OUTPUT_GENERATE_RANDOM_1`

### Job Run Status Values

Monitor these status values when checking job run progress:
- `pending` - Job run created, waiting to be queued
- `queued` - Job run is in the queue, waiting for worker
- `running` - Job run is being executed by a worker
- `completed` - Job run finished successfully
- `failed` - Job run encountered an error

### Troubleshooting

#### Common Issues:

1. **Connection refused**: Ensure the server is running on `http://localhost:8080`
2. **Job creation fails**: Check that all required fields are provided
3. **Job run not processing**: Ensure the worker is running (`make run_worker_dev`)
4. **Database errors**: Verify PostgreSQL is running and migrations are applied
5. **Queue errors**: Ensure Redis is running

#### Checking Logs:

Monitor the server and worker terminal outputs for detailed execution logs and error messages.

### Advanced Testing

For more complex scenarios, you can:

1. **Add secrets**: Use the secrets API to store encrypted values
2. **Parallel job execution**: Create multiple job runs simultaneously  
3. **Long-running jobs**: Create jobs with sleep commands to test cancellation
4. **Error handling**: Create jobs that intentionally fail to test error scenarios

### Next Steps

After successful testing:
1. Explore the web UI at `http://localhost:3000` (if running frontend)
2. Check the codebase structure to understand the system architecture
3. Create your own custom builtin tasks
4. Integrate with external systems using HTTP requests and secrets 
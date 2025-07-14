# Stratal Job Examples

This directory contains examples demonstrating various features of the Stratal job automation system.

## Parallel Task Execution & Output Passing

Stratal now supports:
1. **Parallel execution** of tasks that don't have dependencies
2. **Output passing** between dependent tasks
3. **Flexible dependency references** using either task names or IDs

### Key Features

#### 1. Parallel Execution
Tasks at the same dependency level execute concurrently:
- Tasks with no dependencies run in parallel
- Tasks depending only on level-0 tasks run in parallel after level-0 completes
- Maximum 5 concurrent tasks by default (configurable)

#### 2. Output Passing

**For Builtin Tasks:**
Use `${task_name.output}` syntax in parameters:
```json
{
  "parameters": {
    "subject": "Result: ${previous_task.output}",
    "body": "The calculation returned: ${calculate_task.output}"
  }
}
```

**For Custom Scripts:**
Access outputs via environment variables with `TASK_OUTPUT_` prefix:
```python
# Python example
import os
previous_output = os.environ.get('TASK_OUTPUT_PREVIOUS_TASK', '')
```

```javascript
// JavaScript example
const previousOutput = process.env.TASK_OUTPUT_PREVIOUS_TASK || '';
```

```bash
# Bash example
echo "Previous result: $TASK_OUTPUT_PREVIOUS_TASK"
```

#### 3. Dependency Declaration
Use task names (recommended) or task IDs:
```json
{
  "depends_on": ["fetch_data", "validate_input"]
}
```

### Examples

1. **simple_parallel_example.json** - Basic demonstration:
   - 3 tasks generate random numbers in parallel
   - 1 task sums the results
   - Final task checks if sum > 150

2. **parallel_tasks_example.json** - Complex pipeline:
   - 3 API calls execute in parallel (user, posts, weather)
   - Data processing task depends on user data
   - Aggregation task waits for all 3 API calls
   - Email notification uses aggregated results

### Running Examples

```bash
# Using curl
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d @examples/simple_parallel_example.json

# Check job status
curl http://localhost:8080/api/v1/jobs/{job_id}

# View job runs
curl http://localhost:8080/api/v1/jobs/{job_id}/runs

# View task runs for a specific job run
curl http://localhost:8080/api/v1/job-runs/{job_run_id}/tasks
```

### Performance Benefits

Parallel execution significantly reduces total execution time:
- **Sequential**: Task1 (2s) → Task2 (2s) → Task3 (2s) = 6s total
- **Parallel**: Task1, Task2, Task3 (2s) = 2s total

### Best Practices

1. **Group independent tasks** at the same order level
2. **Use meaningful task names** for easier dependency references
3. **Handle missing outputs gracefully** in custom scripts
4. **Keep task outputs concise** - they're stored in memory during execution
5. **Use structured output** (JSON) for complex data passing between tasks

### Task Execution Flow

```
Level 0: [Task A] [Task B] [Task C]  <- Execute in parallel
           ↓         ↓         ↓
Level 1: [Task D]  [Task E]          <- Execute in parallel after level 0
           ↓         ↓
Level 2: [Task F]                    <- Executes after level 1
```

### Limitations

- Maximum 5 concurrent tasks (configurable in processor)
- Output size should be reasonable (stored in memory)
- Circular dependencies are detected and rejected
- Task outputs are only available within the same job run 
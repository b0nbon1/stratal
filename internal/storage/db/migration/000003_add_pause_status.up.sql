-- Add 'paused' status to job_runs table
ALTER TABLE job_runs DROP CONSTRAINT IF EXISTS job_runs_status_check;
ALTER TABLE job_runs ADD CONSTRAINT job_runs_status_check CHECK (
    status IN ('pending', 'queued', 'running', 'paused', 'failed', 'completed')
);

-- Add 'paused' status to task_runs table  
ALTER TABLE task_runs DROP CONSTRAINT IF EXISTS task_runs_status_check;
ALTER TABLE task_runs ADD CONSTRAINT task_runs_status_check CHECK (
    status IN ('pending', 'running', 'paused', 'failed', 'completed')
);

-- Add paused_at timestamp to job_runs for tracking when job was paused
ALTER TABLE job_runs ADD COLUMN paused_at TIMESTAMP;

-- Add paused_at timestamp to task_runs for tracking when task was paused
ALTER TABLE task_runs ADD COLUMN paused_at TIMESTAMP;
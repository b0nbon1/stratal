-- Remove paused_at columns
ALTER TABLE task_runs DROP COLUMN IF EXISTS paused_at;
ALTER TABLE job_runs DROP COLUMN IF EXISTS paused_at;

-- Revert job_runs status constraint to original
ALTER TABLE job_runs DROP CONSTRAINT IF EXISTS job_runs_status_check;
ALTER TABLE job_runs ADD CONSTRAINT job_runs_status_check CHECK (
    status IN ('pending', 'queued', 'running', 'failed', 'completed')
);

-- Revert task_runs status constraint to original
ALTER TABLE task_runs DROP CONSTRAINT IF EXISTS task_runs_status_check;
ALTER TABLE task_runs ADD CONSTRAINT task_runs_status_check CHECK (
    status IN ('pending', 'running', 'failed', 'completed')
);
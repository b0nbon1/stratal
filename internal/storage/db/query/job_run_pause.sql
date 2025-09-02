-- name: PauseJobRun :exec
UPDATE job_runs 
SET status = 'paused', 
    paused_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND status IN ('running', 'queued');

-- name: ResumeJobRun :exec  
UPDATE job_runs
SET status = CASE 
    WHEN paused_at IS NOT NULL AND started_at IS NULL THEN 'queued'
    ELSE 'running'
END,
paused_at = NULL,
updated_at = NOW()
WHERE id = $1 AND status = 'paused';

-- name: GetPausedJobRuns :many
SELECT * FROM job_runs 
WHERE status = 'paused'
ORDER BY paused_at DESC;

-- name: PauseTaskRun :exec
UPDATE task_runs
SET status = 'paused',
    paused_at = NOW(), 
    updated_at = NOW()
WHERE id = $1 AND status = 'running';

-- name: ResumeTaskRun :exec
UPDATE task_runs  
SET status = 'running',
    paused_at = NULL,
    updated_at = NOW()
WHERE id = $1 AND status = 'paused';

-- name: GetJobRunWithPauseInfo :one
SELECT jr.*, j.name as job_name
FROM job_runs jr
JOIN jobs j ON jr.job_id = j.id  
WHERE jr.id = $1;
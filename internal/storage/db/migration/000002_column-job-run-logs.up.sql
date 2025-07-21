-- Drop the existing task_logs table
DROP TABLE IF EXISTS task_logs;

-- Create the new unified logs table
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

-- Create indexes for performance
CREATE INDEX idx_logs_type ON logs (type);
CREATE INDEX idx_logs_job_run_id ON logs (job_run_id);
CREATE INDEX idx_logs_task_run_id ON logs (task_run_id);
CREATE INDEX idx_logs_timestamp ON logs (timestamp);
CREATE INDEX idx_logs_level ON logs (level);

-- Create a composite index for common queries
CREATE INDEX idx_logs_job_run_timestamp ON logs (job_run_id, timestamp);
CREATE INDEX idx_logs_task_run_timestamp ON logs (task_run_id, timestamp);

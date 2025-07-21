-- Drop the unified logs table
DROP TABLE IF EXISTS logs;

-- Recreate the original task_logs table
CREATE TABLE task_logs (
    id BIGSERIAL PRIMARY KEY,
    task_run_id UUID REFERENCES task_runs (id) ON DELETE CASCADE,
    timestamp TIMESTAMP DEFAULT now(),
    stream TEXT CHECK (stream IN ('stdout', 'stderr')) NOT NULL,
    message TEXT NOT NULL
);

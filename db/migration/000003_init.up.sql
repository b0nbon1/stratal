CREATE TABLE
    job_logs (
        id SERIAL PRIMARY KEY,
        job_id UUID NOT NULL REFERENCES jobs (id),
        log_level TEXT NOT NULL,
        message TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_job_logs_job_id ON task_logs (task_id);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users (id) ON DELETE CASCADE NULL,
    name TEXT NOT NULL,
    description TEXT,
    source TEXT CHECK (source IN ('api', 'cli', 'sdk', 'ui')) NOT NULL,
    raw_payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID REFERENCES jobs (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    config JSONB NOT NULL,
    "order" INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE job_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID REFERENCES jobs (id) ON DELETE CASCADE,
    status TEXT CHECK (
        status IN ('pending', 'queued', 'running', 'failed', 'completed')
    ),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    error_message TEXT,
    triggered_by TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE task_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_run_id UUID REFERENCES job_runs (id) ON DELETE CASCADE,
    task_id UUID REFERENCES tasks (id) ON DELETE CASCADE,
    status TEXT CHECK (
        status IN ('pending', 'running', 'failed', 'completed')
    ),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    exit_code INTEGER,
    output TEXT,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE task_logs (
    id BIGSERIAL PRIMARY KEY,
    task_run_id UUID REFERENCES task_runs (id) ON DELETE CASCADE,
    timestamp TIMESTAMP DEFAULT now (),
    stream TEXT CHECK (stream IN ('stdout', 'stderr')) NOT NULL,
    message TEXT NOT NULL
);


CREATE TABLE secrets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    encrypted_value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE job_status AS ENUM ('pending', 'running', 'failed', 'completed');

CREATE TABLE jobs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  schedule VARCHAR(255),
  type VARCHAR(255),
  config JSONB,
  status job_status DEFAULT 'pending',
  retries INTEGER DEFAULT 0,
  max_retries INTEGER,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE job_runs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  job_id UUID REFERENCES jobs(id),
  status job_status DEFAULT 'pending',
  logs TEXT,
  started_at TIMESTAMP WITH TIME ZONE,
  ended_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE job_runs
ADD CONSTRAINT fk_job_runs_job
FOREIGN KEY (job_id) REFERENCES jobs(id);

-- Down Migration: 000002_init.down.sql

-- Drop index on jobs.user_id if it exists
DROP INDEX IF EXISTS idx_jobs_user_id;

-- Remove user_id column from jobs table
ALTER TABLE jobs
DROP COLUMN IF EXISTS user_id;

-- Drop users table
DROP TABLE IF EXISTS users;
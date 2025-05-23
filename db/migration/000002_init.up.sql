-- Up Migration: 000002_init.up.sql

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add user_id column to jobs table and set up foreign key
ALTER TABLE jobs
ADD COLUMN user_id INTEGER REFERENCES users(id);

-- Optionally, you can add an index for faster lookups
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
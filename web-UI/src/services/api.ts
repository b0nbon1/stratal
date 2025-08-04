import axios from 'axios';

// Base URL for API calls
const API_BASE_URL = 'http://localhost:8080/api/v1';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Types
export interface Job {
  id: string;
  name: string;
  description?: string;
  source: string;
  raw_payload?: any;
  created_at: string;
  updated_at: string;
  tasks?: Task[];
}

export interface Task {
  id: string;
  job_id: string;
  name: string;
  type: string;
  order: number;
  config: any;
  created_at: string;
  updated_at: string;
}

export interface JobRun {
  id: string;
  job_id: string;
  status: string;
  triggered_by: string;
  created_at: string;
  updated_at: string;
  started_at?: string;
  completed_at?: string;
  task_runs?: TaskRun[];
}

export interface TaskRun {
  id: string;
  job_run_id: string;
  task_id: string;
  status: string;
  created_at: string;
  updated_at: string;
  started_at?: string;
  completed_at?: string;
}

export interface Secret {
  id: string;
  name: string;
  created_at: string;
}

export interface LogEntry {
  id: string;
  type: 'system' | 'job' | 'task';
  job_run_id?: string;
  task_run_id?: string;
  timestamp: string;
  level: 'info' | 'error' | 'warn' | 'debug';
  stream: 'stdout' | 'stderr' | 'system';
  message: string;
  task_name?: string;
  metadata?: Record<string, any>;
}

export interface CreateJobRequest {
  name: string;
  description?: string;
  source?: string;
  raw_payload?: any;
  tasks: {
    name: string;
    type: string;
    config: any;
    order: number;
  }[];
  run_immediately?: boolean;
}

export interface CreateJobRunRequest {
  job_id: string;
  triggered_by: string;
}

export interface CreateSecretRequest {
  name: string;
  value: string;
}

// Job API
export const jobsApi = {
  // Create a new job
  create: async (data: CreateJobRequest): Promise<{ job: Job; tasks: Task[]; message: string; job_run_id?: string }> => {
    const response = await api.post('/jobs', data);
    return response.data;
  },

  // Get all jobs with pagination
  list: async (params?: { limit?: number; offset?: number }): Promise<{ jobs: Job[]; limit: number; offset: number; count: number }> => {
    const response = await api.get('/jobs', { params });
    return response.data;
  },

  // Get a specific job by ID
  get: async (id: string): Promise<Job> => {
    const response = await api.get(`/jobs/${id}`);
    return response.data;
  },
};

// Job Run API
export const jobRunsApi = {
  // Create a new job run
  create: async (data: CreateJobRunRequest): Promise<{ message: string; JobRunID: string }> => {
    const response = await api.post('/job-runs', data);
    return response.data;
  },

  // Get job run by ID
  get: async (id: string): Promise<JobRun> => {
    const response = await api.get(`/job-runs/${id}`);
    return response.data;
  },

  // List job runs
  list: async (params?: { limit?: number; offset?: number }): Promise<{ job_runs: JobRun[]; limit: number; offset: number; count: number }> => {
    const response = await api.get('/job-runs', { params });
    return response.data;
  },
};

// Secrets API
export const secretsApi = {
  // Create a new secret
  create: async (data: CreateSecretRequest): Promise<Secret> => {
    const response = await api.post('/secrets', data);
    return response.data;
  },

  // List all secrets
  list: async (): Promise<{ secrets: Secret[] }> => {
    const response = await api.get('/secrets');
    return response.data;
  },

  // Delete a secret
  delete: async (id: string): Promise<{ message: string }> => {
    const response = await api.delete(`/secrets?id=${id}`);
    return response.data;
  },

  // Update a secret
  update: async (id: string, value: string): Promise<{ message: string }> => {
    const response = await api.put(`/secrets?id=${id}`, { value });
    return response.data;
  },
};

// Logs API
export const logsApi = {
  // Get job run logs
  getJobRunLogs: async (jobRunId: string, params?: { limit?: number; offset?: number }): Promise<{ logs: LogEntry[]; pagination: { total: number; limit: number; offset: number; count: number } }> => {
    const response = await api.get(`/logs/job-runs/${jobRunId}`, { params });
    return response.data;
  },

  // Get task run logs
  getTaskRunLogs: async (taskRunId: string): Promise<{ logs: LogEntry[] }> => {
    const response = await api.get(`/logs/task-runs/${taskRunId}`);
    return response.data;
  },

  // Get system logs
  getSystemLogs: async (params?: { limit?: number; offset?: number }): Promise<{ logs: LogEntry[]; pagination: { limit: number; offset: number; count: number } }> => {
    const response = await api.get('/logs/system', { params });
    return response.data;
  },

  // Get logs by type
  getLogsByType: async (type: 'system' | 'job' | 'task', params?: { limit?: number; offset?: number }): Promise<{ logs: LogEntry[]; pagination: { limit: number; offset: number; count: number } }> => {
    const response = await api.get(`/logs/type/${type}`, { params });
    return response.data;
  },

  // Get streaming status
  getStreamingStatus: async (): Promise<{ status: string; connections: any }> => {
    const response = await api.get('/logs/stream/status');
    return response.data;
  },

  // Download job run logs
  downloadJobRunLogs: async (jobRunId: string, date?: string): Promise<Blob> => {
    const params = date ? { date } : {};
    const response = await api.get(`/logs/job-runs/${jobRunId}/download`, {
      params,
      responseType: 'blob',
    });
    return response.data;
  },
};

// Health API
export const healthApi = {
  // Check API health
  check: async (): Promise<{ message: string }> => {
    const response = await api.get('/health');
    return response.data;
  },
};

export default api; 
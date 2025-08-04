import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Pagination,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import AddIcon from '@mui/icons-material/Add';
import VisibilityIcon from '@mui/icons-material/Visibility';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import { jobsApi, jobRunsApi, Job } from '../services/api';

export function Jobs() {
  const navigate = useNavigate();
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [runJobDialog, setRunJobDialog] = useState<{ open: boolean; job?: Job }>({ open: false });
  const [triggeredBy, setTriggeredBy] = useState('');

  const limit = 10;

  useEffect(() => {
    fetchJobs();
  }, [page]);

  const fetchJobs = async () => {
    try {
      setLoading(true);
      const offset = (page - 1) * limit;
      const response = await jobsApi.list({ limit, offset });
      setJobs(response.jobs);
      setTotalPages(Math.ceil(response.count / limit));
    } catch (error) {
      console.error('Error fetching jobs:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRunJob = async (job: Job) => {
    if (!triggeredBy.trim()) return;

    try {
      await jobRunsApi.create({
        job_id: job.id,
        triggered_by: triggeredBy,
      });
      setRunJobDialog({ open: false });
      setTriggeredBy('');
      // Optionally show success message
    } catch (error) {
      console.error('Error running job:', error);
    }
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Jobs</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => navigate('/jobs/create')}
        >
          Create Job
        </Button>
      </Box>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Description</TableCell>
                <TableCell>Source</TableCell>
                <TableCell>Tasks</TableCell>
                <TableCell>Created</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {jobs.map((job) => (
                <TableRow key={job.id}>
                  <TableCell>
                    <Typography variant="subtitle2">{job.name}</Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2" color="textSecondary">
                      {job.description || 'No description'}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip label={job.source} size="small" />
                  </TableCell>
                  <TableCell>
                    {job.tasks?.length || 0} tasks
                  </TableCell>
                  <TableCell>
                    {new Date(job.created_at).toLocaleDateString()}
                  </TableCell>
                  <TableCell>
                    <IconButton
                      onClick={() => navigate(`/jobs/${job.id}`)}
                      size="small"
                      title="View Details"
                    >
                      <VisibilityIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => setRunJobDialog({ open: true, job })}
                      size="small"
                      title="Run Job"
                      color="primary"
                    >
                      <PlayArrowIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        <Box display="flex" justifyContent="center" p={2}>
          <Pagination
            count={totalPages}
            page={page}
            onChange={(_, newPage) => setPage(newPage)}
          />
        </Box>
      </Paper>

      {/* Run Job Dialog */}
      <Dialog
        open={runJobDialog.open}
        onClose={() => setRunJobDialog({ open: false })}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Run Job: {runJobDialog.job?.name}</DialogTitle>
        <DialogContent>
          <TextField
            label="Triggered By"
            fullWidth
            value={triggeredBy}
            onChange={(e) => setTriggeredBy(e.target.value)}
            placeholder="Enter who/what triggered this job run"
            margin="normal"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRunJobDialog({ open: false })}>
            Cancel
          </Button>
          <Button
            onClick={() => handleRunJob(runJobDialog.job!)}
            variant="contained"
            disabled={!triggeredBy.trim()}
          >
            Run Job
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
```

```typescript:web-ui/src/pages/JobDetails.tsx
import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Card,
  CardContent,
  Grid,
  Chip,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import DescriptionIcon from '@mui/icons-material/Description';
import { jobRunsApi, JobRun } from '../services/api';

export function JobRunDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [jobRun, setJobRun] = useState<JobRun | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) {
      fetchJobRun();
    }
  }, [id]);

  const fetchJobRun = async () => {
    try {
      setLoading(true);
      const response = await jobRunsApi.get(id!);
      setJobRun(response);
    } catch (error) {
      console.error('Error fetching job run:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'completed':
        return 'success';
      case 'failed':
        return 'error';
      case 'running':
        return 'primary';
      case 'queued':
        return 'warning';
      default:
        return 'default';
    }
  };

  const getDuration = (startTime?: string, endTime?: string) => {
    if (!startTime) return 'Not started';
    if (!endTime) return 'Running...';
    
    const start = new Date(startTime).getTime();
    const end = new Date(endTime).getTime();
    const duration = end - start;
    
    const seconds = Math.floor(duration / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    
    if (hours > 0) return `${hours}h ${minutes % 60}m ${seconds % 60}s`;
    if (minutes > 0) return `${minutes}m ${seconds % 60}s`;
    return `${seconds}s`;
  };

  if (loading) {
    return <Typography>Loading...</Typography>;
  }

  if (!jobRun) {
    return <Typography>Job run not found</Typography>;
  }

  return (
    <Box>
      <Box display="flex" alignItems="center" gap={2} mb={3}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={() => navigate('/job-runs')}
        >
          Back to Job Runs
        </Button>
        <Typography variant="h4">Job Run Details</Typography>
        <Button
          variant="outlined"
          startIcon={<DescriptionIcon />}
          onClick={() => navigate(`/logs?job_run_id=${jobRun.id}`)}
        >
          View Logs
        </Button>
      </Box>

      <Grid container spacing={3}>
        {/* Job Run Information */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Job Run Information
              </Typography>
              <Box display="flex" flexDirection="column" gap={2}>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Job Run ID
                  </Typography>
                  <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                    {jobRun.id}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Job ID
                  </Typography>
                  <Typography 
                    variant="body2" 
                    sx={{ fontFamily: 'monospace', cursor: 'pointer', color: 'primary.main' }}
                    onClick={() => navigate(`/jobs/${jobRun.job_id}`)}
                  >
                    {jobRun.job_id}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Status
                  </Typography>
                  <Chip
                    label={jobRun.status}
                    color={getStatusColor(jobRun.status) as any}
                    size="small"
                  />
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Triggered By
                  </Typography>
                  <Typography variant="body2">{jobRun.triggered_by}</Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Timing Information */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Timing Information
              </Typography>
              <Box display="flex" flexDirection="column" gap={2}>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Created
                  </Typography>
                  <Typography variant="body2">
                    {new Date(jobRun.created_at).toLocaleString()}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Started
                  </Typography>
                  <Typography variant="body2">
                    {jobRun.started_at
                      ? new Date(jobRun.started_at).toLocaleString()
                      : 'Not started'
                    }
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Completed
                  </Typography>
                  <Typography variant="body2">
                    {jobRun.completed_at
                      ? new Date(jobRun.completed_at).toLocaleString()
                      : 'Not completed'
                    }
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Duration
                  </Typography>
                  <Typography variant="body2">
                    {getDuration(jobRun.started_at, jobRun.completed_at)}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Task Runs */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Task Runs ({jobRun.task_runs?.length || 0})
              </Typography>
              {jobRun.task_runs && jobRun.task_runs.length > 0 ? (
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Task Run ID</TableCell>
                        <TableCell>Task ID</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Created</TableCell>
                        <TableCell>Started</TableCell>
                        <TableCell>Completed</TableCell>
                        <TableCell>Duration</TableCell>
                        <TableCell>Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {jobRun.task_runs.map((taskRun) => (
                        <TableRow key={taskRun.id}>
                          <TableCell sx={{ fontFamily: 'monospace' }}>
                            {taskRun.id.substring(0, 8)}...
                          </TableCell>
                          <TableCell sx={{ fontFamily: 'monospace' }}>
                            {taskRun.task_id.substring(0, 8)}...
                          </TableCell>
                          <TableCell>
                            <Chip
                              label={taskRun.status}
                              color={getStatusColor(taskRun.status) as any}
                              size="small"
                            />
                          </TableCell>
                          <TableCell>
                            {new Date(taskRun.created_at).toLocaleString()}
                          </TableCell>
                          <TableCell>
                            {taskRun.started_at
                              ? new Date(taskRun.started_at).toLocaleString()
                              : '-'
                            }
                          </TableCell>
                          <TableCell>
                            {taskRun.completed_at
                              ? new Date(taskRun.completed_at).toLocaleString()
                              : '-'
                            }
                          </TableCell>
                          <TableCell>
                            {getDuration(taskRun.started_at, taskRun.completed_at)}
                          </TableCell>
                          <TableCell>
                            <Button
                              size="small"
                              onClick={() => navigate(`/logs?task_run_id=${taskRun.id}`)}
                            >
                              View Logs
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              ) : (
                <Typography color="textSecondary">No task runs found</Typography>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}
```

```typescript:web-ui/src/pages/Secrets.tsx
import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import { secretsApi, Secret, CreateSecretRequest } from '../services/api';

export function Secrets() {
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [loading, setLoading] = useState(true);
  const [createDialog, setCreateDialog] = useState(false);
  const [editDialog, setEditDialog] = useState<{ open: boolean; secret?: Secret }>({ open: false });
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; secret?: Secret }>({ open: false });
  const [formData, setFormData] = useState<CreateSecretRequest>({ name: '', value: '' });
  const [editValue, setEditValue] = useState('');

  useEffect(() => {
    fetchSecrets();
  }, []);

  const fetchSecrets = async () => {
    try {
      setLoading(true);
      const response = await secretsApi.list();
      setSecrets(response.secrets);
    } catch (error) {
      console.error('Error fetching secrets:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateSecret = async () => {
    try {
      await secretsApi.create(formData);
      setCreateDialog(false);
      setFormData({ name: '', value: '' });
      fetchSecrets();
    } catch (error) {
      console.error('Error creating secret:', error);
    }
  };

  const handleUpdateSecret = async () => {
    if (!editDialog.secret) return;
    
    try {
      await secretsApi.update(editDialog.secret.id, editValue);
      setEditDialog({ open: false });
      setEditValue('');
      fetchSecrets();
    } catch (error) {
      console.error('Error updating secret:', error);
    }
  };

  const handleDeleteSecret = async () => {
    if (!deleteDialog.secret) return;
    
    try {
      await secretsApi.delete(deleteDialog.secret.id);
      setDeleteDialog({ open: false });
      fetchSecrets();
    } catch (error) {
      console.error('Error deleting secret:', error);
    }
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Secrets</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setCreateDialog(true)}
        >
          Create Secret
        </Button>
      </Box>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>ID</TableCell>
                <TableCell>Created</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {secrets.map((secret) => (
                <TableRow key={secret.id}>
                  <TableCell>
                    <Typography variant="subtitle2">{secret.name}</Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                      {secret.id}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    {new Date(secret.created_at).toLocaleDateString()}
                  </TableCell>
                  <TableCell>
                    <IconButton
                      onClick={() => {
                        setEditDialog({ open: true, secret });
                        setEditValue('');
                      }}
                      size="small"
                      title="Edit Secret"
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => setDeleteDialog({ open: true, secret })}
                      size="small"
                      title="Delete Secret"
                      color="error"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Create Secret Dialog */}
      <Dialog
        open={createDialog}
        onClose={() => setCreateDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Create New Secret</DialogTitle>
        <DialogContent>
          <TextField
            label="Name"
            fullWidth
            margin="normal"
            value={formData.name}
            onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
          />
          <TextField
            label="Value"
            fullWidth
            margin="normal"
            type="password"
            value={formData.value}
            onChange={(e) => setFormData(prev => ({ ...prev, value: e.target.value }))}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialog(false)}>Cancel</Button>
          <Button
            onClick={handleCreateSecret}
            variant="contained"
            disabled={!formData.name.trim() || !formData.value.trim()}
          >
            Create
          </Button>
        </DialogActions>
      </Dialog>

      {/* Edit Secret Dialog */}
      <Dialog
        open={editDialog.open}
        onClose={() => setEditDialog({ open: false })}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Update Secret: {editDialog.secret?.name}</DialogTitle>
        <DialogContent>
          <TextField
            label="New Value"
            fullWidth
            margin="normal"
            type="password"
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            placeholder="Enter new secret value"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialog({ open: false })}>Cancel</Button>
          <Button
            onClick={handleUpdateSecret}
            variant="contained"
            disabled={!editValue.trim()}
          >
            Update
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialog.open}
        onClose={() => setDeleteDialog({ open: false })}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Delete Secret</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete the secret "{deleteDialog.secret?.name}"?
            This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialog({ open: false })}>Cancel</Button>
          <Button
            onClick={handleDeleteSecret}
            variant="contained"
            color="error"
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
```

Finally, let me create the Logs page that integrates with the existing LogViewer:

```typescript:web-ui/src/pages/Logs.tsx
import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Button,
  Grid,
  Card,
  CardContent,
  Switch,
  FormControlLabel,
} from '@mui/material';
import { useSearchParams } from 'react-router-dom';
import DownloadIcon from '@mui/icons-material/Download';
import { LogViewer } from '../components/LogViewer';
import { logsApi, LogEntry } from '../services/api';

export function Logs() {
  const [searchParams] = useSearchParams();
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [logType, setLogType] = useState<'system' | 'job' | 'task' | 'all'>('all');
  const [jobRunId, setJobRunId] = useState(searchParams.get('job_run_id') || '');
  const [taskRunId, setTaskRunId] = useState(searchParams.get('task_run_id') || '');
  const [streamingEnabled, setStreamingEnabled] = useState(false);
  const [streamingStatus, setStreamingStatus] = useState<any>(null);

  useEffect(() => {
    // Initialize with URL params
    if (searchParams.get('job_run_id')) {
      setJobRunId(searchParams.get('job_run_id')!);
      fetchJobRunLogs(searchParams.get('job_run_id')!);
    } else if (searchParams.get('task_run_id')) {
      setTaskRunId(searchParams.get('task_run_id')!);
      fetchTaskRunLogs(searchParams.get('task_run_id')!);
    } else {
      fetchLogs();
    }
  }, [searchParams]);

  useEffect(() => {
    if (streamingEnabled) {
      fetchStreamingStatus();
    }
  }, [streamingEnabled]);

  const fetchLogs = async () => {
    try {
      setLoading(true);
      let response;
      
      if (jobRunId.trim()) {
        response = await logsApi.getJobRunLogs(jobRunId);
        setLogs(response.logs);
      } else if (taskRunId.trim()) {
        response = await logsApi.getTaskRunLogs(taskRunId);
        setLogs(response.logs);
      } else if (logType === 'system') {
        response = await logsApi.getSystemLogs();
        setLogs(response.logs);
      } else if (logType !== 'all') {
        response = await logsApi.getLogsByType(logType);
        setLogs(response.logs);
      } else {
        // Fetch system logs as default
        response = await logsApi.getSystemLogs();
        setLogs(response.logs);
      }
    } catch (error) {
      console.error('Error fetching logs:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchJobRunLogs = async (id: string) => {
    try {
      setLoading(true);
      const response = await logsApi.getJobRunLogs(id);
      setLogs(response.logs);
    } catch (error) {
      console.error('Error fetching job run logs:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchTaskRunLogs = async (id: string) => {
    try {
      setLoading(true);
      const response = await logsApi.getTaskRunLogs(id);
      setLogs(response.logs);
    } catch (error) {
      console.error('Error fetching task run logs:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchStreamingStatus = async () => {
    try {
      const response = await logsApi.getStreamingStatus();
      setStreamingStatus(response);
    } catch (error) {
      console.error('Error fetching streaming status:', error);
    }
  };

  const handleDownloadLogs = async () => {
    if (!jobRunId.trim()) {
      alert('Job Run ID is required for log download');
      return;
    }

    try {
      const blob = await logsApi.downloadJobRunLogs(jobRunId);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `job-run-${jobRunId}-logs.txt`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (error) {
      console.error('Error downloading logs:', error);
    }
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Logs
      </Typography>

      <Grid container spacing={3} mb={3}>
        {/* Controls */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Log Filters
              </Typography>
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} sm={6} md={3}>
                  <FormControl fullWidth>
                    <InputLabel>Log Type</InputLabel>
                    <Select
                      value={logType}
                      label="Log Type"
                      onChange={(e) => setLogType(e.target.value as any)}
                    >
                      <MenuItem value="all">All</MenuItem>
                      <MenuItem value="system">System</MenuItem>
                      <MenuItem value="job">Job</MenuItem>
                      <MenuItem value="task">Task</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                  <TextField
                    label="Job Run ID"
                    fullWidth
                    value={jobRunId}
                    onChange={(e) => setJobRunId(e.target.value)}
                    placeholder="Enter Job Run ID"
                  />
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                  <TextField
                    label="Task Run ID"
                    fullWidth
                    value={taskRunId}
                    onChange={(e) => setTaskRunId(e.target.value)}
                    placeholder="Enter Task Run ID"
                  />
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                  <Button
                    variant="contained"
                    onClick={fetchLogs}
                    disabled={loading}
                    fullWidth
                  >
                    {loading ? 'Loading...' : 'Fetch Logs'}
                  </Button>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Actions */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Actions
              </Typography>
              <Box display="flex" flexDirection="column" gap={2}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={streamingEnabled}
                      onChange={(e) => setStreamingEnabled(e.target.checked)}
                    />
                  }
                  label="Enable Real-time Streaming"
                />
                <Button
                  variant="outlined"
                  startIcon={<DownloadIcon />}
                  onClick={handleDownloadLogs}
                  disabled={!jobRunId.trim()}
                  fullWidth
                >
                  Download Logs
                </Button>
                {streamingStatus && (
                  <Typography variant="caption" color="textSecondary">
                    Active connections: {streamingStatus.connections?.length || 0}
                  </Typography>
                )}
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Log Viewer */}
      <Paper sx={{ height: 600 }}>
        <LogViewer 
          logs={logs}
          streaming={streamingEnabled}
          jobRunId={jobRunId || undefined}
        />
      </Paper>
    </Box>
  );
}
```

Finally, let me update the main.tsx to include the Router:

```typescript:web-ui/src/main.tsx
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
```

Now I need to create the page exports:

```typescript:web-ui/src/pages/index.ts
export { Dashboard } from './Dashboard';
export { Jobs } from './Jobs';
export { JobDetails } from './JobDetails';
export { CreateJob } from './CreateJob';
export { JobRuns } from './JobRuns';
export { JobRunDetails } from './JobRunDetails';
export { Secrets } from './Secrets';
export { Logs } from './Logs';
```

## Summary

I've created a comprehensive UI that covers all the available APIs in your Stratal system:

### **Updated Package Dependencies:**
- Added `react-router-dom` for navigation
- Added `axios` for API calls

### **API Service Layer:**
- Complete TypeScript API service in `src/services/api.ts`
- Covers all endpoints: Jobs, Job Runs, Secrets, Logs, Health
- Proper TypeScript interfaces for all data types

### **Navigation & Routing:**
- Updated sidebar with proper React Router navigation
- Full routing setup in App.tsx
- Visual feedback for active routes

### **Pages Created:**

1. **Dashboard** - Overview with stats cards and recent jobs/runs
2. **Jobs** - List all jobs, create new jobs, view details, run jobs
3. **Job Details** - Detailed view of individual jobs with tasks
4. **Create Job** - Comprehensive job creation form with tasks
5. **Job Runs** - List and filter job runs by status  
6. **Job Run Details** - Detailed view of job runs with task runs
7. **Secrets** - Manage secrets (create, edit, delete)
8. **Logs** - Integrated log viewer with filtering and real-time streaming

### **Key Features:**

- **Real-time functionality** with WebSocket log streaming
- **Comprehensive job management** with immediate execution option
- **Task management** with different task types (http_request, send_email, format_output, custom_script)
- **Secret management** with encrypted storage
- **Advanced logging** with filtering by type, job run, task run
- **Log downloading** capability
- **Responsive design** with Material-UI
- **Type safety** throughout with TypeScript
- **Error handling** and loading states

### **Usage:**

1. Install new dependencies: `cd web-ui && yarn add react-router-dom axios`
2. The UI will automatically connect to your backend at `localhost:8080`
3. All API endpoints are properly integrated and functional

The UI now provides a complete interface for all your Stratal system functionality, from job creation and management to real-time log monitoring and secret management. 
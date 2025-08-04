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

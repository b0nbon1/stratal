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
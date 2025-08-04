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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import { jobsApi, jobRunsApi, Job } from '../services/api';

export function JobDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [job, setJob] = useState<Job | null>(null);
  const [loading, setLoading] = useState(true);
  const [runJobDialog, setRunJobDialog] = useState(false);
  const [triggeredBy, setTriggeredBy] = useState('');

  useEffect(() => {
    if (id) {
      fetchJob();
    }
  }, [id]);

  const fetchJob = async () => {
    try {
      setLoading(true);
      const response = await jobsApi.get(id!);
      setJob(response);
    } catch (error) {
      console.error('Error fetching job:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRunJob = async () => {
    if (!job || !triggeredBy.trim()) return;

    try {
      const response = await jobRunsApi.create({
        job_id: job.id,
        triggered_by: triggeredBy,
      });
      setRunJobDialog(false);
      setTriggeredBy('');
      navigate(`/job-runs/${response.JobRunID}`);
    } catch (error) {
      console.error('Error running job:', error);
    }
  };

  if (loading) {
    return <Typography>Loading...</Typography>;
  }

  if (!job) {
    return <Typography>Job not found</Typography>;
  }

  return (
    <Box>
      <Box display="flex" alignItems="center" gap={2} mb={3}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={() => navigate('/jobs')}
        >
          Back to Jobs
        </Button>
        <Typography variant="h4">{job.name}</Typography>
        <Button
          variant="contained"
          startIcon={<PlayArrowIcon />}
          onClick={() => setRunJobDialog(true)}
        >
          Run Job
        </Button>
      </Box>

      <Grid container spacing={3}>
        {/* Job Information */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Job Information
              </Typography>
              <Box display="flex" flexDirection="column" gap={2}>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    ID
                  </Typography>
                  <Typography variant="body2">{job.id}</Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Description
                  </Typography>
                  <Typography variant="body2">
                    {job.description || 'No description provided'}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Source
                  </Typography>
                  <Chip label={job.source} size="small" />
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Created
                  </Typography>
                  <Typography variant="body2">
                    {new Date(job.created_at).toLocaleString()}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="textSecondary">
                    Last Updated
                  </Typography>
                  <Typography variant="body2">
                    {new Date(job.updated_at).toLocaleString()}
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Raw Payload */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Raw Payload
              </Typography>
              <Paper
                sx={{
                  p: 2,
                  backgroundColor: '#f5f5f5',
                  maxHeight: 300,
                  overflow: 'auto',
                }}
              >
                <Typography variant="body2" component="pre">
                  {job.raw_payload
                    ? JSON.stringify(job.raw_payload, null, 2)
                    : 'No payload data'}
                </Typography>
              </Paper>
            </CardContent>
          </Card>
        </Grid>

        {/* Tasks */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Tasks ({job.tasks?.length || 0})
              </Typography>
              {job.tasks && job.tasks.length > 0 ? (
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Order</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell>Type</TableCell>
                        <TableCell>Config</TableCell>
                        <TableCell>Created</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {job.tasks
                        .sort((a, b) => a.order - b.order)
                        .map((task) => (
                          <TableRow key={task.id}>
                            <TableCell>{task.order}</TableCell>
                            <TableCell>{task.name}</TableCell>
                            <TableCell>
                              <Chip label={task.type} size="small" />
                            </TableCell>
                            <TableCell>
                              <Typography variant="body2" sx={{ maxWidth: 200 }}>
                                {JSON.stringify(task.config)}
                              </Typography>
                            </TableCell>
                            <TableCell>
                              {new Date(task.created_at).toLocaleDateString()}
                            </TableCell>
                          </TableRow>
                        ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              ) : (
                <Typography color="textSecondary">No tasks defined</Typography>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Run Job Dialog */}
      <Dialog
        open={runJobDialog}
        onClose={() => setRunJobDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Run Job: {job.name}</DialogTitle>
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
          <Button onClick={() => setRunJobDialog(false)}>Cancel</Button>
          <Button
            onClick={handleRunJob}
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
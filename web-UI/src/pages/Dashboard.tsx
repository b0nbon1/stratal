import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  CardHeader,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Button,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';
import EventNoteIcon from '@mui/icons-material/EventNote';
import { jobsApi, jobRunsApi, type Job, type JobRun } from '../services/api';

interface DashboardStats {
  totalJobs: number;
  totalJobRuns: number;
  runningJobs: number;
  completedJobs: number;
  failedJobs: number;
  scheduledJobs: number;
}

export function Dashboard() {
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats>({
    totalJobs: 0,
    totalJobRuns: 0,
    runningJobs: 0,
    completedJobs: 0,
    failedJobs: 0,
    scheduledJobs: 0,
  });
  const [recentJobs, setRecentJobs] = useState<Job[]>([]);
  const [recentJobRuns, setRecentJobRuns] = useState<JobRun[]>([]);
  const [, setLoading] = useState(true);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      
      // Fetch recent jobs
      const jobsResponse = await jobsApi.list({ limit: 5 });
      setRecentJobs(jobsResponse.jobs);

      // Fetch recent job runs
      const jobRunsResponse = await jobRunsApi.list({ limit: 10 });
      setRecentJobRuns(jobRunsResponse.job_runs || []);

      // Calculate stats
      const totalJobs = jobsResponse.count;
      const totalJobRuns = jobRunsResponse.count || 0;
      
      // Count job runs by status
      const runningJobs = jobRunsResponse.job_runs?.filter(jr => jr.status === 'running' || jr.status === 'queued').length || 0;
      const completedJobs = jobRunsResponse.job_runs?.filter(jr => jr.status === 'completed').length || 0;
      const failedJobs = jobRunsResponse.job_runs?.filter(jr => jr.status === 'failed').length || 0;
      const scheduledJobs = jobRunsResponse.job_runs?.filter(jr => jr.status === 'scheduled').length || 0;

      setStats({
        totalJobs,
        totalJobRuns,
        runningJobs,
        completedJobs,
        failedJobs,
        scheduledJobs,
      });
    } catch (error) {
      console.error('Error fetching dashboard data:', error);
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

  const statsCards = [
    {
      title: 'Running/Queued',
      value: stats.runningJobs.toString(),
      icon: <AccessTimeIcon fontSize="small" sx={{ color: '#2563eb' }} />,
      bgColor: '#eff6ff',
      change: 'Active jobs',
    },
    {
      title: 'Completed',
      value: stats.completedJobs.toString(),
      icon: <CheckCircleIcon fontSize="small" sx={{ color: '#16a34a' }} />,
      bgColor: '#f0fdf4',
      change: 'Successful runs',
    },
    {
      title: 'Failed',
      value: stats.failedJobs.toString(),
      icon: <CancelIcon fontSize="small" sx={{ color: '#dc2626' }} />,
      bgColor: '#fef2f2',
      change: 'Failed runs',
    },
    {
      title: 'Total Jobs',
      value: stats.totalJobs.toString(),
      icon: <EventNoteIcon fontSize="small" sx={{ color: '#7e22ce' }} />,
      bgColor: '#faf5ff',
      change: 'Total jobs created',
    },
  ];

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      {/* Stats Cards */}
      <Grid container spacing={3} mb={4}>
        {statsCards.map((stat) => (
          <Grid item xs={12} md={6} lg={3} key={stat.title}>
            <Card
              elevation={1}
              sx={{
                transition: 'transform 0.2s',
                '&:hover': { transform: 'scale(1.02)' },
              }}
            >
              <CardHeader
                sx={{ pb: 1 }}
                title={
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Typography variant="subtitle2" color="text.secondary">
                      {stat.title}
                    </Typography>
                    <Box
                      sx={{
                        backgroundColor: stat.bgColor,
                        p: 1,
                        borderRadius: 2,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                      }}
                    >
                      {stat.icon}
                    </Box>
                  </Box>
                }
              />
              <CardContent>
                <Typography variant="h5" fontWeight="bold">
                  {stat.value}
                </Typography>
                <Typography variant="caption" color="text.secondary" mt={1}>
                  {stat.change}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Grid container spacing={3}>
        {/* Recent Jobs */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader
              title="Recent Jobs"
              action={
                <Button onClick={() => navigate('/jobs')} size="small">
                  View All
                </Button>
              }
            />
            <CardContent>
              <TableContainer>
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Source</TableCell>
                      <TableCell>Created</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {recentJobs.map((job) => (
                      <TableRow
                        key={job.id}
                        sx={{ cursor: 'pointer' }}
                        onClick={() => navigate(`/jobs/${job.id}`)}
                      >
                        <TableCell>{job.name}</TableCell>
                        <TableCell>{job.source}</TableCell>
                        <TableCell>
                          {new Date(job.created_at).toLocaleDateString()}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Recent Job Runs */}
        <Grid item xs={12} md={6}>
          <Card>
            <CardHeader
              title="Recent Job Runs"
              action={
                <Button onClick={() => navigate('/job-runs')} size="small">
                  View All
                </Button>
              }
            />
            <CardContent>
              <TableContainer>
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Job Run ID</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Started</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {recentJobRuns.slice(0, 5).map((jobRun) => (
                      <TableRow
                        key={jobRun.id}
                        sx={{ cursor: 'pointer' }}
                        onClick={() => navigate(`/job-runs/${jobRun.id}`)}
                      >
                        <TableCell>
                          {jobRun.id.substring(0, 8)}...
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={jobRun.status}
                            color={getStatusColor(jobRun.status) as any}
                            size="small"
                          />
                        </TableCell>
                        <TableCell>
                          {new Date(jobRun.created_at).toLocaleDateString()}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
} 
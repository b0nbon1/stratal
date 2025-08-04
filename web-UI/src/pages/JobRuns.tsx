import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  Pagination,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import VisibilityIcon from '@mui/icons-material/Visibility';
import DescriptionIcon from '@mui/icons-material/Description';
import { jobRunsApi, JobRun } from '../services/api';

export function JobRuns() {
  const navigate = useNavigate();
  const [jobRuns, setJobRuns] = useState<JobRun[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [statusFilter, setStatusFilter] = useState('all');

  const limit = 10;

  useEffect(() => {
    fetchJobRuns();
  }, [page]);

  const fetchJobRuns = async () => {
    try {
      setLoading(true);
      const offset = (page - 1) * limit;
      const response = await jobRunsApi.list({ limit, offset });
      setJobRuns(response.job_runs || []);
      setTotalPages(Math.ceil((response.count || 0) / limit));
    } catch (error) {
      console.error('Error fetching job runs:', error);
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

  const filteredJobRuns = statusFilter === 'all' 
    ? jobRuns 
    : jobRuns.filter(jr => jr.status.toLowerCase() === statusFilter);

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Job Runs</Typography>
        <FormControl sx={{ minWidth: 150 }}>
          <InputLabel>Status Filter</InputLabel>
          <Select
            value={statusFilter}
            label="Status Filter"
            onChange={(e) => setStatusFilter(e.target.value)}
          >
            <MenuItem value="all">All</MenuItem>
            <MenuItem value="queued">Queued</MenuItem>
            <MenuItem value="running">Running</MenuItem>
            <MenuItem value="completed">Completed</MenuItem>
            <MenuItem value="failed">Failed</MenuItem>
          </Select>
        </FormControl>
      </Box>

      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Job Run ID</TableCell>
                <TableCell>Job ID</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Triggered By</TableCell>
                <TableCell>Created</TableCell>
                <TableCell>Started</TableCell>
                <TableCell>Completed</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {filteredJobRuns.map((jobRun) => (
                <TableRow key={jobRun.id}>
                  <TableCell>
                    <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                      {jobRun.id.substring(0, 8)}...
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                      {jobRun.job_id.substring(0, 8)}...
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={jobRun.status}
                      color={getStatusColor(jobRun.status) as any}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>{jobRun.triggered_by}</TableCell>
                  <TableCell>
                    {new Date(jobRun.created_at).toLocaleString()}
                  </TableCell>
                  <TableCell>
                    {jobRun.started_at 
                      ? new Date(jobRun.started_at).toLocaleString()
                      : '-'
                    }
                  </TableCell>
                  <TableCell>
                    {jobRun.completed_at 
                      ? new Date(jobRun.completed_at).toLocaleString()
                      : '-'
                    }
                  </TableCell>
                  <TableCell>
                    <IconButton
                      onClick={() => navigate(`/job-runs/${jobRun.id}`)}
                      size="small"
                      title="View Details"
                    >
                      <VisibilityIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => navigate(`/logs?job_run_id=${jobRun.id}`)}
                      size="small"
                      title="View Logs"
                      color="primary"
                    >
                      <DescriptionIcon />
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
    </Box>
  );
} 
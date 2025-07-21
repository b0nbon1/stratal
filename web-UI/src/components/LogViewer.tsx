import React, { useState, useEffect, useRef, useCallback } from 'react';
import { 
  Box, 
  Paper, 
  Typography, 
  Button, 
  Select, 
  MenuItem, 
  FormControl, 
  InputLabel, 
  TextField, 
  Chip,
  IconButton,
  Tooltip,
  Switch,
  FormControlLabel
} from '@mui/material';
import { 
  PlayArrow, 
  Stop, 
  Download, 
  Clear, 
  FilterList,
  Search,
  AutoFixHigh
} from '@mui/icons-material';
import { LogStreamingService, LogMessage, WebSocketLogStreamer, SSELogStreamer } from '../utils/LogStreamer';

interface LogViewerProps {
  jobRunId?: string;
  height?: string;
  autoScroll?: boolean;
}

const LogViewer: React.FC<LogViewerProps> = ({ 
  jobRunId, 
  height = '500px', 
  autoScroll = true 
}) => {
  const [logs, setLogs] = useState<LogMessage[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [streamingStatus, setStreamingStatus] = useState<string>('disconnected');
  const [filterLevel, setFilterLevel] = useState<string>('all');
  const [filterStream, setFilterStream] = useState<string>('all');
  const [filterType, setFilterType] = useState<string>('all');
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [useWebSocket, setUseWebSocket] = useState<boolean>(true);
  const [maxLogs, setMaxLogs] = useState<number>(1000);
  
  const logContainerRef = useRef<HTMLDivElement>(null);
  const streamingServiceRef = useRef<LogStreamingService>(new LogStreamingService());
  const streamerRef = useRef<WebSocketLogStreamer | SSELogStreamer | null>(null);

  // Auto-scroll to bottom when new logs arrive
  useEffect(() => {
    if (autoScroll && logContainerRef.current) {
      logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight;
    }
  }, [logs, autoScroll]);

  // Load historical logs on component mount
  useEffect(() => {
    if (jobRunId) {
      loadHistoricalLogs();
    }
  }, [jobRunId]);

  const loadHistoricalLogs = async () => {
    if (!jobRunId) return;
    
    try {
      const response = await streamingServiceRef.current.getJobRunLogs(jobRunId, 100, 0);
      const historicalLogs = response.logs || [];
      setLogs(historicalLogs);
    } catch (error) {
      console.error('Failed to load historical logs:', error);
    }
  };

  const handleMessage = useCallback((message: LogMessage) => {
    setLogs(prevLogs => {
      const newLogs = [...prevLogs, message];
      // Keep only the last maxLogs entries to prevent memory issues
      if (newLogs.length > maxLogs) {
        return newLogs.slice(-maxLogs);
      }
      return newLogs;
    });
  }, [maxLogs]);

  const handleError = useCallback((error: Error) => {
    console.error('Log streaming error:', error);
    setStreamingStatus('error');
  }, []);

  const handleConnect = useCallback(() => {
    setStreamingStatus('connected');
  }, []);

  const handleDisconnect = useCallback(() => {
    setStreamingStatus('disconnected');
  }, []);

  const startStreaming = () => {
    const options = {
      onMessage: handleMessage,
      onError: handleError,
      onConnect: handleConnect,
      onDisconnect: handleDisconnect,
      autoReconnect: true,
      maxReconnectAttempts: 5,
    };

    if (jobRunId) {
      streamerRef.current = streamingServiceRef.current.streamJobRunLogs(
        jobRunId, 
        options, 
        useWebSocket
      );
    } else {
      streamerRef.current = streamingServiceRef.current.streamAllLogs(
        options, 
        useWebSocket
      );
    }

    setIsStreaming(true);
  };

  const stopStreaming = () => {
    if (streamerRef.current) {
      streamerRef.current.disconnect();
      streamerRef.current = null;
    }
    setIsStreaming(false);
    setStreamingStatus('disconnected');
  };

  const clearLogs = () => {
    setLogs([]);
  };

  const downloadLogs = () => {
    if (jobRunId) {
      streamingServiceRef.current.downloadJobRunLogs(jobRunId);
    } else {
      // For global logs, we can export the current view
      const logText = filteredLogs.map(log => 
        `[${log.timestamp}] [${log.level}] [${log.stream}] ${log.message}`
      ).join('\n');
      
      const blob = new Blob([logText], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `logs-${new Date().toISOString().slice(0, 10)}.txt`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    }
  };

  // Filter logs based on level, stream, type, and search term
  const filteredLogs = logs.filter(log => {
    const levelMatch = filterLevel === 'all' || log.level === filterLevel;
    const streamMatch = filterStream === 'all' || log.stream === filterStream;
    const typeMatch = filterType === 'all' || log.type === filterType;
    const searchMatch = searchTerm === '' || 
      log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
      log.job_run_id.toLowerCase().includes(searchTerm.toLowerCase());
    
    return levelMatch && streamMatch && typeMatch && searchMatch;
  });

  const getLogLevelColor = (level: string) => {
    switch (level) {
      case 'error': return '#f44336';
      case 'warn': return '#ff9800';
      case 'info': return '#2196f3';
      case 'debug': return '#9e9e9e';
      default: return '#000000';
    }
  };

  const getStreamColor = (stream: string) => {
    switch (stream) {
      case 'stderr': return '#f44336';
      case 'stdout': return '#4caf50';
      case 'system': return '#9c27b0';
      default: return '#757575';
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'system': return '#ff9800';
      case 'job': return '#2196f3';
      case 'task': return '#4caf50';
      default: return '#757575';
    }
  };

  const getConnectionStatusColor = () => {
    switch (streamingStatus) {
      case 'connected': return 'success';
      case 'error': return 'error';
      default: return 'default';
    }
  };

  return (
    <Paper elevation={2} sx={{ p: 2 }}>
      {/* Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Typography variant="h6">
          Log Viewer {jobRunId && `- Job Run: ${jobRunId.slice(0, 8)}...`}
        </Typography>
        
        <Box display="flex" alignItems="center" gap={1}>
          <Chip 
            label={streamingStatus} 
            color={getConnectionStatusColor()}
            size="small" 
          />
          
          <FormControlLabel
            control={
              <Switch 
                checked={useWebSocket} 
                onChange={(e) => setUseWebSocket(e.target.checked)}
                size="small"
              />
            }
            label="WebSocket"
            sx={{ mr: 1 }}
          />
          
          {!isStreaming ? (
            <Button
              startIcon={<PlayArrow />}
              onClick={startStreaming}
              variant="contained"
              size="small"
              color="primary"
            >
              Start Stream
            </Button>
          ) : (
            <Button
              startIcon={<Stop />}
              onClick={stopStreaming}
              variant="contained"
              size="small"
              color="secondary"
            >
              Stop Stream
            </Button>
          )}
        </Box>
      </Box>

      {/* Controls */}
      <Box display="flex" gap={2} mb={2} flexWrap="wrap" alignItems="center">
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel>Level</InputLabel>
          <Select
            value={filterLevel}
            onChange={(e) => setFilterLevel(e.target.value)}
            label="Level"
          >
            <MenuItem value="all">All Levels</MenuItem>
            <MenuItem value="error">Error</MenuItem>
            <MenuItem value="warn">Warning</MenuItem>
            <MenuItem value="info">Info</MenuItem>
            <MenuItem value="debug">Debug</MenuItem>
          </Select>
        </FormControl>

        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel>Stream</InputLabel>
          <Select
            value={filterStream}
            onChange={(e) => setFilterStream(e.target.value)}
            label="Stream"
          >
            <MenuItem value="all">All Streams</MenuItem>
            <MenuItem value="stdout">stdout</MenuItem>
            <MenuItem value="stderr">stderr</MenuItem>
            <MenuItem value="system">system</MenuItem>
          </Select>
        </FormControl>

        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel>Type</InputLabel>
          <Select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            label="Type"
          >
            <MenuItem value="all">All Types</MenuItem>
            <MenuItem value="system">System</MenuItem>
            <MenuItem value="job">Job</MenuItem>
            <MenuItem value="task">Task</MenuItem>
          </Select>
        </FormControl>

        <TextField
          size="small"
          placeholder="Search logs..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          InputProps={{
            startAdornment: <Search sx={{ mr: 1, color: 'action.active' }} />,
          }}
          sx={{ minWidth: 200 }}
        />

        <Tooltip title="Clear logs">
          <IconButton onClick={clearLogs} size="small">
            <Clear />
          </IconButton>
        </Tooltip>

        <Tooltip title="Download logs">
          <IconButton onClick={downloadLogs} size="small">
            <Download />
          </IconButton>
        </Tooltip>

        <Typography variant="caption" color="textSecondary">
          {filteredLogs.length} / {logs.length} logs
        </Typography>
      </Box>

      {/* Log Display */}
      <Paper 
        variant="outlined" 
        sx={{ 
          height, 
          overflow: 'auto', 
          p: 1, 
          backgroundColor: '#1e1e1e',
          fontFamily: 'Monaco, Consolas, "Courier New", monospace'
        }}
        ref={logContainerRef}
      >
        {filteredLogs.length === 0 ? (
          <Box 
            display="flex" 
            alignItems="center" 
            justifyContent="center" 
            height="100%"
            color="textSecondary"
          >
            <Typography variant="body2">
              {logs.length === 0 ? 'No logs available' : 'No logs match the current filters'}
            </Typography>
          </Box>
        ) : (
          filteredLogs.map((log, index) => (
            <Box 
              key={index} 
              sx={{ 
                fontSize: '0.8rem', 
                lineHeight: 1.4, 
                mb: 0.5,
                color: '#ffffff',
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-all'
              }}
            >
              <span style={{ color: '#888888' }}>
                [{new Date(log.timestamp).toLocaleTimeString()}]
              </span>
              {' '}
              <span style={{ 
                color: getLogLevelColor(log.level), 
                fontWeight: 'bold',
                textTransform: 'uppercase'
              }}>
                [{log.level}]
              </span>
              {' '}
              <span style={{ color: getStreamColor(log.stream) }}>
                [{log.stream}]
              </span>
              {log.task_name && (
                <>
                  {' '}
                  <span style={{ color: '#ffeb3b' }}>
                    [{log.task_name}]
                  </span>
                </>
              )}
              {' '}
              <span style={{ color: getTypeColor(log.type) }}>
                [{log.type}]
              </span>
              {' '}
              <span>{log.message}</span>
            </Box>
          ))
        )}
      </Paper>
    </Paper>
  );
};

export default LogViewer; 
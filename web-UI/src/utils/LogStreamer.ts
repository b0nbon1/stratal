export interface LogMessage {
  type: 'system' | 'job' | 'task';
  job_run_id: string;
  task_run_id?: string;
  timestamp: string;
  level: 'info' | 'error' | 'warn' | 'debug';
  stream: 'stdout' | 'stderr' | 'system';
  message: string;
  task_name?: string;
  metadata?: Record<string, any>;
}

export interface LogStreamOptions {
  jobRunId?: string;
  onMessage: (message: LogMessage) => void;
  onError?: (error: Error) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  autoReconnect?: boolean;
  maxReconnectAttempts?: number;
}

export class WebSocketLogStreamer {
  private ws: WebSocket | null = null;
  private options: LogStreamOptions;
  private reconnectAttempts = 0;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private url: string;

  constructor(baseUrl: string, options: LogStreamOptions) {
    this.options = {
      autoReconnect: true,
      maxReconnectAttempts: 5,
      ...options,
    };
    
    const params = new URLSearchParams();
    if (this.options.jobRunId) {
      params.append('job_run_id', this.options.jobRunId);
    }
    
    this.url = `${baseUrl}/api/v1/logs/stream/ws?${params.toString()}`;
  }

  connect(): void {
    try {
      this.ws = new WebSocket(this.url);
      
      this.ws.onopen = () => {
        console.log('WebSocket connected to log stream');
        this.reconnectAttempts = 0;
        this.options.onConnect?.();
      };

      this.ws.onmessage = (event) => {
        try {
          const message: LogMessage = JSON.parse(event.data);
          this.options.onMessage(message);
        } catch (error) {
          console.error('Failed to parse log message:', error);
          this.options.onError?.(new Error('Invalid log message format'));
        }
      };

      this.ws.onclose = (event) => {
        console.log('WebSocket connection closed:', event.code, event.reason);
        this.options.onDisconnect?.();
        
        if (this.options.autoReconnect && 
            this.reconnectAttempts < (this.options.maxReconnectAttempts || 5)) {
          this.scheduleReconnect();
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.options.onError?.(new Error('WebSocket connection error'));
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      this.options.onError?.(new Error('Failed to create WebSocket connection'));
    }
  }

  disconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    
    if (this.ws) {
      this.ws.close(1000, 'Client disconnecting');
      this.ws = null;
    }
  }

  private scheduleReconnect(): void {
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000); // Exponential backoff, max 30s
    
    this.reconnectTimer = setTimeout(() => {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.options.maxReconnectAttempts})...`);
      this.connect();
    }, delay);
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}

export class SSELogStreamer {
  private eventSource: EventSource | null = null;
  private options: LogStreamOptions;
  private reconnectAttempts = 0;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private url: string;

  constructor(baseUrl: string, options: LogStreamOptions) {
    this.options = {
      autoReconnect: true,
      maxReconnectAttempts: 5,
      ...options,
    };
    
    const params = new URLSearchParams();
    if (this.options.jobRunId) {
      params.append('job_run_id', this.options.jobRunId);
    }
    
    this.url = `${baseUrl}/api/v1/logs/stream/sse?${params.toString()}`;
  }

  connect(): void {
    try {
      this.eventSource = new EventSource(this.url);
      
      this.eventSource.onopen = () => {
        console.log('SSE connected to log stream');
        this.reconnectAttempts = 0;
        this.options.onConnect?.();
      };

      this.eventSource.addEventListener('connected', (event) => {
        console.log('SSE connection confirmed:', event.data);
      });

      this.eventSource.addEventListener('log', (event) => {
        try {
          const message: LogMessage = JSON.parse(event.data);
          this.options.onMessage(message);
        } catch (error) {
          console.error('Failed to parse log message:', error);
          this.options.onError?.(new Error('Invalid log message format'));
        }
      });

      this.eventSource.addEventListener('ping', (event) => {
        console.log('SSE keep-alive ping received');
      });

      this.eventSource.onerror = (error) => {
        console.error('SSE error:', error);
        this.options.onError?.(new Error('SSE connection error'));
        this.options.onDisconnect?.();
        
        if (this.options.autoReconnect && 
            this.reconnectAttempts < (this.options.maxReconnectAttempts || 5)) {
          this.scheduleReconnect();
        }
      };
    } catch (error) {
      console.error('Failed to create SSE connection:', error);
      this.options.onError?.(new Error('Failed to create SSE connection'));
    }
  }

  disconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
  }

  private scheduleReconnect(): void {
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000); // Exponential backoff, max 30s
    
    this.reconnectTimer = setTimeout(() => {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.options.maxReconnectAttempts})...`);
      this.connect();
    }, delay);
  }

  isConnected(): boolean {
    return this.eventSource?.readyState === EventSource.OPEN;
  }
}

export class LogStreamingService {
  private baseUrl: string;
  private streamers: Map<string, WebSocketLogStreamer | SSELogStreamer> = new Map();

  constructor(baseUrl: string = '') {
    this.baseUrl = baseUrl || window.location.origin;
  }

  // Stream logs for a specific job run
  streamJobRunLogs(
    jobRunId: string, 
    options: Omit<LogStreamOptions, 'jobRunId'>,
    preferWebSocket: boolean = true
  ): WebSocketLogStreamer | SSELogStreamer {
    const key = `job-${jobRunId}`;
    
    // Close existing streamer if any
    this.stopStreaming(key);
    
    const streamOptions = { ...options, jobRunId };
    
    const streamer = preferWebSocket 
      ? new WebSocketLogStreamer(this.baseUrl, streamOptions)
      : new SSELogStreamer(this.baseUrl, streamOptions);
    
    this.streamers.set(key, streamer);
    streamer.connect();
    
    return streamer;
  }

  // Stream all logs (global)
  streamAllLogs(
    options: Omit<LogStreamOptions, 'jobRunId'>,
    preferWebSocket: boolean = true
  ): WebSocketLogStreamer | SSELogStreamer {
    const key = 'global';
    
    // Close existing streamer if any
    this.stopStreaming(key);
    
    const streamer = preferWebSocket 
      ? new WebSocketLogStreamer(this.baseUrl, options)
      : new SSELogStreamer(this.baseUrl, options);
    
    this.streamers.set(key, streamer);
    streamer.connect();
    
    return streamer;
  }

  // Stop streaming for a specific key
  stopStreaming(key: string): void {
    const streamer = this.streamers.get(key);
    if (streamer) {
      streamer.disconnect();
      this.streamers.delete(key);
    }
  }

  // Stop all streaming
  stopAllStreaming(): void {
    this.streamers.forEach((streamer) => {
      streamer.disconnect();
    });
    this.streamers.clear();
  }

  // Get streaming status
  async getStreamingStatus(): Promise<any> {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/logs/stream/status`);
      return await response.json();
    } catch (error) {
      console.error('Failed to get streaming status:', error);
      throw error;
    }
  }

  // Fetch historical logs for a job run
  async getJobRunLogs(jobRunId: string, limit: number = 100, offset: number = 0): Promise<any> {
    try {
      const params = new URLSearchParams({
        limit: limit.toString(),
        offset: offset.toString(),
      });
      
      const response = await fetch(`${this.baseUrl}/api/v1/logs/job-runs/${jobRunId}?${params}`);
      return await response.json();
    } catch (error) {
      console.error('Failed to fetch job run logs:', error);
      throw error;
    }
  }

  // Fetch logs for a specific task run
  async getTaskRunLogs(taskRunId: string): Promise<any> {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/logs/task-runs/${taskRunId}`);
      return await response.json();
    } catch (error) {
      console.error('Failed to fetch task run logs:', error);
      throw error;
    }
  }

  // Fetch system logs
  async getSystemLogs(limit: number = 100, offset: number = 0): Promise<any> {
    try {
      const params = new URLSearchParams({
        limit: limit.toString(),
        offset: offset.toString(),
      });
      
      const response = await fetch(`${this.baseUrl}/api/v1/logs/system?${params}`);
      return await response.json();
    } catch (error) {
      console.error('Failed to fetch system logs:', error);
      throw error;
    }
  }

  // Fetch logs by type
  async getLogsByType(type: 'system' | 'job' | 'task', limit: number = 100, offset: number = 0): Promise<any> {
    try {
      const params = new URLSearchParams({
        limit: limit.toString(),
        offset: offset.toString(),
      });
      
      const response = await fetch(`${this.baseUrl}/api/v1/logs/type/${type}?${params}`);
      return await response.json();
    } catch (error) {
      console.error('Failed to fetch logs by type:', error);
      throw error;
    }
  }

  // Download log file for a job run
  downloadJobRunLogs(jobRunId: string, date?: string): void {
    const params = new URLSearchParams();
    if (date) {
      params.append('date', date);
    }
    
    const url = `${this.baseUrl}/api/v1/logs/job-runs/${jobRunId}/download?${params}`;
    
    // Create a temporary link and click it to trigger download
    const link = document.createElement('a');
    link.href = url;
    link.download = `${jobRunId}-${date || 'latest'}.txt`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
} 
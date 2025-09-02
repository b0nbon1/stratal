import React, { useState } from 'react';
import { jobRunsApi, JobRun } from '../services/api';

interface JobRunControlsProps {
  jobRun: JobRun;
  onStatusChange: (newStatus: string) => void;
}

const JobRunControls: React.FC<JobRunControlsProps> = ({ jobRun, onStatusChange }) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const canPause = jobRun.status === 'running' || jobRun.status === 'queued';
  const canResume = jobRun.status === 'paused';

  const handlePause = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await jobRunsApi.pause(jobRun.id);
      onStatusChange('paused');
      console.log('Job paused:', response.message);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to pause job run');
      console.error('Error pausing job run:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleResume = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await jobRunsApi.resume(jobRun.id);
      onStatusChange(response.new_status);
      console.log('Job resumed:', response.message);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to resume job run');
      console.error('Error resuming job run:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'running': return 'text-blue-600';
      case 'completed': return 'text-green-600';
      case 'failed': return 'text-red-600';
      case 'paused': return 'text-yellow-600';
      case 'queued': return 'text-gray-600';
      case 'pending': return 'text-gray-500';
      default: return 'text-gray-600';
    }
  };

  const getStatusBadge = (status: string) => {
    const baseClasses = 'px-2 py-1 rounded-full text-xs font-medium';
    switch (status) {
      case 'running': return `${baseClasses} bg-blue-100 text-blue-800`;
      case 'completed': return `${baseClasses} bg-green-100 text-green-800`;
      case 'failed': return `${baseClasses} bg-red-100 text-red-800`;
      case 'paused': return `${baseClasses} bg-yellow-100 text-yellow-800`;
      case 'queued': return `${baseClasses} bg-gray-100 text-gray-800`;
      case 'pending': return `${baseClasses} bg-gray-100 text-gray-600`;
      default: return `${baseClasses} bg-gray-100 text-gray-600`;
    }
  };

  return (
    <div className="flex items-center gap-4">
      <div className="flex items-center gap-2">
        <span className="text-sm font-medium text-gray-700">Status:</span>
        <span className={getStatusBadge(jobRun.status)}>
          {jobRun.status.toUpperCase()}
        </span>
      </div>

      {error && (
        <div className="text-red-600 text-sm bg-red-50 px-3 py-1 rounded">
          {error}
        </div>
      )}

      <div className="flex gap-2">
        {canPause && (
          <button
            onClick={handlePause}
            disabled={isLoading}
            className="flex items-center gap-1 px-3 py-1 bg-yellow-500 text-white text-sm rounded hover:bg-yellow-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? (
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            ) : (
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 9v6m4-6v6" />
              </svg>
            )}
            Pause
          </button>
        )}

        {canResume && (
          <button
            onClick={handleResume}
            disabled={isLoading}
            className="flex items-center gap-1 px-3 py-1 bg-green-500 text-white text-sm rounded hover:bg-green-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? (
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            ) : (
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            )}
            Resume
          </button>
        )}
      </div>

      {jobRun.paused_at && (
        <div className="text-xs text-gray-500">
          Paused: {new Date(jobRun.paused_at).toLocaleString()}
        </div>
      )}
    </div>
  );
};

export default JobRunControls;
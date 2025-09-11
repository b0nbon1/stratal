import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { useToast } from '../use-toast';

export interface TaskConfig {
  script: {
    language: string;
    code: string;
  };
  depends_on?: string[];
}

export interface JobTask {
  name: string;
  type: string;
  order: number;
  config: TaskConfig;
}

export interface CreateJobRequest {
  name: string;
  description: string;
  source: string;
  run_immediately: boolean;
  raw_payload: string;
  tasks: JobTask[];
}

export interface JobTaskResponse {
  id: string;
  job_id: string;
  name: string;
  type: string;
  config: {
    script: {
      language: string;
      code: string;
    };
  };
  order: number;
  created_at: string;
}

export interface CreateJobResponse {
  job: {
    id: string;
    user_id: string | null;
    name: string;
    description: string;
    source: string;
    created_at: string;
  };
  job_run_id: string;
  message: string;
  status: string;
  tasks: JobTaskResponse[];
}

const createJob = async (jobData: CreateJobRequest): Promise<CreateJobResponse> => {
  const baseUrl = 'http://localhost:8083';
  const response = await fetch(`${baseUrl}/api/v1/jobs`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(jobData),
  });

  if (!response.ok) {
    throw new Error(`Failed to create job: ${response.statusText}`);
  }

  return response.json();
};

export const useCreateJob = () => {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { toast } = useToast();

  return useMutation({
    mutationFn: createJob,
    onSuccess: () => {
      toast({
      title: "Job Configuration Saved",
      description: "Your job configuration has been saved successfully.",
    });
      queryClient.invalidateQueries({ queryKey: ['jobs', 'jobRuns'] });
      navigate('/');
    },
    onError: (error: any) => {
      toast({
        title: "Error",
        description: error.message || "An error occurred while saving the job configuration.",
        variant: "destructive",
      });
    }
  });
};
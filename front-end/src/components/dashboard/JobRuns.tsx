
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Clock, CheckCircle, XCircle, PlayCircle, Eye, ChevronDown, ChevronRight } from "lucide-react";
import { useState } from "react";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";

const jobRuns = [
  {
    id: 1,
    jobName: "Simple Parallel Task Example",
    status: "running",
    progress: 65,
    startTime: "2024-01-15 10:30:00",
    duration: "2m 15s",
    taskRuns: [
      {
        id: 1,
        taskName: "generate_random_1",
        status: "completed",
        progress: 100,
        duration: "0m 5s",
        output: "42"
      },
      {
        id: 2,
        taskName: "generate_random_2",
        status: "completed", 
        progress: 100,
        duration: "0m 3s",
        output: "78"
      },
      {
        id: 3,
        taskName: "generate_random_3",
        status: "completed",
        progress: 100,
        duration: "0m 4s",
        output: "33"
      },
      {
        id: 4,
        taskName: "calculate_sum",
        status: "running",
        progress: 60,
        duration: "0m 45s",
        output: null
      },
      {
        id: 5,
        taskName: "check_result",
        status: "queued",
        progress: 0,
        duration: "0s",
        output: null
      }
    ]
  },
  {
    id: 2,
    jobName: "Database Backup Job",
    status: "completed",
    progress: 100,
    startTime: "2024-01-15 09:15:00",
    duration: "5m 32s",
    taskRuns: [
      {
        id: 6,
        taskName: "backup_users_table",
        status: "completed",
        progress: 100,
        duration: "2m 15s",
        output: "Backup completed: 15,432 rows"
      },
      {
        id: 7,
        taskName: "backup_orders_table",
        status: "completed",
        progress: 100,
        duration: "3m 17s",
        output: "Backup completed: 8,901 rows"
      }
    ]
  },
  {
    id: 3,
    jobName: "Code Quality Check",
    status: "failed",
    progress: 45,
    startTime: "2024-01-15 08:45:00",
    duration: "1m 23s",
    taskRuns: [
      {
        id: 8,
        taskName: "lint_check",
        status: "failed",
        progress: 100,
        duration: "1m 23s",
        output: "Error: 5 linting errors found"
      }
    ]
  }
];

const getStatusIcon = (status: string) => {
  switch (status) {
    case "running":
      return <PlayCircle className="w-4 h-4 text-blue-600" />;
    case "completed":
      return <CheckCircle className="w-4 h-4 text-green-600" />;
    case "failed":
      return <XCircle className="w-4 h-4 text-red-600" />;
    default:
      return <Clock className="w-4 h-4 text-gray-400" />;
  }
};

const getStatusColor = (status: string) => {
  switch (status) {
    case "running":
      return "bg-blue-100 text-blue-800";
    case "completed":
      return "bg-green-100 text-green-800";
    case "failed":
      return "bg-red-100 text-red-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

export function JobRuns() {
  const [expandedJobs, setExpandedJobs] = useState<number[]>([1]);

  const toggleJobExpansion = (jobId: number) => {
    setExpandedJobs(prev => 
      prev.includes(jobId) 
        ? prev.filter(id => id !== jobId)
        : [...prev, jobId]
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Clock className="w-5 h-5" />
          Job Runs
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {jobRuns.map((jobRun) => (
            <div key={jobRun.id} className="border border-border rounded-lg">
              <Collapsible 
                open={expandedJobs.includes(jobRun.id)}
                onOpenChange={() => toggleJobExpansion(jobRun.id)}
              >
                <CollapsibleTrigger className="w-full p-4 hover:bg-accent/50 transition-colors">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      {expandedJobs.includes(jobRun.id) ? (
                        <ChevronDown className="w-4 h-4" />
                      ) : (
                        <ChevronRight className="w-4 h-4" />
                      )}
                      {getStatusIcon(jobRun.status)}
                      <div className="text-left">
                        <h4 className="font-medium">{jobRun.jobName}</h4>
                        <p className="text-sm text-muted-foreground">
                          Started: {jobRun.startTime} • Duration: {jobRun.duration}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant="secondary" className={getStatusColor(jobRun.status)}>
                        {jobRun.status}
                      </Badge>
                      <Button variant="ghost" size="sm">
                        <Eye className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                  <div className="flex items-center gap-3 mt-2">
                    <Progress value={jobRun.progress} className="flex-1" />
                    <span className="text-sm font-medium min-w-[3rem]">{jobRun.progress}%</span>
                  </div>
                </CollapsibleTrigger>
                
                <CollapsibleContent className="px-4 pb-4">
                  <div className="pl-8 space-y-2">
                    <h5 className="font-medium text-sm text-muted-foreground mb-2">Task Runs</h5>
                    {jobRun.taskRuns.map((taskRun) => (
                      <div key={taskRun.id} className="p-3 bg-accent/30 rounded-lg">
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-2">
                            {getStatusIcon(taskRun.status)}
                            <div>
                              <div className="font-medium text-sm">{taskRun.taskName}</div>
                              <div className="text-xs text-muted-foreground">Duration: {taskRun.duration}</div>
                            </div>
                          </div>
                          <Badge variant="outline" className={getStatusColor(taskRun.status)}>
                            {taskRun.status}
                          </Badge>
                        </div>
                        <div className="flex items-center gap-3 mb-2">
                          <Progress value={taskRun.progress} className="flex-1" />
                          <span className="text-xs font-medium min-w-[3rem]">{taskRun.progress}%</span>
                        </div>
                        {taskRun.output && (
                          <div className="text-xs bg-muted p-2 rounded font-mono">
                            Output: {taskRun.output}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                </CollapsibleContent>
              </Collapsible>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

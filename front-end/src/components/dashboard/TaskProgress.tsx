
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { Clock, CheckCircle, XCircle, PlayCircle } from "lucide-react";

const tasks = [
  {
    id: 1,
    name: "Deploy Production Build",
    status: "running",
    progress: 75,
    duration: "2m 30s",
    type: "deployment"
  },
  {
    id: 2,
    name: "Database Backup",
    status: "completed",
    progress: 100,
    duration: "5m 12s",
    type: "maintenance"
  },
  {
    id: 3,
    name: "Code Quality Check",
    status: "failed",
    progress: 45,
    duration: "1m 45s",
    type: "validation"
  },
  {
    id: 4,
    name: "Email Newsletter",
    status: "queued",
    progress: 0,
    duration: "0s",
    type: "automation"
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

export function TaskProgress() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Clock className="w-5 h-5" />
          Task Progress
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {tasks.map((task) => (
            <div key={task.id} className="p-4 border border-border rounded-lg hover:bg-accent/50 transition-colors">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  {getStatusIcon(task.status)}
                  <div>
                    <h4 className="font-medium">{task.name}</h4>
                    <p className="text-sm text-muted-foreground">Duration: {task.duration}</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary" className={getStatusColor(task.status)}>
                    {task.status}
                  </Badge>
                  <Badge variant="outline">
                    {task.type}
                  </Badge>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <Progress value={task.progress} className="flex-1" />
                <span className="text-sm font-medium min-w-[3rem]">{task.progress}%</span>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

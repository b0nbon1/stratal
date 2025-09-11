
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CalendarClock, Play, Pause, Trash2 } from "lucide-react";

const scheduledTasks = [
  {
    id: 1,
    name: "Daily Database Backup",
    schedule: "0 2 * * *",
    nextRun: "Tomorrow at 2:00 AM",
    status: "active",
    lastRun: "Today at 2:00 AM",
    success: true
  },
  {
    id: 2,
    name: "Weekly Report Generation",
    schedule: "0 9 * * 1",
    nextRun: "Monday at 9:00 AM",
    status: "active",
    lastRun: "Last Monday at 9:00 AM",
    success: true
  },
  {
    id: 3,
    name: "Code Quality Scan",
    schedule: "0 12 * * 1-5",
    nextRun: "Today at 12:00 PM",
    status: "paused",
    lastRun: "Yesterday at 12:00 PM",
    success: false
  },
  {
    id: 4,
    name: "Security Audit",
    schedule: "0 0 1 * *",
    nextRun: "1st of next month",
    status: "active",
    lastRun: "1st of this month",
    success: true
  }
];

export function ScheduledTasks() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <CalendarClock className="w-5 h-5" />
          Scheduled Tasks
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {scheduledTasks.map((task) => (
            <div key={task.id} className="p-4 border border-border rounded-lg">
              <div className="flex items-center justify-between mb-3">
                <div>
                  <h4 className="font-medium">{task.name}</h4>
                  <p className="text-sm text-muted-foreground">Cron: {task.schedule}</p>
                </div>
                <div className="flex items-center gap-2">
                  <Badge 
                    variant={task.status === "active" ? "default" : "secondary"}
                    className={task.status === "active" ? "bg-green-100 text-green-800" : ""}
                  >
                    {task.status}
                  </Badge>
                  <Badge 
                    variant="outline"
                    className={task.success ? "border-green-200 text-green-700" : "border-red-200 text-red-700"}
                  >
                    {task.success ? "Success" : "Failed"}
                  </Badge>
                </div>
              </div>
              
              <div className="grid grid-cols-2 gap-4 text-sm mb-3">
                <div>
                  <span className="text-muted-foreground">Next run:</span>
                  <p className="font-medium">{task.nextRun}</p>
                </div>
                <div>
                  <span className="text-muted-foreground">Last run:</span>
                  <p className="font-medium">{task.lastRun}</p>
                </div>
              </div>
              
              <div className="flex gap-2">
                <Button variant="outline" size="sm">
                  {task.status === "active" ? (
                    <>
                      <Pause className="w-4 h-4 mr-1" />
                      Pause
                    </>
                  ) : (
                    <>
                      <Play className="w-4 h-4 mr-1" />
                      Resume
                    </>
                  )}
                </Button>
                <Button variant="outline" size="sm">
                  <Play className="w-4 h-4 mr-1" />
                  Run Now
                </Button>
                <Button variant="outline" size="sm" className="text-red-600 hover:text-red-700">
                  <Trash2 className="w-4 h-4 mr-1" />
                  Delete
                </Button>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}


import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PlayCircle, Save, FileText, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { useToast } from "@/hooks/use-toast";
import { YamlEditor } from "./YamlEditor";

interface TaskConfig {
  script: {
    language: string;
    code: string;
  };
  depends_on?: string[];
}

interface Task {
  name: string;
  type: string;
  order: number;
  config: TaskConfig;
}

interface Job {
  name: string;
  description: string;
  source: string;
  run_immediately: boolean;
  tasks: Task[];
}

const defaultJob: Job = {
  name: "New Job",
  description: "Job description",
  source: "api",
  run_immediately: false,
  tasks: [
    {
      name: "example_task",
      type: "custom",
      order: 1,
      config: {
        script: {
          language: "bash",
          code: "echo 'Hello World'"
        }
      }
    }
  ]
};

export function JobEditor() {
  const [job, setJob] = useState<Job>(defaultJob);
  const [selectedTask, setSelectedTask] = useState(0);
  const [jsonContent, setJsonContent] = useState(JSON.stringify(job, null, 2));
  const { toast } = useToast();

  const handleSave = () => {
    toast({
      title: "Job Configuration Saved",
      description: "Your job configuration has been saved successfully.",
    });
  };

  const handleRunJob = () => {
    toast({
      title: "Job Started",
      description: "Your job has been queued and will start shortly.",
    });
  };

  const addTask = () => {
    const newTask: Task = {
      name: `task_${job.tasks.length + 1}`,
      type: "custom",
      order: job.tasks.length + 1,
      config: {
        script: {
          language: "bash",
          code: "echo 'New task'"
        }
      }
    };
    setJob(prev => ({ ...prev, tasks: [...prev.tasks, newTask] }));
  };

  const removeTask = (index: number) => {
    setJob(prev => ({ ...prev, tasks: prev.tasks.filter((_, i) => i !== index) }));
    if (selectedTask >= index && selectedTask > 0) {
      setSelectedTask(selectedTask - 1);
    }
  };

  const updateTask = (index: number, field: keyof Task, value: any) => {
    setJob(prev => ({
      ...prev,
      tasks: prev.tasks.map((task, i) => 
        i === index ? { ...task, [field]: value } : task
      )
    }));
  };

  const updateTaskConfig = (index: number, field: keyof TaskConfig, value: any) => {
    setJob(prev => ({
      ...prev,
      tasks: prev.tasks.map((task, i) => 
        i === index ? { 
          ...task, 
          config: { ...task.config, [field]: value }
        } : task
      )
    }));
  };

  const updateTaskScript = (index: number, field: string, value: any) => {
    setJob(prev => ({
      ...prev,
      tasks: prev.tasks.map((task, i) => 
        i === index ? { 
          ...task, 
          config: { 
            ...task.config, 
            script: { ...task.config.script, [field]: value }
          }
        } : task
      )
    }));
  };

  const handleJsonChange = (value: string) => {
    setJsonContent(value);
    try {
      const parsedJob = JSON.parse(value);
      setJob(parsedJob);
    } catch (error) {
      // Invalid JSON, don't update the job
    }
  };

  const UIEditor = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FileText className="w-5 h-5" />
            Job Configuration
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label htmlFor="job-name">Job Name</Label>
              <Input
                id="job-name"
                value={job.name}
                onChange={(e) => setJob(prev => ({ ...prev, name: e.target.value }))}
              />
            </div>
            <div>
              <Label htmlFor="job-source">Source</Label>
              <Select value={job.source} onValueChange={(value) => setJob(prev => ({ ...prev, source: value }))}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="api">API</SelectItem>
                  <SelectItem value="manual">Manual</SelectItem>
                  <SelectItem value="scheduled">Scheduled</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          
          <div>
            <Label htmlFor="job-description">Description</Label>
            <Textarea
              id="job-description"
              value={job.description}
              onChange={(e) => setJob(prev => ({ ...prev, description: e.target.value }))}
              rows={2}
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="run-immediately"
              checked={job.run_immediately}
              onChange={(e) => setJob(prev => ({ ...prev, run_immediately: e.target.checked }))}
            />
            <Label htmlFor="run-immediately">Run Immediately</Label>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Tasks</CardTitle>
            <Button size="sm" onClick={addTask}>
              <Plus className="w-4 h-4 mr-2" />
              Add Task
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex gap-4">
            <div className="w-1/3 space-y-2">
              {job.tasks.map((task, index) => (
                <div
                  key={index}
                  className={`p-3 border rounded-lg cursor-pointer transition-colors ${
                    selectedTask === index ? 'border-primary bg-accent' : 'border-border hover:bg-accent/50'
                  }`}
                  onClick={() => setSelectedTask(index)}
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <div className="font-medium">{task.name}</div>
                      <div className="text-sm text-muted-foreground">Order: {task.order}</div>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={(e) => {
                        e.stopPropagation();
                        removeTask(index);
                      }}
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>

            <div className="flex-1 space-y-4">
              {job.tasks[selectedTask] && (
                <>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <Label>Task Name</Label>
                      <Input
                        value={job.tasks[selectedTask].name}
                        onChange={(e) => updateTask(selectedTask, 'name', e.target.value)}
                      />
                    </div>
                    <div>
                      <Label>Order</Label>
                      <Input
                        type="number"
                        value={job.tasks[selectedTask].order}
                        onChange={(e) => updateTask(selectedTask, 'order', parseInt(e.target.value))}
                      />
                    </div>
                  </div>

                  <div>
                    <Label>Script Language</Label>
                    <Select 
                      value={job.tasks[selectedTask].config.script.language}
                      onValueChange={(value) => updateTaskScript(selectedTask, 'language', value)}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="bash">Bash</SelectItem>
                        <SelectItem value="python">Python</SelectItem>
                        <SelectItem value="javascript">JavaScript</SelectItem>
                        <SelectItem value="sql">SQL</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <Label>Script Code</Label>
                    <Textarea
                      value={job.tasks[selectedTask].config.script.code}
                      onChange={(e) => updateTaskScript(selectedTask, 'code', e.target.value)}
                      className="min-h-[300px] font-mono text-sm"
                      placeholder="Enter your script code here..."
                    />
                  </div>

                  <div>
                    <Label>Dependencies (comma-separated task names)</Label>
                    <Input
                      value={job.tasks[selectedTask].config.depends_on?.join(', ') || ''}
                      onChange={(e) => updateTaskConfig(selectedTask, 'depends_on', e.target.value.split(',').map(s => s.trim()).filter(Boolean))}
                      placeholder="task1, task2, task3"
                    />
                  </div>
                </>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  const JSONEditor = () => (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FileText className="w-5 h-5" />
          JSON Configuration
        </CardTitle>
      </CardHeader>
      <CardContent>
        <Textarea
          value={jsonContent}
          onChange={(e) => handleJsonChange(e.target.value)}
          className="min-h-[500px] font-mono text-sm"
          placeholder="Enter your JSON configuration here..."
        />
        <div className="mt-4 p-3 bg-muted rounded-lg">
          <p className="text-sm text-muted-foreground">
            <strong>Tip:</strong> Use valid JSON format to define your job configuration.
            Changes will be reflected in the UI editor automatically.
          </p>
        </div>
      </CardContent>
    </Card>
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Job Editor</h2>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={handleSave}>
            <Save className="w-4 h-4 mr-2" />
            Save
          </Button>
          <Button size="sm" onClick={handleRunJob} className="bg-gradient-to-r from-blue-500 to-purple-600">
            <PlayCircle className="w-4 h-4 mr-2" />
            Run Job
          </Button>
        </div>
      </div>

      <Tabs defaultValue="ui" className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="ui">UI Editor</TabsTrigger>
          <TabsTrigger value="json">JSON Editor</TabsTrigger>
          <TabsTrigger value="yaml">YAML Editor</TabsTrigger>
        </TabsList>
        
        <TabsContent value="ui" className="space-y-4">
          <UIEditor />
        </TabsContent>
        
        <TabsContent value="json" className="space-y-4">
          <JSONEditor />
        </TabsContent>
        
        <TabsContent value="yaml" className="space-y-4">
          <YamlEditor />
        </TabsContent>
      </Tabs>
    </div>
  );
}

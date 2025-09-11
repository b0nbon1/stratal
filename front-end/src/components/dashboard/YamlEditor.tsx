
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { PlayCircle, Save, FileText } from "lucide-react";
import { useState } from "react";
import { useToast } from "@/hooks/use-toast";

const defaultYaml = `# TaskFlow Configuration
name: "My Automation Task"
description: "Sample automation workflow"

steps:
  - name: "Setup Environment"
    type: "command"
    run: "npm install"
    
  - name: "Run Tests"
    type: "command"
    run: "npm test"
    
  - name: "Build Project"
    type: "command"
    run: "npm run build"
    
  - name: "Deploy"
    type: "deployment"
    target: "production"
    
schedule:
  cron: "0 9 * * 1-5"  # Weekdays at 9 AM
  
notifications:
  on_success: true
  on_failure: true
  email: "admin@example.com"`;

export function YamlEditor() {
  const [yamlContent, setYamlContent] = useState(defaultYaml);
  const { toast } = useToast();

  const handleSave = () => {
    toast({
      title: "Configuration Saved",
      description: "Your YAML configuration has been saved successfully.",
    });
  };

  const handleRun = () => {
    toast({
      title: "Task Started",
      description: "Your automation task has been queued and will start shortly.",
    });
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <FileText className="w-5 h-5" />
            YAML Configuration
          </CardTitle>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" onClick={handleSave}>
              <Save className="w-4 h-4 mr-2" />
              Save
            </Button>
            <Button size="sm" onClick={handleRun} className="bg-gradient-to-r from-blue-500 to-purple-600">
              <PlayCircle className="w-4 h-4 mr-2" />
              Run Task
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <Textarea
          value={yamlContent}
          onChange={(e) => setYamlContent(e.target.value)}
          className="min-h-[400px] font-mono text-sm"
          placeholder="Enter your YAML configuration here..."
        />
        <div className="mt-4 p-3 bg-muted rounded-lg">
          <p className="text-sm text-muted-foreground">
            <strong>Tip:</strong> Use the YAML format to define your automation tasks. 
            Include steps, scheduling, and notification preferences.
          </p>
        </div>
      </CardContent>
    </Card>
  );
}

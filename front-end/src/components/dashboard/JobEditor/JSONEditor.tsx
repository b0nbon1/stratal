import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { FileText } from "lucide-react";
import React, { useState } from "react";

function JSONEditor({ job, setJob }) {
  const [jsonContent, setJsonContent] = useState(JSON.stringify(job, null, 2));

  const handleJsonChange = (value: string) => {
    setJsonContent(value);
    try {
      const parsedJob = JSON.parse(value);
      setJob(parsedJob);
    } catch (error) {
      // Invalid JSON, don't update the job
    }
  };
  return (
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
            <strong>Tip:</strong> Use valid JSON format to define your job
            configuration. Changes will be reflected in the UI editor
            automatically.
          </p>
        </div>
      </CardContent>
    </Card>
  );
}

export default JSONEditor;

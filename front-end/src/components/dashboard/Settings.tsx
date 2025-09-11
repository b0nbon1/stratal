import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Plus, Trash2, Eye, EyeOff, Key, Lock, Settings as SettingsIcon } from "lucide-react";
import { useState } from "react";
import { useToast } from "@/hooks/use-toast";

interface EnvironmentVariable {
  id: string;
  name: string;
  value: string;
  encrypted: boolean;
  description?: string;
}

const defaultEnvVars: EnvironmentVariable[] = [
  {
    id: "1",
    name: "API_BASE_URL",
    value: "https://api.example.com",
    encrypted: false,
    description: "Base URL for API endpoints"
  },
  {
    id: "2",
    name: "DATABASE_PASSWORD",
    value: "••••••••••••",
    encrypted: true,
    description: "Database connection password"
  }
];

export function Settings() {
  const [envVars, setEnvVars] = useState<EnvironmentVariable[]>(defaultEnvVars);
  const [showValues, setShowValues] = useState<Record<string, boolean>>({});
  const [newVar, setNewVar] = useState({
    name: "",
    value: "",
    encrypted: false,
    description: ""
  });
  const { toast } = useToast();

  const addEnvironmentVariable = () => {
    if (!newVar.name || !newVar.value) {
      toast({
        title: "Validation Error",
        description: "Please provide both name and value for the environment variable.",
        variant: "destructive",
      });
      return;
    }

    const newEnvVar: EnvironmentVariable = {
      id: Date.now().toString(),
      name: newVar.name,
      value: newVar.value,
      encrypted: newVar.encrypted,
      description: newVar.description
    };

    setEnvVars(prev => [...prev, newEnvVar]);
    setNewVar({ name: "", value: "", encrypted: false, description: "" });
    
    toast({
      title: "Environment Variable Added",
      description: `${newVar.name} has been added successfully.`,
    });
  };

  const removeEnvironmentVariable = (id: string) => {
    setEnvVars(prev => prev.filter(env => env.id !== id));
    toast({
      title: "Environment Variable Removed",
      description: "The environment variable has been removed.",
    });
  };

  const toggleValueVisibility = (id: string) => {
    setShowValues(prev => ({ ...prev, [id]: !prev[id] }));
  };

  const updateEnvironmentVariable = (id: string, field: keyof EnvironmentVariable, value: any) => {
    setEnvVars(prev => prev.map(env => 
      env.id === id ? { ...env, [field]: value } : env
    ));
  };

  const handleSave = () => {
    toast({
      title: "Settings Saved",
      description: "Environment variables have been saved successfully.",
    });
  };

  const handleExport = () => {
    const exportData = {
      environment_variables: envVars.map(env => ({
        name: env.name,
        value: env.encrypted ? "[ENCRYPTED]" : env.value,
        encrypted: env.encrypted,
        description: env.description
      })),
      exported_at: new Date().toISOString()
    };

    const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'environment-variables.json';
    a.click();
    URL.revokeObjectURL(url);

    toast({
      title: "Export Complete",
      description: "Environment variables exported successfully.",
    });
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <SettingsIcon className="w-5 h-5" />
              Environment Variables
            </CardTitle>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" onClick={handleExport}>
                Export
              </Button>
              <Button size="sm" onClick={handleSave}>
                Save Settings
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Add New Environment Variable */}
          <div className="space-y-4 p-4 border border-dashed rounded-lg">
            <h3 className="font-medium">Add New Environment Variable</h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="new-env-name">Name</Label>
                <Input
                  id="new-env-name"
                  value={newVar.name}
                  onChange={(e) => setNewVar(prev => ({ ...prev, name: e.target.value }))}
                  placeholder="VARIABLE_NAME"
                />
              </div>
              <div>
                <Label htmlFor="new-env-value">Value</Label>
                <Input
                  id="new-env-value"
                  type={newVar.encrypted ? "password" : "text"}
                  value={newVar.value}
                  onChange={(e) => setNewVar(prev => ({ ...prev, value: e.target.value }))}
                  placeholder="Enter value"
                />
              </div>
            </div>
            <div>
              <Label htmlFor="new-env-description">Description (Optional)</Label>
              <Textarea
                id="new-env-description"
                value={newVar.description}
                onChange={(e) => setNewVar(prev => ({ ...prev, description: e.target.value }))}
                placeholder="Brief description of this variable"
                rows={2}
              />
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Switch
                  id="new-env-encrypted"
                  checked={newVar.encrypted}
                  onCheckedChange={(checked) => setNewVar(prev => ({ ...prev, encrypted: checked }))}
                />
                <Label htmlFor="new-env-encrypted" className="flex items-center gap-2">
                  <Lock className="w-4 h-4" />
                  Encrypt Value
                </Label>
              </div>
              <Button onClick={addEnvironmentVariable}>
                <Plus className="w-4 h-4 mr-2" />
                Add Variable
              </Button>
            </div>
          </div>

          <Separator />

          {/* Existing Environment Variables */}
          <div className="space-y-4">
            <h3 className="font-medium">Existing Environment Variables</h3>
            {envVars.map((envVar) => (
              <div key={envVar.id} className="p-4 border rounded-lg space-y-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Badge variant={envVar.encrypted ? "default" : "secondary"}>
                      {envVar.encrypted ? (
                        <>
                          <Key className="w-3 h-3 mr-1" />
                          Encrypted
                        </>
                      ) : (
                        "Plain Text"
                      )}
                    </Badge>
                    <span className="font-mono text-sm">{envVar.name}</span>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => removeEnvironmentVariable(envVar.id)}
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </div>
                
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label>Name</Label>
                    <Input
                      value={envVar.name}
                      onChange={(e) => updateEnvironmentVariable(envVar.id, 'name', e.target.value)}
                    />
                  </div>
                  <div>
                    <Label>Value</Label>
                    <div className="flex gap-2">
                      <Input
                        type={envVar.encrypted && !showValues[envVar.id] ? "password" : "text"}
                        value={envVar.value}
                        onChange={(e) => updateEnvironmentVariable(envVar.id, 'value', e.target.value)}
                      />
                      {envVar.encrypted && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => toggleValueVisibility(envVar.id)}
                        >
                          {showValues[envVar.id] ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                        </Button>
                      )}
                    </div>
                  </div>
                </div>

                {envVar.description && (
                  <div>
                    <Label>Description</Label>
                    <Textarea
                      value={envVar.description}
                      onChange={(e) => updateEnvironmentVariable(envVar.id, 'description', e.target.value)}
                      rows={2}
                    />
                  </div>
                )}

                <div className="flex items-center gap-2">
                  <Switch
                    checked={envVar.encrypted}
                    onCheckedChange={(checked) => updateEnvironmentVariable(envVar.id, 'encrypted', checked)}
                  />
                  <Label className="flex items-center gap-2">
                    <Lock className="w-4 h-4" />
                    Encrypt Value
                  </Label>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Encryption Settings</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="p-4 bg-muted rounded-lg">
            <h4 className="font-medium mb-2">Encryption Status</h4>
            <p className="text-sm text-muted-foreground mb-3">
              Environment variables marked as encrypted are stored securely and can only be accessed during job execution.
            </p>
            <div className="flex items-center gap-2">
              <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
                <Lock className="w-3 h-3 mr-1" />
                Encryption Enabled
              </Badge>
            </div>
          </div>
          
          <div className="space-y-2">
            <h4 className="font-medium">Security Notes</h4>
            <ul className="text-sm text-muted-foreground space-y-1">
              <li>• Encrypted variables are masked in the UI and logs</li>
              <li>• Values are encrypted at rest using AES-256</li>
              <li>• Access is restricted to authorized job execution contexts</li>
              <li>• Regular rotation of sensitive variables is recommended</li>
            </ul>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

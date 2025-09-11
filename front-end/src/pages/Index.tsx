import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/layout/AppSidebar";
import { StatsCards } from "@/components/dashboard/StatsCards";
import { JobRuns } from "@/components/dashboard/JobRuns";
import { JobEditor } from "@/components/dashboard/JobEditor";
import { ScheduledTasks } from "@/components/dashboard/ScheduledTasks";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { RefreshCw } from "lucide-react";
import { Settings } from "@/components/dashboard/Settings";

const Index = () => {
  return (
    <SidebarProvider>
      <div className="min-h-screen flex w-full bg-background">
        <AppSidebar />
        <main className="flex-1 p-6">
          <div className="max-w-7xl mx-auto">
            <div className="flex items-center justify-between mb-8">
              <div className="flex items-center gap-4">
                <SidebarTrigger />
                <div>
                  <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                    Stratal Dashboard
                  </h1>
                  <p className="text-muted-foreground">
                    Automate jobs, schedule tasks, and monitor execution runs
                  </p>
                </div>
              </div>
              <Button variant="outline" size="sm">
                <RefreshCw className="w-4 h-4 mr-2" />
                Refresh
              </Button>
            </div>

            <StatsCards />

            <Tabs defaultValue="overview" className="space-y-6">
              <TabsList className="grid w-full grid-cols-5">
                <TabsTrigger value="overview">Overview</TabsTrigger>
                <TabsTrigger value="editor">Job Editor</TabsTrigger>
                <TabsTrigger value="scheduled">Scheduled</TabsTrigger>
                <TabsTrigger value="history">History</TabsTrigger>
                <TabsTrigger value="settings">Settings</TabsTrigger>
              </TabsList>

              <TabsContent value="overview" className="space-y-6">
                <div className="grid grid-cols-1 gap-6">
                  <JobRuns />
                  <ScheduledTasks />
                </div>
              </TabsContent>

              <TabsContent value="editor" className="space-y-6">
                <JobEditor />
              </TabsContent>

              <TabsContent value="scheduled" className="space-y-6">
                <ScheduledTasks />
              </TabsContent>

              <TabsContent value="history" className="space-y-6">
                <div className="text-center py-12">
                  <h3 className="text-lg font-medium mb-2">Job History</h3>
                  <p className="text-muted-foreground">
                    View detailed logs and execution history of your jobs and tasks
                  </p>
                </div>
              </TabsContent>

              <TabsContent value="settings" className="space-y-6">
                <Settings />
              </TabsContent>
            </Tabs>
          </div>
        </main>
      </div>
    </SidebarProvider>
  );
};

export default Index;

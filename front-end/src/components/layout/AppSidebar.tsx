
import {
  Calendar,
  Clock,
  List,
  SquareKanban,
  SquareCheck,
  SquareX,
  CalendarCheck,
  CalendarPlus
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarHeader,
  SidebarFooter,
} from "@/components/ui/sidebar";

const navigationItems = [
  {
    title: "Dashboard",
    url: "#dashboard",
    icon: SquareKanban,
  },
  {
    title: "Tasks",
    url: "#tasks",
    icon: List,
  },
  {
    title: "Completed",
    url: "#completed",
    icon: SquareCheck,
  },
  {
    title: "Failed",
    url: "#failed",
    icon: SquareX,
  },
  {
    title: "Scheduled",
    url: "#scheduled",
    icon: CalendarCheck,
  },
  {
    title: "Create Task",
    url: "#create",
    icon: CalendarPlus,
  },
];

export function AppSidebar() {
  return (
    <Sidebar className="border-r border-border">
      <SidebarHeader className="p-6">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
            <SquareKanban className="w-4 h-4 text-white" />
          </div>
          <h2 className="text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
            Stratal
          </h2>
        </div>
      </SidebarHeader>
      
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navigationItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild className="hover:bg-accent transition-colors">
                    <a href={item.url} className="flex items-center gap-3">
                      <item.icon className="w-4 h-4" />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      
      <SidebarFooter className="p-6">
        <div className="text-sm text-muted-foreground">
          <div className="flex items-center gap-2 mb-2">
            <Clock className="w-4 h-4" />
            <span>System Status: Online</span>
          </div>
          <div className="text-xs">
            Last updated: {new Date().toLocaleTimeString()}
          </div>
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}

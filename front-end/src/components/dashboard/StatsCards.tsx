
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Clock, SquareCheck, SquareX, Calendar } from "lucide-react";

const stats = [
  {
    title: "Running Jobs",
    value: "3",
    icon: Clock,
    color: "text-blue-600",
    bgColor: "bg-blue-50",
    change: "+1 from yesterday"
  },
  {
    title: "Completed Jobs",
    value: "127",
    icon: SquareCheck,
    color: "text-green-600",
    bgColor: "bg-green-50",
    change: "+8 from yesterday"
  },
  {
    title: "Failed Jobs",
    value: "4",
    icon: SquareX,
    color: "text-red-600",
    bgColor: "bg-red-50",
    change: "-1 from yesterday"
  },
  {
    title: "Scheduled Jobs",
    value: "15",
    icon: Calendar,
    color: "text-purple-600",
    bgColor: "bg-purple-50",
    change: "+2 from yesterday"
  }
];

export function StatsCards() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      {stats.map((stat) => (
        <Card key={stat.title} className="hover-scale">
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              <div className={`p-2 rounded-lg ${stat.bgColor}`}>
                <stat.icon className={`w-4 h-4 ${stat.color}`} />
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stat.value}</div>
            <p className="text-xs text-muted-foreground mt-1">{stat.change}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

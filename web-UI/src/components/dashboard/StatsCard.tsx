import {
  Card,
  CardContent,
  CardHeader,
  Typography,
  Box,
  Grid,
} from "@mui/material";

import AccessTimeIcon from "@mui/icons-material/AccessTime";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import CancelIcon from "@mui/icons-material/Cancel";
import EventNoteIcon from "@mui/icons-material/EventNote";

const stats = [
  {
    title: "Running Tasks",
    value: "12",
    icon: <AccessTimeIcon fontSize="small" sx={{ color: "#2563eb" }} />,
    bgColor: "#eff6ff",
    change: "+2 from yesterday",
  },
  {
    title: "Completed",
    value: "234",
    icon: <CheckCircleIcon fontSize="small" sx={{ color: "#16a34a" }} />,
    bgColor: "#f0fdf4",
    change: "+15 from yesterday",
  },
  {
    title: "Failed",
    value: "8",
    icon: <CancelIcon fontSize="small" sx={{ color: "#dc2626" }} />,
    bgColor: "#fef2f2",
    change: "-3 from yesterday",
  },
  {
    title: "Scheduled",
    value: "45",
    icon: <EventNoteIcon fontSize="small" sx={{ color: "#7e22ce" }} />,
    bgColor: "#faf5ff",
    change: "+5 from yesterday",
  },
];

export function StatsCards() {
  return (
    <Grid container spacing={3} mb={6}>
      {stats.map((stat) => (
        <Grid size={{ xs: 12, md: 6, lg: 3 }} key={stat.title}>
          <Card
            elevation={1}
            sx={{
              transition: "transform 0.2s",
              "&:hover": { transform: "scale(1.02)" },
            }}
          >
            <CardHeader
              sx={{
                pb: 1,
              }}
              title={
                <Box display="flex" justifyContent="space-between" alignItems="center">
                  <Typography variant="subtitle2" color="text.secondary">
                    {stat.title}
                  </Typography>
                  <Box
                    sx={{
                      backgroundColor: stat.bgColor,
                      p: 1,
                      borderRadius: 2,
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                    }}
                  >
                    {stat.icon}
                  </Box>
                </Box>
              }
            />
            <CardContent>
              <Typography variant="h5" fontWeight="bold">
                {stat.value}
              </Typography>
              <Typography variant="caption" color="text.secondary" mt={1}>
                {stat.change}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
}

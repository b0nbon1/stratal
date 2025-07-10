import DashboardIcon from "@mui/icons-material/Dashboard";
import ListAltIcon from "@mui/icons-material/ListAlt";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import CancelIcon from "@mui/icons-material/Cancel";
import EventAvailableIcon from "@mui/icons-material/EventAvailable";
import AddCircleIcon from "@mui/icons-material/AddCircle";
import QueryBuilderIcon from "@mui/icons-material/QueryBuilder";
import LogoutIcon from "@mui/icons-material/Logout";

import {
  Drawer,
  List as MUIList,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Divider,
  Typography,
  Box,
  Button,
} from "@mui/material";

const navigationItems = [
  {
    title: "Dashboard",
    url: "#dashboard",
    icon: <DashboardIcon fontSize="small" />,
  },
  {
    title: "Tasks",
    url: "#tasks",
    icon: <ListAltIcon fontSize="small" />,
  },
  {
    title: "Completed",
    url: "#completed",
    icon: <CheckCircleIcon fontSize="small" />,
  },
  {
    title: "Failed",
    url: "#failed",
    icon: <CancelIcon fontSize="small" />,
  },
  {
    title: "Scheduled",
    url: "#scheduled",
    icon: <EventAvailableIcon fontSize="small" />,
  },
  {
    title: "Create Task",
    url: "#create",
    icon: <AddCircleIcon fontSize="small" />,
  },
];

export function AppSidebar() {
  return (
    <Drawer
      variant="permanent"
      anchor="left"
      sx={{
        width: 240,
        flexShrink: 0,
        "& .MuiDrawer-paper": {
          width: 240,
          boxSizing: "border-box",
          borderRight: "1px solid #e0e0e0",
          display: "flex",
          flexDirection: "column",
        },
      }}
    >
      {/* Header */}
      <Box p={2} display="flex" alignItems="center" gap={1}>
        <Box
          width={32}
          height={32}
          display="flex"
          alignItems="center"
          justifyContent="center"
          borderRadius={1}
          sx={{
            background: "linear-gradient(to bottom right, #3b82f6, #8b5cf6)",
          }}
        >
          <DashboardIcon style={{ color: "white", fontSize: 16 }} />
        </Box>
        <Typography
          variant="h6"
          fontWeight="bold"
          sx={{
            background: "linear-gradient(to right, #2563eb, #7c3aed)",
            WebkitBackgroundClip: "text",
            WebkitTextFillColor: "transparent",
          }}
        >
          Stratal
        </Typography>
      </Box>

      {/* Navigation */}
      <Divider />
      <Box px={2} pt={2}>
        <Typography variant="caption" color="textSecondary" gutterBottom>
          Navigation
        </Typography>
      </Box>
      <MUIList>
        {navigationItems.map((item) => (
          <ListItem key={item.title} disablePadding>
            <ListItemButton component="a" href={item.url}>
              <ListItemIcon sx={{ minWidth: 32 }}>{item.icon}</ListItemIcon>
              <ListItemText primary={item.title} />
            </ListItemButton>
          </ListItem>
        ))}
      </MUIList>

      <Box flexGrow={1} />

      {/* Logout */}
      <Box px={2} py={1}>
        <Button
          variant="text"
          startIcon={<LogoutIcon />}
          onClick={() => {
            console.log("Logging out...");
          }}
        >
          Logout
        </Button>
      </Box>

      {/* Footer */}
      <Box px={2} pb={2}>
        <Box display="flex" alignItems="center" gap={1} mb={1}>
          <QueryBuilderIcon fontSize="small" />
          <Typography variant="body2" color="textSecondary">
            System Status: Online
          </Typography>
        </Box>
        <Typography variant="caption" color="textSecondary">
          Last updated: {new Date().toLocaleTimeString()}
        </Typography>
      </Box>
    </Drawer>
  );
}

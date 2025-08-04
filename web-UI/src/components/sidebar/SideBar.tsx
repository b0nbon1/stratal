import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
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
} from '@mui/material';
import DashboardIcon from '@mui/icons-material/Dashboard';
import WorkIcon from '@mui/icons-material/Work';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import VpnKeyIcon from '@mui/icons-material/VpnKey';
import DescriptionIcon from '@mui/icons-material/Description';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import QueryBuilderIcon from '@mui/icons-material/QueryBuilder';
import LogoutIcon from '@mui/icons-material/Logout';

const navigationItems = [
  {
    title: 'Dashboard',
    path: '/dashboard',
    icon: <DashboardIcon fontSize="small" />,
  },
  {
    title: 'Jobs',
    path: '/jobs',
    icon: <WorkIcon fontSize="small" />,
  },
  {
    title: 'Job Runs',
    path: '/job-runs',
    icon: <PlayArrowIcon fontSize="small" />,
  },
  {
    title: 'Create Job',
    path: '/jobs/create',
    icon: <AddCircleIcon fontSize="small" />,
  },
  {
    title: 'Secrets',
    path: '/secrets',
    icon: <VpnKeyIcon fontSize="small" />,
  },
  {
    title: 'Logs',
    path: '/logs',
    icon: <DescriptionIcon fontSize="small" />,
  },
];

export function AppSidebar() {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <Drawer
      variant="permanent"
      anchor="left"
      sx={{
        width: 240,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: 240,
          boxSizing: 'border-box',
          borderRight: '1px solid #e0e0e0',
          display: 'flex',
          flexDirection: 'column',
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
            background: 'linear-gradient(to bottom right, #3b82f6, #8b5cf6)',
          }}
        >
          <DashboardIcon style={{ color: 'white', fontSize: 16 }} />
        </Box>
        <Typography
          variant="h6"
          fontWeight="bold"
          sx={{
            background: 'linear-gradient(to right, #2563eb, #7c3aed)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
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
            <ListItemButton
              onClick={() => navigate(item.path)}
              selected={location.pathname === item.path}
              sx={{
                '&.Mui-selected': {
                  backgroundColor: '#f3f4f6',
                  '&:hover': {
                    backgroundColor: '#e5e7eb',
                  },
                },
              }}
            >
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
            console.log('Logging out...');
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

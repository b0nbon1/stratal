import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Box, CssBaseline, ThemeProvider } from '@mui/material';
import { AppSidebar } from './components/sidebar/SideBar';
import { Dashboard } from './pages/Dashboard';
import { Jobs } from './pages/Jobs';
import { JobDetails } from './pages/JobDetails';
import { CreateJob } from './pages/CreateJob';
import { JobRuns } from './pages/JobRuns';
import { JobRunDetails } from './pages/JobRunDetails';
import { Secrets } from './pages/Secrets';
import { Logs } from './pages/Logs';
import { theme } from './styles/theme';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Box sx={{ display: 'flex' }}>
          <AppSidebar />
          <Box 
            component="main" 
            sx={{ 
              flexGrow: 1, 
              p: 3, 
              width: { sm: `calc(100% - 240px)` },
              ml: { sm: `240px` }
            }}
          >
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/jobs" element={<Jobs />} />
              <Route path="/jobs/create" element={<CreateJob />} />
              <Route path="/jobs/:id" element={<JobDetails />} />
              <Route path="/job-runs" element={<JobRuns />} />
              <Route path="/job-runs/:id" element={<JobRunDetails />} />
              <Route path="/secrets" element={<Secrets />} />
              <Route path="/logs" element={<Logs />} />
            </Routes>
          </Box>
        </Box>
      </Router>
    </ThemeProvider>
  );
}

export default App;

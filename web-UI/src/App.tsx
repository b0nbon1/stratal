import { Box } from "@mui/material";
import { StatsCards } from "./components/dashboard/StatsCard";
import { AppSidebar } from "./components/sidebar/SideBar";

function App() {
  return (
    <>
      <AppSidebar />
      <Box sx={{ pl: {xs: 4, sm: 5,  md: 32 }, mt: 4, pr: 4 }}>
        <StatsCards />
      </Box>
    </>
  );
}

export default App;

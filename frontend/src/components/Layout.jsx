import React, { useState } from 'react';
import { Box, AppBar, Toolbar, IconButton, Typography, Drawer, CssBaseline, Button } from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import BugReportIcon from '@mui/icons-material/BugReport';
import { useAuth } from '../context/AuthContext';
import SidebarContent from './SidebarContent';

const drawerWidth = 240;

const Layout = ({ children }) => {
  const [mobileOpen, setMobileOpen] = useState(false);
  const { user, logout } = useAuth();

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const drawer = <SidebarContent />;

  const drawerStyles = {
    '& .MuiDrawer-paper': {
      boxSizing: 'border-box',
      width: drawerWidth,
      backgroundColor: 'background.default', 
      borderRight: '1px solid',
      borderColor: 'divider',
    },
  };

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar
        position="fixed"
        elevation={0}
        sx={{
          width: { md: `calc(100% - ${drawerWidth}px)` },
          ml: { md: `${drawerWidth}px` },
          width: { xs: '100%' },
          backgroundColor: 'background.paper',
          borderBottom: '1px solid',
          borderColor: 'divider',
        }}
      >
        <Toolbar>
          <IconButton
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { md: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          
          <Typography 
            variant="h6" 
            noWrap 
            component="div" 
            color="text.primary"
            sx={{ flexGrow: 1, display: { xs: 'block', md: 'none' } }}
          >
            BugTracker
          </Typography>

          <Box sx={{ flexGrow: 1, display: { xs: 'none', md: 'block' } }} />

          <Button onClick={logout}>
            Выйти
          </Button>
        </Toolbar>
      </AppBar>

      <Box
        component="nav"
        sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
        aria-label="folders"
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{ keepMounted: true }}
          sx={{
            ...drawerStyles,
            display: { xs: 'block', md: 'none' },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            ...drawerStyles,
            display: { xs: 'none', md: 'block' },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { xs: '100%', md: `calc(100% - ${drawerWidth}px)` },
        }}
      >
        <Toolbar />
        {children}
      </Box>
    </Box>
  );
};

export default Layout;
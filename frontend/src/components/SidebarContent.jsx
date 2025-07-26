import React from 'react';
import { NavLink } from 'react-router-dom';
import { Box, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Divider, Toolbar, Typography } from '@mui/material';
import DashboardIcon from '@mui/icons-material/Dashboard';
import FolderIcon from '@mui/icons-material/Folder';
import AssignmentIndIcon from '@mui/icons-material/AssignmentInd';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import BugReportIcon from '@mui/icons-material/BugReport';
import { useAuth } from '../context/AuthContext';

const SidebarContent = () => {
  const { isAdmin } = useAuth();

  const navLinkStyles = ({ isActive }) => ({
    backgroundColor: isActive ? 'rgba(144, 202, 249, 0.16)' : 'transparent',
    borderRight: isActive ? '3px solid #90caf9' : 'none',
  });

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <Box>
        <Toolbar>
          <BugReportIcon sx={{ mr: 1, color: 'primary.main' }} />
          <Typography variant="h6" noWrap component="div">
            BugTracker
          </Typography>
        </Toolbar>
        <Divider />
        <List>
          <ListItem disablePadding>
            <NavLink to="/dashboard" style={{ textDecoration: 'none', color: 'inherit', width: '100%' }}>
              {({ isActive }) => (
                <ListItemButton sx={navLinkStyles({ isActive })}>
                  <ListItemIcon><DashboardIcon /></ListItemIcon>
                  <ListItemText primary="Главная" />
                </ListItemButton>
              )}
            </NavLink>
          </ListItem>
          
          <ListItem disablePadding>
            <NavLink end to="/tasks" style={{ textDecoration: 'none', color: 'inherit', width: '100%' }}>
              {({ isActive }) => (
                <ListItemButton sx={navLinkStyles({ isActive })}>
                  <ListItemIcon><AssignmentIndIcon /></ListItemIcon>
                  <ListItemText primary="Мои задачи" />
                </ListItemButton>
              )}
            </NavLink>
          </ListItem>

          {isAdmin && (
            <ListItem disablePadding>
              <NavLink end to="/folders" style={{ textDecoration: 'none', color: 'inherit', width: '100%' }}>
                {({ isActive }) => (
                  <ListItemButton sx={navLinkStyles({ isActive })}>
                    <ListItemIcon><FolderIcon /></ListItemIcon>
                    <ListItemText primary="Папки" />
                  </ListItemButton>
                )}
              </NavLink>
            </ListItem>
          )}
        </List>
      </Box>

      <Box sx={{ marginTop: 'auto' }}>
        <Divider />
        <List>
          <ListItem disablePadding>
            <NavLink end to="/profile" style={{ textDecoration: 'none', color: 'inherit', width: '100%' }}>
              {({ isActive }) => (
                <ListItemButton sx={navLinkStyles({ isActive })}>
                  <ListItemIcon><AccountCircleIcon /></ListItemIcon>
                  <ListItemText primary="Мой профиль" />
                </ListItemButton>
              )}
            </NavLink>
          </ListItem>
        </List>
      </Box>
    </Box>
  );
};

export default SidebarContent;
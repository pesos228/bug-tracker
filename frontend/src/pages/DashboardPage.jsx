import React from 'react';
import { useAuth } from '../context/AuthContext.jsx';
import { Box, Typography, Container, Button, AppBar, Toolbar } from '@mui/material';
import BugReportIcon from '@mui/icons-material/BugReport';

const DashboardPage = () => {
  const { user, isAdmin, logout } = useAuth();

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <AppBar position="static">
        <Toolbar>
          <BugReportIcon sx={{ mr: 1 }} />
          <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
            BugTracker Dashboard
          </Typography>
          <Button color="inherit" onClick={logout}>
            Выйти
          </Button>
        </Toolbar>
      </AppBar>

      <Container component="main" sx={{ py: 4 }}>
        <Typography variant="h4" gutterBottom>
          Добро пожаловать, {user?.FirstName}!
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Это ваша панель управления. Здесь вы можете управлять задачами.
        </Typography>

        {isAdmin && (
          <Box sx={{ mt: 4, p: 2, border: '1px dashed grey', borderRadius: 1 }}>
            <Typography variant="h6" color="secondary">
              Панель администратора
            </Typography>
            <Typography>Здесь будут кнопки и функции, доступные только админам.</Typography>
            <Button variant="contained" color="secondary" sx={{ mt: 2 }}>
              Управление пользователями
            </Button>
          </Box>
        )}

        <Box sx={{ mt: 4 }}>
          <Typography variant="h5">Ваши задачи</Typography>
        </Box>
      </Container>
    </Box>
  );
};

export default DashboardPage;
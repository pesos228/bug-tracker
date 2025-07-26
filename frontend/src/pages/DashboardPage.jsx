import React from 'react';
import { useAuth } from '../context/AuthContext.jsx';
import { Box, Typography, Container, Button } from '@mui/material';

const DashboardPage = () => {
  const { user, isAdmin } = useAuth();

  return (
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
  );
};

export default DashboardPage;
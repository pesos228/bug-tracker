import React, { useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { Box, CircularProgress, Typography } from '@mui/material';

const LoginPage = () => {
  const { login, isAuthenticated, isInitializing } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isInitializing) {
      if (isAuthenticated) {
        navigate('/dashboard', { replace: true });
      } else {
        login();
      }
    }
  }, [isInitializing, isAuthenticated, login, navigate]);

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        flexDirection: 'column',
        gap: 2,
      }}
    >
      <CircularProgress />
      <Typography variant="body1" color="text.secondary">
        {isInitializing ? 'Проверка статуса...' : 'Перенаправление на страницу входа...'}
      </Typography>
    </Box>
  );
};

export default LoginPage;
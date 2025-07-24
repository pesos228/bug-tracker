import React from 'react';
import { useAuth } from '../context/AuthContext';
import { Box, CircularProgress, Alert } from '@mui/material';

const ProtectedRoute = ({ children }) => {
  const { isAuthenticated, isInitializing, authError } = useAuth();

  if (isInitializing) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (authError) {
    return (
      <Box sx={{ p: 3 }}>
        <Alert severity="error">
          {authError === 'network' ? 'Ошибка соединения с сервером.' : 'Произошла ошибка авторизации.'}
          <button onClick={() => window.location.reload()} style={{ marginLeft: '10px' }}>
            Обновить
          </button>
        </Alert>
      </Box>
    );
  }
  
  if (isAuthenticated) {
    return children;
  }

  return null; 
};

export default ProtectedRoute;
import React from 'react';
import { useAuth } from '../context/AuthContext';
import { 
  Box, 
  CircularProgress, 
  Alert, 
  Button, 
  Paper,
  Typography,
  Stack 
} from '@mui/material';
import { Refresh, Login } from '@mui/icons-material';

const ProtectedRoute = ({ children }) => {
  const { isAuthenticated, authError, fetchUser, login } = useAuth();

  if (authError) {
    return (
      <Box 
        sx={{ 
          display: 'flex', 
          justifyContent: 'center', 
          alignItems: 'center', 
          height: '100vh',
          p: 3 
        }}
      >
        <Paper elevation={3} sx={{ p: 4, maxWidth: 500, textAlign: 'center' }}>
          <Stack spacing={3}>
            <Alert 
              severity="error" 
              sx={{ textAlign: 'left' }}
            >
              {authError === 'network'
                ? 'Ошибка соединения с сервером. Проверьте подключение к интернету.'
                : 'Произошла ошибка при загрузке данных пользователя.'}
            </Alert>
            
            <Stack direction="row" spacing={2} justifyContent="center">
              <Button 
                variant="contained" 
                startIcon={<Refresh />}
                onClick={fetchUser}
              >
                Повторить попытку
              </Button>
              
              {authError === 'network' && (
                <Button 
                  variant="outlined" 
                  startIcon={<Login />}
                  onClick={login}
                >
                  Войти заново
                </Button>
              )}
            </Stack>
          </Stack>
        </Paper>
      </Box>
    );
  }
 
  if (!isAuthenticated) {
    return (
      <Box 
        sx={{ 
          display: 'flex', 
          justifyContent: 'center', 
          alignItems: 'center', 
          height: '100vh' 
        }}
      >
        <Paper elevation={3} sx={{ p: 4, textAlign: 'center' }}>
          <Stack spacing={3}>
            <Typography variant="h5">
              Требуется авторизация
            </Typography>
            <Typography variant="body1" color="text.secondary">
              Для доступа к этой странице необходимо войти в систему
            </Typography>
            <Button 
              variant="contained" 
              size="large"
              startIcon={<Login />}
              onClick={login}
            >
              Войти в систему
            </Button>
          </Stack>
        </Paper>
      </Box>
    );
  }

  return children;
};

export default ProtectedRoute;
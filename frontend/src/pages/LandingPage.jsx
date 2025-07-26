import React from 'react';
import {
  Box,
  Button,
  Container,
  Typography,
  AppBar,
  Toolbar,
  keyframes,
} from '@mui/material';
import BugReportIcon from '@mui/icons-material/BugReport';
import { useAuth } from '../context/AuthContext.jsx';
import { useNavigate } from 'react-router-dom';
import { useTheme } from '@mui/material/styles';

const fadeIn = keyframes`
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
`;
 
const LandingPage = () => {
  const { isAuthenticated, login } = useAuth();
  const navigate = useNavigate();
  const theme = useTheme();
  
  const handleStart = () => {
    if (isAuthenticated) {
      navigate('/dashboard');
    } else {
      login();
    }
  };
 
  return (
    <Box sx={{ position: 'relative', minHeight: '100vh' }}>
      <AppBar position="absolute" color="transparent" elevation={0} sx={{ zIndex: 1 }}>
        <Toolbar>
          <BugReportIcon sx={{ mr: 1, color: 'primary.main' }} />
          <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
            BugTracker
          </Typography>
          {!isAuthenticated && (
            <Button color="primary" variant="outlined" onClick={login}>
              Войти
            </Button>
          )}
        </Toolbar>
      </AppBar>
 
      <Container
        component="main"
        sx={{
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          textAlign: 'center',
          py: { xs: 4, sm: 8 },
        }}
      >
        <Typography
          variant="h2"
          component="h1"
          gutterBottom
          sx={{
            fontWeight: 'bold',
            fontSize: { xs: '2.5rem', sm: '3.75rem' },
            animation: `${fadeIn} 1s ease-out`,
            color: 'text.primary',
          }}
        >
          Эффективное отслеживание задач
        </Typography>
 
        <Typography
          variant="h5"
          color="text.secondary"
          paragraph
          sx={{
            maxWidth: '600px',
            mb: 4,
            fontSize: { xs: '1rem', sm: '1.5rem' },
            animation: `${fadeIn} 1s ease-out 0.2s`,
            animationFillMode: 'backwards',
          }}
        >
          Инструмент для регистрации, отслеживания и управления задачами в рамках проектной работы.
        </Typography>
 
        <Button
          variant="contained"
          size="large"
          color="primary"
          onClick={handleStart}
          startIcon={<BugReportIcon />}
          sx={{
            animation: `${fadeIn} 1s ease-out 0.4s`,
            animationFillMode: 'backwards',
          }}
        >
          {isAuthenticated ? 'Перейти в панель' : 'Начать работу'}
        </Button>
      </Container>
    </Box>
  );
};
 
export default LandingPage;
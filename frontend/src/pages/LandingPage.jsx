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

  const handleStart = () => {
    if (isAuthenticated) {
      navigate('/dashboard');
    } else {
      login();
    }
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh', overflow: 'hidden' }}>
      <AppBar position="static" color="transparent" elevation={0}>
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
          flexGrow: 1,
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
            background: 'linear-gradient(90deg, #90caf9, #f48fb1)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
          }}
        >
          Систематизируйте хаос. Исправляйте баги.
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
          Наш BugTracker — это простое и мощное решение для отслеживания ошибок,
          управления задачами и совместной работы вашей команды.
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
import React, { useEffect } from 'react';
import { Routes, Route, useNavigate } from 'react-router-dom';
import LandingPage from '../pages/LandingPage.jsx';
import DashboardPage from '../pages/DashboardPage.jsx';
import LoginPage from '../pages/LoginPage.jsx';
import ProtectedRoute from './ProtectedRoute.jsx';
import { useAuth } from '../context/AuthContext.jsx';
import { Box, CircularProgress } from '@mui/material';
import NotFoundPage from '../pages/NotFoundPage.jsx';

const PostLoginHandler = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const redirectPath = sessionStorage.getItem('redirectAfterLogin') || '/dashboard';
    sessionStorage.removeItem('redirectAfterLogin');
    
    navigate(redirectPath, { replace: true });
  }, [navigate]);

  return (
    <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <CircularProgress />
    </Box>
  );
};


const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      
      <Route path="/login" element={<LoginPage />} />
      
      <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <DashboardPage />
          </ProtectedRoute>
        }
      />
      
      <Route path="/auth-callback-success" element={<PostLoginHandler />} />
      
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
};

export default AppRoutes;
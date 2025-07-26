import React, { useEffect } from 'react';
import { Routes, Route, useNavigate, Outlet } from 'react-router-dom';
import LandingPage from '../pages/LandingPage.jsx';
import DashboardPage from '../pages/DashboardPage.jsx';
import LoginPage from '../pages/LoginPage.jsx';
import ProtectedRoute from './ProtectedRoute.jsx';
import { Box, CircularProgress, Typography, Button } from '@mui/material';
import NotFoundPage from '../pages/NotFoundPage.jsx';
import Layout from '../components/Layout.jsx';
import FoldersPage from '../pages/FoldersPage.jsx';
import FolderTasksPage from '../pages/FolderTasksPage.jsx';
import TaskDetailsPage from '../pages/TaskDetailsPage.jsx';
import CreateTaskPage from '../pages/CreateTaskPage.jsx';
import MyTasksPage from '../pages/MyTasksPage.jsx';

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

const ProtectedLayout = () => {
  return (
    <ProtectedRoute>
      <Layout>
        <Outlet />
      </Layout>
    </ProtectedRoute>
  );
};

const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/auth-callback-success" element={<PostLoginHandler />} />
      
      <Route element={<ProtectedLayout />}>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/tasks" element={<MyTasksPage />} />
        <Route path="/folders" element={<FoldersPage />} />
        <Route path="/folders/:folderId/tasks" element={<FolderTasksPage />} /> 
        <Route path="/tasks/:taskId" element={<TaskDetailsPage />} />
        <Route path="/folders/:folderId/tasks/create" element={<CreateTaskPage />} /> 
      </Route>
      
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
};

export default AppRoutes;
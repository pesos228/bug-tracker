import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext.jsx';
import { Box, Typography, Container, Paper, Grid, CircularProgress, Alert } from '@mui/material';
import { getUserStats } from '../api/userService';

const DashboardPage = () => {
  const { user } = useAuth();
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        const data = await getUserStats();
        setStats(data);
      } catch (err) {
        setError('Не удалось загрузить статистику по задачам.');
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  return (
    <Container component="main" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>
        Добро пожаловать, {user?.FirstName}!
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
        Это ваша панель управления. Вот краткая сводка по вашим задачам.
      </Typography>

      {loading && <CircularProgress />}
      {error && <Alert severity="error">{error}</Alert>}
      
      {stats && (
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6}>
            <Paper sx={{ p: 3, textAlign: 'center' }}>
              <Typography variant="h2" color="primary">{stats.inProgressTasksCount}</Typography>
              <Typography variant="subtitle1" color="text.secondary">Задач в работе</Typography>
            </Paper>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Paper sx={{ p: 3, textAlign: 'center' }}>
              <Typography variant="h2" color="success.main">{stats.completedTasksCount}</Typography>
              <Typography variant="subtitle1" color="text.secondary">Завершено задач</Typography>
            </Paper>
          </Grid>
        </Grid>
      )}
    </Container>
  );
};

export default DashboardPage;
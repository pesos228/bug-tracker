import React from 'react';
import { Box, Typography, Button, Container } from '@mui/material';
import { Link } from 'react-router-dom';
import ReportProblemIcon from '@mui/icons-material/ReportProblem';

const NotFoundPage = () => {
  return (
    <Container component="main" maxWidth="sm">
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          textAlign: 'center',
          height: '100vh',
        }}
      >
        <Typography variant="h2" component="h1" gutterBottom sx={{ fontWeight: 'bold' }}>
          404
        </Typography>
        <Typography variant="h5" color="text.secondary" paragraph>
          Ой! Страница, которую вы ищете, не существует.
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
          Возможно, вы ошиблись в адресе или страница была перемещена.
        </Typography>
        <Button
          component={Link}
          to="/"
          variant="contained"
          size="large"
        >
          Вернуться на главную
        </Button>
      </Box>
    </Container>
  );
};

export default NotFoundPage;
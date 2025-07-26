import React from 'react';
import { Box, Typography, Button, Paper, Avatar, Stack, Divider, FormGroup, FormControlLabel, Switch, Tooltip } from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import { useAuth } from '../context/AuthContext';
import { useThemeContext } from '../context/ThemeContext';
import { useSnackbar } from 'notistack';

const KEYCLOAK_ACCOUNT_URL = import.meta.env.VITE_KEYCLOAK_ACCOUNT_URL;

const ProfilePage = () => {
  const { user, isAdmin, logout } = useAuth();
  const { mode, toggleTheme } = useThemeContext();
  const { enqueueSnackbar } = useSnackbar();

  if (!user) return <Typography>Загрузка...</Typography>;

  const handleEditProfile = () => {
    if (!KEYCLOAK_ACCOUNT_URL) {
      enqueueSnackbar('URL для редактирования профиля не настроен администратором.', { variant: 'error' });
      return;
    }
    window.open(KEYCLOAK_ACCOUNT_URL, '_blank');
    logout();
  };

  const isEditDisabled = !KEYCLOAK_ACCOUNT_URL;

  return (
    <Paper sx={{ p: { xs: 2, md: 4 }, maxWidth: 700, mx: 'auto' }}>
      <Stack spacing={4}>
        <Box>
          <Typography variant="h5" gutterBottom>Профиль</Typography>
          <Stack direction="row" spacing={3} alignItems="center">
            <Avatar sx={{ width: 80, height: 80, bgcolor: 'primary.main', fontSize: '2.5rem' }}>
              {user.FirstName.charAt(0)}{user.LastName.charAt(0)}
            </Avatar>
            <Box>
              <Typography variant="h6">{user.FirstName} {user.LastName}</Typography>
              <Typography variant="body1" color="text.secondary">
                {isAdmin ? 'Администратор' : 'Пользователь'}
              </Typography>
              
              <Tooltip title={isEditDisabled ? "Функция не настроена" : "Редактировать в Keycloak"}>
                <span>
                  <Button
                    variant="outlined"
                    size="small"
                    startIcon={<EditIcon />}
                    onClick={handleEditProfile}
                    sx={{ mt: 1 }}
                  >
                    Редактировать
                  </Button>
                </span>
              </Tooltip>
            </Box>
          </Stack>
        </Box>
        
        <Divider />

        <Box>
          <Typography variant="h5" gutterBottom>Настройки</Typography>
          <FormGroup>
            <FormControlLabel
              control={<Switch checked={mode === 'dark'} onChange={toggleTheme} />}
              label="Темная тема"
            />
          </FormGroup>
        </Box>
      </Stack>
    </Paper>
  );
};

export default ProfilePage;
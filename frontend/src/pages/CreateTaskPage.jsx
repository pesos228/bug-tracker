import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Paper, TextField, Button, Typography, Box, Stack, CircularProgress } from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import dayjs from 'dayjs';
import { useSnackbar } from 'notistack';
import { createTask } from '../api/taskService';
import UserSelectionDialog from '../components/selectors/UserSelectionDialog';
import SelectionInput from '../components/selectors/SelectionInput';

const CreateTaskPage = () => {
  const { folderId } = useParams();
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();

  const [formData, setFormData] = useState({
    softName: '',
    requestId: '',
    description: '',
    testEnvDateUpdate: null,
    assigneeId: '',
  });
  const [selectedUserName, setSelectedUserName] = useState('');
  const [isUserDialogOpen, setUserDialogOpen] = useState(false);
  const [isSubmitting, setSubmitting] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleDateChange = (newValue) => {
    setFormData(prev => ({ ...prev, testEnvDateUpdate: newValue ? dayjs(newValue).toISOString() : null }));
  };

  const handleUserSelect = (user) => {
    setFormData(prev => ({ ...prev, assigneeId: user.id }));
    setSelectedUserName(user.fullName);
  };
  
  const isFormValid = () => {
    return formData.softName && formData.requestId && formData.description && formData.testEnvDateUpdate && formData.assigneeId;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!isFormValid()) {
      enqueueSnackbar('Пожалуйста, заполните все обязательные поля', { variant: 'warning' });
      return;
    }
    
    setSubmitting(true);
    try {
      await createTask(folderId, formData);
      enqueueSnackbar('Задача успешно создана!', { variant: 'success' });
      navigate(`/folders/${folderId}/tasks`);
    } catch (error) {
      enqueueSnackbar(error.response?.data || 'Ошибка при создании задачи', { variant: 'error' });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <Paper sx={{ p: 3, maxWidth: '800px', mx: 'auto' }}>
        <Typography variant="h4" gutterBottom>
          Новая задача
        </Typography>
        <form onSubmit={handleSubmit}>
          <Stack spacing={3}>
            <TextField
              name="softName"
              label="ПО"
              value={formData.softName}
              onChange={handleChange}
              fullWidth
              required
            />
            <TextField
              name="requestId"
              label="Номер заявки"
              value={formData.requestId}
              onChange={handleChange}
              fullWidth
              required
            />
            <TextField
              name="description"
              label="Описание"
              value={formData.description}
              onChange={handleChange}
              fullWidth
              multiline
              rows={4}
              required
            />
            <DatePicker
              label="Дата обновления ТС"
              value={formData.testEnvDateUpdate ? dayjs(formData.testEnvDateUpdate) : null}
              onChange={handleDateChange}
              format="DD.MM.YYYY"
              sx={{ width: '100%' }}
              slotProps={{ textField: { required: true } }}
            />
            <SelectionInput
              label="Ответственный"
              value={selectedUserName}
              onClick={() => setUserDialogOpen(true)}
              required
            />
            <Box sx={{ textAlign: 'right' }}>
              <Button
                type="submit"
                variant="contained"
                disabled={isSubmitting || !isFormValid()}
              >
                {isSubmitting ? <CircularProgress size={24} /> : 'Создать задачу'}
              </Button>
            </Box>
          </Stack>
        </form>
      </Paper>

      <UserSelectionDialog
        open={isUserDialogOpen}
        onClose={() => setUserDialogOpen(false)}
        onSelect={handleUserSelect}
        currentValue={formData.assigneeId}
      />
    </>
  );
};

export default CreateTaskPage;
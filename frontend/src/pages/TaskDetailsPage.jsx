import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Paper, TextField, Button, Typography, CircularProgress, Alert, Box,
  Select, MenuItem, FormControl, InputLabel, Stack, Grid, Divider
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import dayjs from 'dayjs';
import { useSnackbar } from 'notistack';
import { getTaskDetails, updateTaskByAdmin, updateTaskByUser, deleteTask } from '../api/taskService';
import { useAuth } from '../context/AuthContext';
import { statusNameMapping } from '../utils/statusUtils';
import UserSelectionDialog from '../components/selectors/UserSelectionDialog';
import FolderSelectionDialog from '../components/selectors/FolderSelectionDialog';
import SelectionInput from '../components/selectors/SelectionInput';
import { searchUsers } from '../api/userService';
import { searchFolders } from '../api/folderService';
import ConfirmDialog from '../components/ConfirmDialog';
import PageBreadcrumbs from '../components/PageBreadcrumbs';

const resultNameMapping = { success: 'Успешно', failure: 'Неуспешно', warning: 'Есть замечания' };

const TaskDetailsPage = () => {
  const { taskId } = useParams();
  const navigate = useNavigate();
  const { isAdmin } = useAuth();
  const { enqueueSnackbar } = useSnackbar();

  const [initialTask, setInitialTask] = useState(null);
  const [formData, setFormData] = useState({});
  const [loading, setLoading] = useState(true);
  const [isSubmitting, setSubmitting] = useState(false);
  const [error, setError] = useState(null);
  const [selectedUserName, setSelectedUserName] = useState('');
  const [selectedFolderName, setSelectedFolderName] = useState('');
  
  const [isUserDialogOpen, setUserDialogOpen] = useState(false);
  const [isFolderDialogOpen, setFolderDialogOpen] = useState(false);
  const [isConfirmDeleteDialogOpen, setConfirmDeleteDialogOpen] = useState(false);

  const isAdminView = isAdmin;
  
  useEffect(() => {
    const fetchDetails = async () => {
      setLoading(true);
      setError(null);
      try {
        const view = isAdmin ? 'full' : 'short';
        const data = await getTaskDetails(taskId, view);
        
        const isZeroDate = (dateStr) => dateStr && dateStr.startsWith('0001-01-01');
        if (isZeroDate(data.testEnvDateUpdate)) data.testEnvDateUpdate = null;
        if (isZeroDate(data.checkDate)) data.checkDate = null;

        setInitialTask(data);
        setFormData(data);

        if (data.assigneeID) {
          const userData = await searchUsers(1, 1, data.assigneeID);
          if (userData.data.length > 0) setSelectedUserName(userData.data[0].fullName);
        }
        if (data.folderID) {
          const folderData = await searchFolders(1, 1, data.folderID);
          if (folderData.data.length > 0) setSelectedFolderName(folderData.data[0].name);
        }
      } catch (err) {
        setError('Не удалось загрузить задачу.');
        enqueueSnackbar('Ошибка при загрузке задачи', { variant: 'error' });
      } finally {
        setLoading(false);
      }
    };
    fetchDetails();
  }, [taskId, isAdmin, enqueueSnackbar]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };
  
  const handleDateChange = (name, newValue) => {
    setFormData(prev => ({ ...prev, [name]: newValue ? dayjs(newValue).toISOString() : null }));
  };

  const handleUserSelect = (user) => {
    setFormData(prev => ({ ...prev, assigneeID: user.id }));
    setSelectedUserName(user.fullName);
  };
  
  const handleFolderSelect = (folder) => {
    setFormData(prev => ({ ...prev, folderID: folder.id }));
    setSelectedFolderName(folder.name);
  };
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    const changedData = {};
    for (const key in formData) {
      if (JSON.stringify(formData[key]) !== JSON.stringify(initialTask[key])) {
        changedData[key] = formData[key];
      }
    }

    if (Object.keys(changedData).length === 0) {
      enqueueSnackbar('Нет изменений для сохранения', { variant: 'info' });
      return;
    }

    setSubmitting(true);
    try {
      if (isAdminView) {
        await updateTaskByAdmin(taskId, changedData);
      } else {
        const userUpdateData = {
          checkStatus: changedData.checkStatus,
          checkResult: changedData.checkResult,
          comment: changedData.comment,
        };
        await updateTaskByUser(taskId, userUpdateData);
      }
      enqueueSnackbar('Задача успешно обновлена!', { variant: 'success' });
      const view = isAdminView ? 'full' : 'short';
      const updatedData = await getTaskDetails(taskId, view);
      setInitialTask(updatedData);
      setFormData(updatedData);
    } catch (err) {
      enqueueSnackbar('Ошибка при обновлении задачи', { variant: 'error' });
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async () => {
    setConfirmDeleteDialogOpen(false);
    setSubmitting(true);
    try {
      await deleteTask(taskId);
      enqueueSnackbar('Задача успешно удалена', { variant: 'success' });
      navigate(-1);
    } catch (error) {
      enqueueSnackbar('Ошибка при удалении задачи', { variant: 'error' });
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) return <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}><CircularProgress /></Box>;
  if (error) return <Alert severity="error" sx={{ m: 2 }}>{error}</Alert>;

  const isDirty = JSON.stringify(initialTask) !== JSON.stringify(formData);

  const renderAdminView = () => (
    <form onSubmit={handleSubmit}>
      <Stack spacing={2}>
        <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
          <TextField name="softName" label="ПО" value={formData.softName || ''} onChange={handleChange} fullWidth />
          <TextField name="requestID" label="Номер заявки" value={formData.requestID || ''} onChange={handleChange} fullWidth />
          <DatePicker label="Дата обновления ТС" value={formData.testEnvDateUpdate ? dayjs(formData.testEnvDateUpdate) : null} onChange={(val) => handleDateChange('testEnvDateUpdate', val)} sx={{ width: '100%' }} format="DD.MM.YYYY" />
          <DatePicker label="Дата проверки" value={formData.checkDate ? dayjs(formData.checkDate) : null} onChange={(val) => handleDateChange('checkDate', val)} sx={{ width: '100%' }} format="DD.MM.YYYY" />
        </Stack>
        <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
          <TextField name="description" label="Описание" value={formData.description || ''} onChange={handleChange} fullWidth multiline rows={4} />
          <TextField name="comment" label="Комментарий" value={formData.comment || ''} onChange={handleChange} fullWidth multiline rows={4} />
        </Stack>
        <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems="center">
          <FormControl fullWidth><InputLabel>Статус проверки</InputLabel><Select name="checkStatus" value={formData.checkStatus || ''} label="Статус проверки" onChange={handleChange}><MenuItem value=""><em>Не выбрано</em></MenuItem>{Object.entries(statusNameMapping).map(([key, name]) => (<MenuItem key={key} value={key}>{name}</MenuItem>))}</Select></FormControl>
          <FormControl fullWidth><InputLabel>Результат проверки</InputLabel><Select name="checkResult" value={formData.checkResult || ''} label="Результат проверки" onChange={handleChange}><MenuItem value=""><em>Не выбрано</em></MenuItem>{Object.entries(resultNameMapping).map(([key, name]) => (<MenuItem key={key} value={key}>{name}</MenuItem>))}</Select></FormControl>
          <SelectionInput label="Ответственный" value={selectedUserName} onClick={() => setUserDialogOpen(true)} />
          <SelectionInput label="Папка" value={selectedFolderName} onClick={() => setFolderDialogOpen(true)} />
        </Stack>
        <Stack direction="row" justifyContent="space-between" sx={{ mt: 2 }}>
            <Button variant="outlined" color="error" onClick={() => setConfirmDeleteDialogOpen(true)} disabled={isSubmitting}>Удалить задачу</Button>
            <Button type="submit" variant="contained" disabled={isSubmitting || !isDirty}>{isSubmitting ? <CircularProgress size={24} /> : 'Сохранить'}</Button>
        </Stack>
      </Stack>
    </form>
  );

  const renderUserView = () => (
    <form onSubmit={handleSubmit}>
      <Stack spacing={3}>
        <Box>
          <Typography variant="h6" gutterBottom>Информация о задаче</Typography>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <Typography variant="subtitle2" color="text.secondary">ПО</Typography>
              <Typography>{formData.softName}</Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="subtitle2" color="text.secondary">Номер заявки</Typography>
              <Typography>{formData.requestID}</Typography>
            </Grid>
            <Grid item xs={12}>
              <Typography variant="subtitle2" color="text.secondary">Описание</Typography>
              <Typography sx={{ whiteSpace: 'pre-wrap', mt: 0.5, wordBreak: 'break-word' }}>{formData.description}</Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="subtitle2" color="text.secondary">Дата обновления ТС</Typography>
              <Typography>{formData.testEnvDateUpdate ? dayjs(formData.testEnvDateUpdate).format('DD.MM.YYYY') : 'Не указана'}</Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="subtitle2" color="text.secondary">Дата проверки</Typography>
              <Typography>{formData.checkDate ? dayjs(formData.checkDate).format('DD.MM.YYYY') : 'Не указана'}</Typography>
            </Grid>
          </Grid>
        </Box>
        <Divider />
  
        <Box>
          <Typography variant="h6" gutterBottom>Ваш отзыв</Typography>
          <Stack spacing={2}>
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
              <FormControl fullWidth>
                <InputLabel>Статус проверки</InputLabel>
                <Select name="checkStatus" value={formData.checkStatus || ''} label="Статус проверки" onChange={handleChange}>
                  <MenuItem value=""><em>Не выбрано</em></MenuItem>
                  {Object.entries(statusNameMapping).map(([key, name]) => (<MenuItem key={key} value={key}>{name}</MenuItem>))}
                </Select>
              </FormControl>
              <FormControl fullWidth>
                <InputLabel>Результат проверки</InputLabel>
                <Select name="checkResult" value={formData.checkResult || ''} label="Результат проверки" onChange={handleChange}>
                  <MenuItem value=""><em>Не выбрано</em></MenuItem>
                  {Object.entries(resultNameMapping).map(([key, name]) => (<MenuItem key={key} value={key}>{name}</MenuItem>))}
                </Select>
              </FormControl>
            </Stack>
            <TextField name="comment" label="Комментарий" value={formData.comment || ''} onChange={handleChange} fullWidth multiline rows={4} />
            <Box sx={{ textAlign: 'right' }}>
              <Button type="submit" variant="contained" disabled={isSubmitting || !isDirty}>{isSubmitting ? <CircularProgress size={24} /> : 'Отправить отзыв'}</Button>
            </Box>
          </Stack>
        </Box>
      </Stack>
    </form>
  );
  
  return (
    <Box>
      <PageBreadcrumbs 
        items={[{ label: 'Задачи', to: -1 }]}
        currentPage={formData.requestID || 'Детали'}
      />
      <Paper sx={{ p: 3 }}>
        <Typography variant="h4" gutterBottom>
          {isAdminView ? 'Редактирование задачи' : 'Детали задачи'}
        </Typography>
        {isAdminView ? renderAdminView() : renderUserView()}
      </Paper>
      <UserSelectionDialog open={isUserDialogOpen} onClose={() => setUserDialogOpen(false)} onSelect={handleUserSelect} currentValue={formData.assigneeID} />
      <FolderSelectionDialog open={isFolderDialogOpen} onClose={() => setFolderDialogOpen(false)} onSelect={handleFolderSelect} currentValue={formData.folderID} />
      <ConfirmDialog
        open={isConfirmDeleteDialogOpen}
        onCancel={() => setConfirmDeleteDialogOpen(false)}
        onConfirm={handleDelete}
        title="Удалить задачу?"
        content="Вы уверены, что хотите удалить эту задачу? Это действие необратимо."
      />
    </Box>
  );
};

export default TaskDetailsPage;
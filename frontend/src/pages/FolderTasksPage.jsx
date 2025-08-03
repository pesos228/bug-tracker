import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box, Typography, Button, TextField, Paper, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, TablePagination, CircularProgress, Alert, InputAdornment, useTheme,
  useMediaQuery, Select, MenuItem, FormControl, InputLabel, Chip, Stack, LinearProgress,
  IconButton, Tooltip
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import SearchIcon from '@mui/icons-material/Search';
import DeleteIcon from '@mui/icons-material/Delete';
import FileDownloadIcon from '@mui/icons-material/FileDownload';
import { useSnackbar } from 'notistack';
import { useDebounce } from '../hooks/useDebounce';
import { getTasksByFolderId } from '../api/taskService';
import { getFolderDetails, deleteFolder, exportFolder } from '../api/folderService';
import { getStatusName, statusNameMapping, getStatusChipColor } from '../utils/statusUtils';
import ConfirmDialog from '../components/ConfirmDialog';
import dayjs from 'dayjs';
import PageBreadcrumbs from '../components/PageBreadcrumbs';

const truncateText = (text, maxLength) => {
  if (text.length <= maxLength) {
    return text;
  }
  return text.substring(0, maxLength) + '...';
};
  
const FolderTasksPage = () => {
  const { folderId } = useParams();
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  
  const [folderDetails, setFolderDetails] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [isInitialLoading, setInitialLoading] = useState(true);
  const [isTableLoading, setTableLoading] = useState(false);
  
  const [error, setError] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const debouncedSearchQuery = useDebounce(searchQuery, 500);
  const [isConfirmOpen, setConfirmOpen] = useState(false);
  
  const fetchTasks = useCallback(async () => {
    setTableLoading(true);
    setError(null);
    try {
      const tasksData = await getTasksByFolderId(folderId, {
        page: page + 1,
        pageSize: rowsPerPage,
        checkStatus: statusFilter,
        requestID: debouncedSearchQuery,
      });
      setTasks(tasksData.data || []);
      setTotalCount(tasksData.pagination.totalCount || 0);
    } catch (err) {
      setError('Не удалось загрузить задачи.');
    } finally {
      setTableLoading(false);
    }
  }, [folderId, page, rowsPerPage, statusFilter, debouncedSearchQuery]);
  
  useEffect(() => {
    const fetchInitialData = async () => {
      setInitialLoading(true);
      try {
        const detailsData = await getFolderDetails(folderId);
        setFolderDetails(detailsData);
      } catch (err) {
        setError('Не удалось загрузить данные папки.');
        enqueueSnackbar('Не удалось загрузить данные папки.', { variant: 'error' });
        navigate('/folders');
      } finally {
        setInitialLoading(false);
      }
    };
    fetchInitialData();
  }, [folderId, navigate, enqueueSnackbar]);
  
  useEffect(() => {
    if (folderDetails) {
      fetchTasks();
    }
  }, [page, rowsPerPage, statusFilter, debouncedSearchQuery, folderDetails, fetchTasks]);
  
  const handleSearchChange = (e) => setSearchQuery(e.target.value);
  const handleStatusChange = (e) => {
    setPage(0);
    setStatusFilter(e.target.value);
  };
  const handleChangePage = (e, newPage) => setPage(newPage);
  const handleChangeRowsPerPage = (e) => {
    setRowsPerPage(parseInt(e.target.value, 10));
    setPage(0);
  };
  const handleRowClick = (taskId) => navigate(`/tasks/${taskId}`, { state: { fromFolder: true } });
  
  const handleDelete = async () => {
    setConfirmOpen(false);
    try {
      await deleteFolder(folderId);
      enqueueSnackbar('Папка успешно удалена', { variant: 'success' });
      navigate('/folders');
    } catch (error) {
      enqueueSnackbar('Ошибка при удалении папки', { variant: 'error' });
    }
  };
  
  const handleExport = async () => {
    try {
        const response = await exportFolder(folderId);
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        const contentDisposition = response.headers['content-disposition'];
        let fileName = 'report.xlsx';
        if (contentDisposition) {
            const fileNameMatch = contentDisposition.match(/filename="(.+)"/);
            if (fileNameMatch && fileNameMatch.length === 2)
                fileName = fileNameMatch[1];
        }
        link.href = url;
        link.setAttribute('download', fileName);
        document.body.appendChild(link);
        link.click();
        link.remove();
        enqueueSnackbar('Отчет успешно скачан', { variant: 'success' });
    } catch (error) {
        enqueueSnackbar('Ошибка при скачивании отчета', { variant: 'error' });
    }
  };
  
  if (isInitialLoading) {
    return <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '50vh' }}><CircularProgress /></Box>;
  }
  
  if (error && !folderDetails) {
    return <Alert severity="error" sx={{ m: 2 }}>{error}</Alert>;
  }
  
  return (
    <Box>
      <PageBreadcrumbs 
        items={[{ label: 'Папки', to: '/folders' }]}
        currentPage={folderDetails?.name || 'Задачи'}
      />
      <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2, flexWrap: 'wrap', gap: 2 }}>
          <Box>
            <Typography variant="h4" component="h1">{folderDetails?.name}</Typography>
            {folderDetails && (<Typography variant="caption" color="text.secondary">Создана: {dayjs(folderDetails.createdAt).format('DD.MM.YYYY')} | Ответственный: {folderDetails.assigneePerson}</Typography>)}
          </Box>
          <Stack direction="row" spacing={1} alignItems="center">
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => navigate(`/folders/${folderId}/tasks/create`)}
              size={isMobile ? 'small' : 'medium'}
            >
              {!isMobile && 'Создать задачу'}
            </Button>
            <Tooltip title="Экспорт в Excel">
              {isMobile ? (
                <IconButton onClick={handleExport} size="small">
                  <FileDownloadIcon />
                </IconButton>
              ) : (
                <Button variant="outlined" startIcon={<FileDownloadIcon />} onClick={handleExport}>
                  Экспорт
                </Button>
              )}
            </Tooltip>
            <Tooltip title="Удалить папку">
              {isMobile ? (
                <IconButton onClick={() => setConfirmOpen(true)} size="small" color="error">
                  <DeleteIcon />
                </IconButton>
              ) : (
                <Button variant="outlined" color="error" startIcon={<DeleteIcon />} onClick={() => setConfirmOpen(true)}>
                  Удалить папку
                </Button>
              )}
            </Tooltip>
          </Stack>
        </Box>
    
        <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} sx={{ mb: 2 }}>
          <TextField fullWidth variant="outlined" placeholder="Поиск по номеру заявки..." value={searchQuery} onChange={handleSearchChange} InputProps={{ startAdornment: <InputAdornment position="start"><SearchIcon /></InputAdornment> }} />
          <FormControl fullWidth>
            <InputLabel>Статус</InputLabel>
            <Select value={statusFilter} label="Статус" onChange={handleStatusChange}>
              <MenuItem value=""><em>Все статусы</em></MenuItem>
              {Object.entries(statusNameMapping).map(([key, name]) => (<MenuItem key={key} value={key}>{name}</MenuItem>))}
            </Select>
          </FormControl>
        </Stack>
    
        <Box sx={{ position: 'relative' }}>
          {isTableLoading && <LinearProgress sx={{ position: 'absolute', top: 0, left: 0, width: '100%' }} />}
          <Box sx={{ opacity: isTableLoading ? 0.6 : 1, transition: 'opacity 0.3s' }}>
            {error && !isTableLoading && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
            <TableContainer>
              {!isMobile ? (
                <Table stickyHeader>
                  <TableHead><TableRow><TableCell sx={{ width: '20%' }}>ПО</TableCell><TableCell sx={{ width: '15%' }}>Номер заявки</TableCell><TableCell>Описание</TableCell><TableCell sx={{ width: '15%' }}>Дата создания</TableCell></TableRow></TableHead>
                  <TableBody>
                    {tasks.map((task) => (
                      <TableRow hover key={task.id} onClick={() => handleRowClick(task.id)} sx={{ cursor: 'pointer' }}>
                        <TableCell><Stack direction="row" alignItems="center" spacing={1}><Chip label={getStatusName(task.checkStatus)} size="small" variant="outlined" color={getStatusChipColor(task.checkStatus)}/><Typography variant="body2">{task.softName}</Typography></Stack></TableCell>
                        <TableCell>{task.requestId}</TableCell>
                        <TableCell>{truncateText(task.description, 100)}</TableCell>
                        <TableCell>{dayjs(task.createdAt).format('DD.MM.YYYY')}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              ) : (
                <Box>
                  {tasks.map((task) => (
                    <Paper key={task.id} sx={{ p: 2, mb: 2, cursor: 'pointer' }} onClick={() => handleRowClick(task.id)}>
                      <Stack direction="row" justifyContent="space-between" alignItems="center" mb={1}><Typography variant="h6">{task.softName}</Typography>
                      <Chip label={getStatusName(task.checkStatus)} size="small" variant="outlined" color={getStatusChipColor(task.checkStatus)}/></Stack>
                      <Typography variant="body2" color="text.secondary">Заявка: {task.requestId}</Typography>
                      <Typography variant="body1" noWrap sx={{ my: 1 }}>{truncateText(task.description, 50)}</Typography>
                      <Typography variant="caption" color="text.secondary">{dayjs(task.createdAt).format('DD.MM.YYYY')}</Typography>
                    </Paper>
                  ))}
                </Box>
              )}
            </TableContainer>
    
            <TablePagination
              rowsPerPageOptions={isMobile ? [] : [5, 10, 25]}
              component="div"
              count={totalCount}
              rowsPerPage={rowsPerPage}
              page={page}
              onPageChange={handleChangePage}
              onRowsPerPageChange={handleChangeRowsPerPage}
              labelRowsPerPage={isMobile ? '' : 'Строк на странице:'}
              labelDisplayedRows={({ from, to, count }) => `${from}-${to} из ${count}`}
            />
          </Box>
        </Box>
    
        <ConfirmDialog
          open={isConfirmOpen}
          title="Удалить папку?"
          content={`Вы уверены, что хотите удалить папку "${folderDetails?.name}"? Это действие необратимо.`}
          onConfirm={handleDelete}
          onCancel={() => setConfirmOpen(false)}
        />
      </Paper>
    </Box>
  );
};
  
export default FolderTasksPage;
import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box, Typography, TextField, Paper, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, TablePagination, CircularProgress, Alert, InputAdornment, useTheme,
  useMediaQuery, Select, MenuItem, FormControl, InputLabel, Chip, Stack, LinearProgress
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import { useSnackbar } from 'notistack';
import { useDebounce } from '../hooks/useDebounce';
import { getMyTasks } from '../api/taskService';
import { getStatusName, statusNameMapping, getStatusChipColor } from '../utils/statusUtils';
import dayjs from 'dayjs';

const MyTasksPage = () => {
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const [tasks, setTasks] = useState([]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('not_checked'); 
  const debouncedSearchQuery = useDebounce(searchQuery, 500);

  const fetchTasks = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const tasksData = await getMyTasks({
        page: page + 1,
        pageSize: rowsPerPage,
        checkStatus: statusFilter,
        requestID: debouncedSearchQuery,
      });
      setTasks(tasksData.data || []);
      setTotalCount(tasksData.pagination.totalCount || 0);
    } catch (err) {
      setError('Не удалось загрузить задачи.');
      enqueueSnackbar('Ошибка при загрузке задач', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  }, [page, rowsPerPage, statusFilter, debouncedSearchQuery, enqueueSnackbar]);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

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
  const handleRowClick = (taskId) => navigate(`/tasks/${taskId}`);

  return (
    <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
      <Typography variant="h4" component="h1" sx={{ mb: 2 }}>
        Мои задачи
      </Typography>

      <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} sx={{ mb: 2 }}>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Поиск по номеру заявки..."
          value={searchQuery}
          onChange={handleSearchChange}
          InputProps={{ startAdornment: <InputAdornment position="start"><SearchIcon /></InputAdornment> }}
        />
        <FormControl fullWidth>
          <InputLabel>Статус</InputLabel>
          <Select value={statusFilter} label="Статус" onChange={handleStatusChange}>
            <MenuItem value=""><em>Все статусы</em></MenuItem>
            {Object.entries(statusNameMapping).map(([key, name]) => (
              <MenuItem key={key} value={key}>{name}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Stack>

      <Box sx={{ position: 'relative' }}>
        {loading && <LinearProgress sx={{ position: 'absolute', top: 0, left: 0, width: '100%' }} />}
        <Box sx={{ opacity: loading ? 0.6 : 1, transition: 'opacity 0.3s' }}>
          {error && !loading && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
          
          <TableContainer>
            {!isMobile ? (
              <Table stickyHeader>
                <TableHead>
                  <TableRow>
                    <TableCell sx={{ width: '20%' }}>ПО</TableCell>
                    <TableCell sx={{ width: '15%' }}>Номер заявки</TableCell>
                    <TableCell>Описание</TableCell>
                    <TableCell sx={{ width: '15%' }}>Дата создания</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {tasks.map((task) => (
                    <TableRow hover key={task.id} onClick={() => handleRowClick(task.id)} sx={{ cursor: 'pointer' }}>
                      <TableCell>
                        <Stack direction="row" alignItems="center" spacing={1}>
                          <Chip label={getStatusName(task.checkStatus)} size="small" variant="outlined" color={getStatusChipColor(task.checkStatus)} />
                          <Typography variant="body2">{task.softName}</Typography>
                        </Stack>
                      </TableCell>
                      <TableCell>{task.requestId}</TableCell>
                      <TableCell>{task.description}</TableCell>
                      <TableCell>{dayjs(task.createdAt).format('DD.MM.YYYY')}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : (
              <Box>
                {tasks.map((task) => (
                  <Paper key={task.id} sx={{ p: 2, mb: 2, cursor: 'pointer' }} onClick={() => handleRowClick(task.id)}>
                    <Stack direction="row" justifyContent="space-between" alignItems="center" mb={1}>
                      <Typography variant="h6">{task.softName}</Typography>
                      <Chip label={getStatusName(task.checkStatus)} size="small" variant="outlined" color={getStatusChipColor(task.checkStatus)}/>
                    </Stack>
                    <Typography variant="body2" color="text.secondary">Заявка: {task.requestId}</Typography>
                    <Typography variant="body1" noWrap sx={{ my: 1 }}>{task.description}</Typography>
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
    </Paper>
  );
};

export default MyTasksPage;
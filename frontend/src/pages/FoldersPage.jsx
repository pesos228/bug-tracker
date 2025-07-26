import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box, Typography, Button, TextField, Paper, Table, TableBody, TableCell,
  TableContainer, TableHead, TableRow, TablePagination, CircularProgress,
  Alert, InputAdornment, useTheme, useMediaQuery,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import SearchIcon from '@mui/icons-material/Search';
import { useSnackbar } from 'notistack';
import { searchFolders } from '../api/folderService';
import { useDebounce } from '../hooks/useDebounce';
import CreateFolderDialog from '../components/CreateFolderDialog';
import dayjs from 'dayjs';

const FoldersPage = () => {
  const [folders, setFolders] = useState([]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [isCreateDialogOpen, setCreateDialogOpen] = useState(false);

  const debouncedSearchTerm = useDebounce(searchTerm, 500);
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const fetchFolders = useCallback(async (currentPage, currentRowsPerPage, search) => {
    setLoading(true);
    setError(null);
    try {
      const data = await searchFolders(currentPage + 1, currentRowsPerPage, search);
      setFolders(data.data || []);
      setTotalCount(data.pagination.totalCount || 0);
    } catch (err) {
      setError('Не удалось загрузить папки. Попробуйте обновить страницу.');
      enqueueSnackbar('Ошибка при загрузке данных', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  }, [enqueueSnackbar]);

  useEffect(() => {
    fetchFolders(page, rowsPerPage, debouncedSearchTerm);
  }, [page, rowsPerPage, debouncedSearchTerm, fetchFolders]);

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleSearchChange = (event) => {
    setSearchTerm(event.target.value);
    setPage(0);
  };

  const handleRowClick = (folderId) => {
    navigate(`/folders/${folderId}/tasks`);
  };
  
  const handleCreateSuccess = (newFolder) => {
    setFolders(prevFolders => [newFolder, ...prevFolders]);
    setTotalCount(prevCount => prevCount + 1);
  };

  return (
    <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h4" component="h1">
          Папки
        </Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => setCreateDialogOpen(true)}>
          Создать папку
        </Button>
      </Box>

      <TextField
        fullWidth
        variant="outlined"
        placeholder="Поиск по названию папки..."
        value={searchTerm}
        onChange={handleSearchChange}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon />
            </InputAdornment>
          ),
        }}
        sx={{ mb: 2 }}
      />

      {loading && <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}><CircularProgress /></Box>}
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
      
      {!loading && !error && (
        <>
          <TableContainer>
            {!isMobile ? (
              <Table stickyHeader aria-label="folders table">
                <TableHead>
                  <TableRow>
                    <TableCell>Название</TableCell>
                    <TableCell align="right">Количество задач</TableCell>
                    <TableCell align="right">Дата создания</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {folders.map((folder) => (
                    <TableRow hover key={folder.id} onClick={() => handleRowClick(folder.id)} sx={{ cursor: 'pointer' }}>
                      <TableCell component="th" scope="row">{folder.name}</TableCell>
                      <TableCell align="right">{folder.taskCount}</TableCell>
                      <TableCell align="right">{dayjs(folder.createdAt).format('DD.MM.YYYY')}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : (
              <Box>
                {folders.map((folder) => (
                  <Paper key={folder.id} sx={{ p: 2, mb: 2, cursor: 'pointer' }} onClick={() => handleRowClick(folder.id)}>
                    <Typography variant="h6" component="div" gutterBottom>{folder.name}</Typography>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <Typography variant="body2" color="text.secondary">Кол-во задач:</Typography>
                      <Typography variant="body1">{folder.taskCount}</Typography>
                    </Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mt: 1 }}>
                      <Typography variant="body2" color="text.secondary">Создана:</Typography>
                      <Typography variant="body1">{new Date(folder.createdAt).toLocaleDateString()}</Typography>
                    </Box>
                  </Paper>
                ))}
              </Box>
            )}
          </TableContainer>
          <TablePagination
            rowsPerPageOptions={[5, 10, 25]}
            component="div"
            count={totalCount}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handleChangePage}
            onRowsPerPageChange={handleChangeRowsPerPage}
            labelRowsPerPage="Строк на странице:"
            labelDisplayedRows={({ from, to, count }) => `${from}-${to} из ${count}`}
          />
        </>
      )}
      
      <CreateFolderDialog 
        open={isCreateDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onSuccess={handleCreateSuccess}
      />
    </Paper>
  );
};

export default FoldersPage;
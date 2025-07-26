import React, { useState, useEffect, useCallback } from 'react';
import { Dialog, DialogTitle, DialogContent, TextField, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, TablePagination, CircularProgress, Radio } from '@mui/material';
import { useDebounce } from '../../hooks/useDebounce';
import { searchFolders } from '../../api/folderService';
import dayjs from 'dayjs';

const FolderSelectionDialog = ({ open, onClose, onSelect, currentValue }) => {
  const [folders, setFolders] = useState([]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const debouncedSearchTerm = useDebounce(searchTerm, 500);

  const fetchFolders = useCallback(async () => {
    setLoading(true);
    try {
      const data = await searchFolders(page + 1, rowsPerPage, debouncedSearchTerm);
      setFolders(data.data || []);
      setTotalCount(data.pagination.totalCount || 0);
    } catch (error) {
      console.error("Failed to fetch folders", error);
    } finally {
      setLoading(false);
    }
  }, [page, rowsPerPage, debouncedSearchTerm]);

  useEffect(() => {
    if (open) {
      fetchFolders();
    }
  }, [open, fetchFolders]);

  const handleSelect = (folder) => {
    onSelect(folder);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle>Выберите папку</DialogTitle>
      <DialogContent>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Поиск по названию..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          sx={{ mb: 2 }}
        />
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell />
                <TableCell>Название</TableCell>
                <TableCell>Количество задач</TableCell>
                <TableCell>Дата создания</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow><TableCell colSpan={4} align="center"><CircularProgress /></TableCell></TableRow>
              ) : (
                folders.map((folder) => (
                  <TableRow hover key={folder.id} onClick={() => handleSelect(folder)} sx={{ cursor: 'pointer' }}>
                    <TableCell padding="checkbox">
                      <Radio checked={currentValue === folder.id} />
                    </TableCell>
                    <TableCell>{folder.name}</TableCell>
                    <TableCell>{folder.taskCount}</TableCell>
                    <TableCell>{dayjs(folder.createdAt).format('DD.MM.YYYY')}</TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          rowsPerPageOptions={[5, 10]}
          component="div"
          count={totalCount}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={(e, newPage) => setPage(newPage)}
          onRowsPerPageChange={(e) => { setRowsPerPage(parseInt(e.target.value, 10)); setPage(0); }}
          labelRowsPerPage="Строк на странице:"
          labelDisplayedRows={({ from, to, count }) => `${from}-${to} из ${count}`}
        />
      </DialogContent>
    </Dialog>
  );
};

export default FolderSelectionDialog;
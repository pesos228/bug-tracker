import React, { useState, useEffect, useCallback } from 'react';
import { Dialog, DialogTitle, DialogContent, TextField, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, TablePagination, CircularProgress, Radio } from '@mui/material';
import { useDebounce } from '../../hooks/useDebounce';
import { searchUsers } from '../../api/userService';

const UserSelectionDialog = ({ open, onClose, onSelect, currentValue }) => {
  const [users, setUsers] = useState([]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const debouncedSearchTerm = useDebounce(searchTerm, 500);

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    try {
      const data = await searchUsers(page + 1, rowsPerPage, debouncedSearchTerm);
      setUsers(data.data || []);
      setTotalCount(data.pagination.totalCount || 0);
    } catch (error) {
      console.error("Failed to fetch users", error);
    } finally {
      setLoading(false);
    }
  }, [page, rowsPerPage, debouncedSearchTerm]);

  useEffect(() => {
    if (open) {
      fetchUsers();
    }
  }, [open, fetchUsers]);

  const handleSelect = (user) => {
    onSelect(user);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle>Выберите пользователя</DialogTitle>
      <DialogContent>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Поиск по имени..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          sx={{ mb: 2 }}
        />
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell />
                <TableCell>Имя</TableCell>
                <TableCell>Задач в работе</TableCell>
                <TableCell>Завершено задач</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow><TableCell colSpan={4} align="center"><CircularProgress /></TableCell></TableRow>
              ) : (
                users.map((user) => (
                  <TableRow hover key={user.id} onClick={() => handleSelect(user)} sx={{ cursor: 'pointer' }}>
                    <TableCell padding="checkbox">
                      <Radio checked={currentValue === user.id} />
                    </TableCell>
                    <TableCell>{user.fullName}</TableCell>
                    <TableCell>{user.inProgressTasksCount}</TableCell>
                    <TableCell>{user.completedTasksCount}</TableCell>
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

export default UserSelectionDialog;
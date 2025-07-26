import React, { useState } from 'react';
import {
  Button,
  TextField,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  CircularProgress,
} from '@mui/material';
import { useSnackbar } from 'notistack';
import { createFolder } from '../api/folderService';

const CreateFolderDialog = ({ open, onClose, onSuccess }) => {
  const [folderName, setFolderName] = useState('');
  const [isCreating, setCreating] = useState(false);
  const { enqueueSnackbar } = useSnackbar();

  const handleCreate = async () => {
    if (!folderName.trim()) {
      enqueueSnackbar('Название папки не может быть пустым', { variant: 'warning' });
      return;
    }

    setCreating(true);
    try {
      const newFolder = await createFolder(folderName);
      enqueueSnackbar('Папка успешно создана!', { variant: 'success' });
      onSuccess(newFolder);
      handleClose();
    } catch (error) {
      enqueueSnackbar('Не удалось создать папку', { variant: 'error' });
    } finally {
      setCreating(false);
    }
  };

  const handleClose = () => {
    setFolderName('');
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} fullWidth maxWidth="sm">
      <DialogTitle>Создание новой папки</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          id="name"
          label="Название папки"
          type="text"
          fullWidth
          variant="standard"
          value={folderName}
          onChange={(e) => setFolderName(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleCreate()}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={isCreating}>Отмена</Button>
        <Button onClick={handleCreate} disabled={isCreating}>
          {isCreating ? <CircularProgress size={24} /> : 'Создать'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateFolderDialog;
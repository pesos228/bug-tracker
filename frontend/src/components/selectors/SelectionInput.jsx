import React from 'react';
import { TextField, InputAdornment, IconButton } from '@mui/material';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';

const SelectionInput = ({ label, value, onClick }) => {
  return (
    <TextField
      label={label}
      value={value || ''}
      onClick={onClick}
      readOnly
      fullWidth
      InputProps={{
        endAdornment: (
          <InputAdornment position="end">
            <IconButton size="small">
              <ArrowDropDownIcon />
            </IconButton>
          </InputAdornment>
        ),
      }}
    />
  );
};

export default SelectionInput;
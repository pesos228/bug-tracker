import React from 'react';
import { Breadcrumbs, Link, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';

const PageBreadcrumbs = ({ items, currentPage }) => {
  return (
    <Breadcrumbs
      separator={<NavigateNextIcon fontSize="small" />}
      aria-label="breadcrumb"
      sx={{ mb: 2 }}
    >
      {items.map((item, index) => (
        <Link
          key={index}
          component={RouterLink}
          underline="hover"
          color="inherit"
          to={item.to}
        >
          {item.label}
        </Link>
      ))}
      <Typography color="text.primary">{currentPage}</Typography>
    </Breadcrumbs>
  );
};

export default PageBreadcrumbs;
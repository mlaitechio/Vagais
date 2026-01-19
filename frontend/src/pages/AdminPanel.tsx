import React from 'react';
import { Box, Typography, useTheme } from '@mui/material';
import { AdminPanelSettings } from '@mui/icons-material';

const AdminPanel: React.FC = () => {
  const theme = useTheme();

  return (
    <Box sx={{ p: 3, minHeight: '100vh', background: theme.palette.background.default }}>
      <Box textAlign="center" py={8}>
        <AdminPanelSettings sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
        <Typography variant="h4" color="text.secondary" mb={1}>
          Admin Panel
        </Typography>
        <Typography variant="body1" color="text.secondary">
          This page will contain administrative controls and system management
        </Typography>
      </Box>
    </Box>
  );
};

export default AdminPanel;

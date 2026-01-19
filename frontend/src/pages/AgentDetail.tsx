import React from 'react';
import { Box, Typography, useTheme } from '@mui/material';
import { SmartToy } from '@mui/icons-material';

const AgentDetail: React.FC = () => {
  const theme = useTheme();

  return (
    <Box sx={{ p: 3, minHeight: '100vh', background: theme.palette.background.default }}>
      <Box textAlign="center" py={8}>
        <SmartToy sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
        <Typography variant="h4" color="text.secondary" mb={1}>
          Agent Detail Page
        </Typography>
        <Typography variant="body1" color="text.secondary">
          This page will show detailed information about a specific AI agent
        </Typography>
      </Box>
    </Box>
  );
};

export default AgentDetail;

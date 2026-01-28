import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  Button,
  Avatar,
  Chip,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  IconButton,
  useTheme,
  CircularProgress,
} from '@mui/material';
import {
  SmartToy,
  PlayArrow,
  Edit,
  Delete,
  Visibility,
  Download,
  Star,
  Search,
  Add,
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import { useQuery } from '@tanstack/react-query';
import apiService from '../services/api';
import { Agent } from '../types/api';

const MyAgents: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('');
  const [sortBy, setSortBy] = useState('latest');

  const { data: agents, isLoading } = useQuery({
    queryKey: ['userAgents', { search: searchQuery, category: categoryFilter, sort: sortBy }],
    queryFn: () => apiService.getAgents(1, 50, { status: 'published' }),
    retry: false,
  });

  const AgentCard = ({ agent }: { agent: Agent }) => (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.3 }}
      whileHover={{ scale: 1.02 }}
    >
      <Card
        sx={{
          background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
          borderRadius: 3,
          height: '100%',
          position: 'relative',
          overflow: 'hidden',
          cursor: 'pointer',
          transition: 'all 0.3s ease',
          '&:hover': {
            border: '1px solid rgba(152, 23, 126, 0.5)',
            background: 'linear-gradient(135deg, rgba(152, 23, 126, 0.15) 0%, rgba(0, 212, 255, 0.1) 100%)',
          },
          '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            height: '2px',
            background: `linear-gradient(90deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
          },
        }}
        onClick={() => navigate(`/agents/${agent.id}`)}
      >
        <CardContent>
          <Box display="flex" alignItems="flex-start" justifyContent="space-between" mb={2}>
            <Box display="flex" alignItems="center" gap={2} flex={1}>
              <Avatar
                src={agent.icon}
                sx={{
                  width: 48,
                  height: 48,
                  background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                }}
              >
                <SmartToy />
              </Avatar>
              <Box>
                <Typography variant="h6" fontWeight="bold">
                  {agent.name}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {agent.category}
                </Typography>
              </Box>
            </Box>
            <Box display="flex" gap={0.5}>
              <IconButton size="small" onClick={(e) => {
                e.stopPropagation();
                navigate(`/agents/${agent.id}/edit`);
              }}>
                <Edit fontSize="small" />
              </IconButton>
              <IconButton size="small" onClick={(e) => e.stopPropagation()}>
                <Delete fontSize="small" />
              </IconButton>
            </Box>
          </Box>

          <Typography variant="body2" color="text.secondary" mb={2} sx={{ minHeight: '40px' }}>
            {agent.description}
          </Typography>

          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Box display="flex" alignItems="center" gap={3}>
              <Box display="flex" alignItems="center" gap={0.5}>
                <Star sx={{ color: 'warning.main', fontSize: 16 }} />
                <Typography variant="body2" fontWeight="medium">{agent.rating}</Typography>
              </Box>
              <Box display="flex" alignItems="center" gap={0.5}>
                <Visibility sx={{ fontSize: 16 }} />
                <Typography variant="body2">{agent.usage_count}</Typography>
              </Box>
              <Box display="flex" alignItems="center" gap={0.5}>
                <Download sx={{ fontSize: 16 }} />
                <Typography variant="body2">{agent.downloads}</Typography>
              </Box>
            </Box>
          </Box>

          <Box display="flex" gap={1} flexWrap="wrap" mb={2}>
            {(() => {
              let tags: string[] = [];
              if (agent.tags) {
                if (Array.isArray(agent.tags)) {
                  tags = agent.tags;
                } else if (typeof agent.tags === 'string') {
                  try {
                    const decoded = atob(agent.tags);
                    const parsed = JSON.parse(decoded);
                    tags = Array.isArray(parsed) ? parsed : [];
                  } catch {
                    const tagStr = agent.tags as string;
                    tags = tagStr.includes(',') ? tagStr.split(',').map((t: string) => t.trim()) : [tagStr];
                  }
                }
              }
              return tags.slice(0, 3).map((tag, index) => (
                <Chip
                  key={index}
                  label={tag}
                  size="small"
                  variant="outlined"
                  sx={{ fontSize: '0.7rem' }}
                />
              ));
            })()}
          </Box>

          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="h6" color="primary" fontWeight="bold">
              {agent.price === 0 ? 'Free' : `$${agent.price.toFixed(2)}`}
            </Typography>
            <Button
              variant="contained"
              size="small"
              startIcon={<PlayArrow />}
              onClick={(e) => {
                e.stopPropagation();
                navigate(`/chat/${agent.id}`);
              }}
              sx={{
                background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                borderRadius: 2,
              }}
            >
              Execute
            </Button>
          </Box>
        </CardContent>
      </Card>
    </motion.div>
  );

  return (
    <Box sx={{ p: 3, minHeight: '100vh', background: theme.palette.background.default }}>
      {/* Header */}
      <Box display="flex" alignItems="center" justifyContent="space-between" mb={4}>
        <Box>
          <Typography variant="h4" fontWeight="bold" mb={1}>
            My Agents
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Manage and monitor all your AI agents
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={() => navigate('/agents/create')}
          sx={{
            background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
            borderRadius: 2,
            px: 3,
          }}
        >
          Create Agent
        </Button>
      </Box>

      {/* Filters */}
      <Card
        sx={{
          background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
          borderRadius: 3,
          mb: 4,
        }}
      >
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={12} md={4}>
              <TextField
                fullWidth
                placeholder="Search agents..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                InputProps={{
                  startAdornment: <Search sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
                sx={{
                  '& .MuiOutlinedInput-root': {
                    borderRadius: 2,
                  },
                }}
              />
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControl fullWidth>
                <InputLabel>Category</InputLabel>
                <Select
                  value={categoryFilter}
                  label="Category"
                  onChange={(e) => setCategoryFilter(e.target.value)}
                  sx={{ borderRadius: 2 }}
                >
                  <MenuItem value="">All Categories</MenuItem>
                  <MenuItem value="text-generation">Text Generation</MenuItem>
                  <MenuItem value="code-generation">Code Generation</MenuItem>
                  <MenuItem value="data-analysis">Data Analysis</MenuItem>
                  <MenuItem value="content-creation">Content Creation</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={4}>
              <FormControl fullWidth>
                <InputLabel>Sort By</InputLabel>
                <Select
                  value={sortBy}
                  label="Sort By"
                  onChange={(e) => setSortBy(e.target.value)}
                  sx={{ borderRadius: 2 }}
                >
                  <MenuItem value="latest">Latest</MenuItem>
                  <MenuItem value="rating">Highest Rated</MenuItem>
                  <MenuItem value="popular">Most Popular</MenuItem>
                  <MenuItem value="trending">Trending</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Agents Grid */}
      {isLoading ? (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={3}>
          {(agents as any)?.data?.map((agent: Agent) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={agent.id}>
              <AgentCard agent={agent} />
            </Grid>
          ))}
        </Grid>
      )}

      {!isLoading && (!agents || (agents as any)?.data?.length === 0) && (
        <Box textAlign="center" py={8}>
          <SmartToy sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
          <Typography variant="h5" color="text.secondary" mb={1}>
            No agents found
          </Typography>
          <Typography variant="body2" color="text.secondary" mb={3}>
            Create your first agent to get started
          </Typography>
          <Button
            variant="contained"
            startIcon={<Add />}
            onClick={() => navigate('/agents/create')}
            sx={{
              background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
              borderRadius: 2,
            }}
          >
            Create Agent
          </Button>
        </Box>
      )}
    </Box>
  );
};

export default MyAgents;

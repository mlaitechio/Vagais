import React from 'react';
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  Button,
  Avatar,
  Chip,
  LinearProgress,
  Divider,
  useTheme,
  Container,
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  Person,
  SmartToy,
  PlayArrow,
  Add,
  Download,
  Star,
  AttachMoney,
  Assessment,
  Speed,
  Security,
  Dashboard as DashboardIcon,
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import apiService from '../services/api';
import { Agent } from '../types/api';
import { useCategories } from '../hooks/useCategories';
import { useAuth } from '../contexts/AuthContext';

const DemoDashboard: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const { user } = useAuth();

  // Fetch marketplace data for authenticated user
  const { data: agents, isLoading: agentsLoading } = useQuery({
    queryKey: ['marketplaceAgents'],
    queryFn: async () => {
      const response = await apiService.searchAgents({ page: 1, limit: 6 });
      return response;
    },
    retry: false,
  });

  const { data: categories } = useCategories();

  // Get top categories from backend data
  const topCategories = categories && typeof categories === 'object' 
    ? Object.keys(categories as Record<string, number>).slice(0, 4)
    : ['Customer facing', 'Sales', 'Research', 'Document processing']; // Fallback

  // Mock public statistics for demo
  const publicStats = {
    totalAgents: 1250,
    totalUsers: 8500,
    totalExecutions: 125000,
    totalRevenue: 450000,
    growthRate: 15.3,
    successRate: 94.2,
    avgResponseTime: 2.3,
    topCategories,
  };

  const StatCard = ({ title, value, icon, trend, color = 'primary' }: any) => {
    const colorKey = color as 'primary' | 'secondary' | 'error' | 'warning' | 'info' | 'success';
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
      <Card
        sx={{
          height: '100%',
          background: `linear-gradient(135deg, ${theme.palette[colorKey].main}15, ${theme.palette[colorKey].main}05)`,
          border: `1px solid ${theme.palette[colorKey].main}20`,
          '&:hover': {
            transform: 'translateY(-2px)',
            boxShadow: theme.shadows[8],
            transition: 'all 0.3s ease',
          },
        }}
      >
        <CardContent>
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Avatar
              sx={{
                bgcolor: `${theme.palette[colorKey].main}20`,
                color: theme.palette[colorKey].main,
              }}
            >
              {icon}
            </Avatar>
            {trend && (
              <Chip
                icon={trend > 0 ? <TrendingUp /> : <TrendingDown />}
                label={`${Math.abs(trend)}%`}
                size="small"
                color={trend > 0 ? 'success' : 'error'}
                variant="outlined"
              />
            )}
          </Box>
          <Typography variant="h4" component="div" fontWeight="bold" color="text.primary">
            {value}
          </Typography>
          <Typography variant="body2" color="text.secondary" mt={1}>
            {title}
          </Typography>
        </CardContent>
      </Card>
    </motion.div>
    );
  };

  const AgentCard = ({ agent }: { agent: Agent }) => (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.3 }}
    >
      <Card
        sx={{
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          '&:hover': {
            transform: 'translateY(-4px)',
            boxShadow: theme.shadows[8],
            transition: 'all 0.3s ease',
          },
        }}
      >
        <Box
          sx={{
            height: 120,
            background: `linear-gradient(135deg, ${theme.palette.primary.main}20, ${theme.palette.secondary.main}20)`,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            position: 'relative',
          }}
        >
          <SmartToy sx={{ fontSize: 48, color: theme.palette.primary.main }} />
          <Chip
            label={agent.pricing_model === 'free' ? 'Free' : `$${agent.price}`}
            size="small"
            color={agent.pricing_model === 'free' ? 'success' : 'primary'}
            sx={{ position: 'absolute', top: 8, right: 8 }}
          />
        </Box>
        <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
          <Typography variant="h6" component="div" fontWeight="bold" gutterBottom>
            {agent.name}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2, flexGrow: 1 }}>
            {agent.description}
          </Typography>
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={1}>
            <Box display="flex" alignItems="center">
              <Star sx={{ fontSize: 16, color: theme.palette.warning.main, mr: 0.5 }} />
              <Typography variant="body2" color="text.secondary">
                {agent.rating} ({agent.review_count})
              </Typography>
            </Box>
            <Chip label={agent.category} size="small" variant="outlined" />
          </Box>
          <Box display="flex" alignItems="center" justifyContent="space-between">
            <Typography variant="body2" color="text.secondary">
              {agent.usage_count} uses
            </Typography>
            <Button
              variant="contained"
              size="small"
              startIcon={<Download />}
              sx={{ minWidth: 'auto' }}
            >
              Try
            </Button>
          </Box>
        </CardContent>
      </Card>
    </motion.div>
  );

  return (
    <Container maxWidth="xl" sx={{ py: 4 }}>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={4}>
          <Box>
            <Typography variant="h3" component="h1" fontWeight="bold" gutterBottom>
              vagais.ai Platform Overview
            </Typography>
            <Typography variant="h6" color="text.secondary">
              Public marketplace statistics and trending agents
            </Typography>
          </Box>
          <Box display="flex" gap={2}>
            <Button
              variant="outlined"
              startIcon={<DashboardIcon />}
              onClick={() => window.location.href = '/demo-marketplace'}
            >
              View Marketplace
            </Button>
            <Button
              variant="contained"
              startIcon={<Add />}
              onClick={() => window.location.href = '/register'}
            >
              Get Started
            </Button>
          </Box>
        </Box>

        {/* Public Statistics */}
        <Grid container spacing={3} mb={4}>
          <Grid item xs={12} sm={6} md={3}>
            <StatCard
              title="Total Agents"
              value={publicStats.totalAgents.toLocaleString()}
              icon={<SmartToy />}
              trend={12.5}
              color="primary"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <StatCard
              title="Active Users"
              value={publicStats.totalUsers.toLocaleString()}
              icon={<Person />}
              trend={8.7}
              color="secondary"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <StatCard
              title="Total Executions"
              value={publicStats.totalExecutions.toLocaleString()}
              icon={<PlayArrow />}
              trend={15.3}
              color="success"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <StatCard
              title="Platform Revenue"
              value={`$${(publicStats.totalRevenue / 1000).toFixed(0)}K`}
              icon={<AttachMoney />}
              trend={22.1}
              color="warning"
            />
          </Grid>
        </Grid>

        {/* Performance Metrics */}
        <Grid container spacing={3} mb={4}>
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Platform Performance
                </Typography>
                <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Success Rate
                  </Typography>
                  <Typography variant="body2" fontWeight="bold">
                    {publicStats.successRate}%
                  </Typography>
                </Box>
                <LinearProgress
                  variant="determinate"
                  value={publicStats.successRate}
                  sx={{ mb: 2 }}
                />
                <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                  <Typography variant="body2" color="text.secondary">
                    Avg Response Time
                  </Typography>
                  <Typography variant="body2" fontWeight="bold">
                    {publicStats.avgResponseTime}s
                  </Typography>
                </Box>
                <LinearProgress
                  variant="determinate"
                  value={(publicStats.avgResponseTime / 5) * 100}
                  sx={{ mb: 2 }}
                />
                <Box display="flex" justifyContent="space-between" alignItems="center">
                  <Typography variant="body2" color="text.secondary">
                    Growth Rate
                  </Typography>
                  <Typography variant="body2" fontWeight="bold">
                    {publicStats.growthRate}%
                  </Typography>
                </Box>
                <LinearProgress
                  variant="determinate"
                  value={publicStats.growthRate}
                  color="success"
                />
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Top Categories
                </Typography>
                <Box display="flex" flexWrap="wrap" gap={1}>
                  {publicStats.topCategories.map((category, index) => (
                    <Chip
                      key={index}
                      label={category}
                      variant="outlined"
                      color="primary"
                      size="small"
                    />
                  ))}
                </Box>
                <Divider sx={{ my: 2 }} />
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  Platform Highlights
                </Typography>
                <Box display="flex" flexDirection="column" gap={1}>
                  <Box display="flex" alignItems="center" gap={1}>
                    <Security sx={{ fontSize: 16, color: 'success.main' }} />
                    <Typography variant="body2">Enterprise-grade security</Typography>
                  </Box>
                  <Box display="flex" alignItems="center" gap={1}>
                    <Speed sx={{ fontSize: 16, color: 'primary.main' }} />
                    <Typography variant="body2">High-performance infrastructure</Typography>
                  </Box>
                  <Box display="flex" alignItems="center" gap={1}>
                    <Assessment sx={{ fontSize: 16, color: 'secondary.main' }} />
                    <Typography variant="body2">Advanced analytics & insights</Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Trending Agents */}
        <Box mb={4}>
          <Typography variant="h5" component="h2" gutterBottom fontWeight="bold">
            Trending Agents
          </Typography>
          <Typography variant="body1" color="text.secondary" mb={3}>
            Most popular agents in the marketplace
          </Typography>
          <Grid container spacing={3}>
            {agentsLoading ? (
              Array.from({ length: 6 }).map((_, index) => (
                <Grid item xs={12} sm={6} md={4} key={index}>
                  <Card sx={{ height: 300 }}>
                    <CardContent>
                      <Box sx={{ height: 200, bgcolor: 'grey.100' }} />
                      <Typography variant="h6" sx={{ mt: 2 }}>Loading...</Typography>
                    </CardContent>
                  </Card>
                </Grid>
              ))
            ) : (
              (agents as any)?.data?.slice(0, 6).map((agent: Agent) => (
                <Grid item xs={12} sm={6} md={4} key={agent.id}>
                  <AgentCard agent={agent} />
                </Grid>
              ))
            )}
          </Grid>
        </Box>

        {/* Call to Action */}
        <Card
          sx={{
            background: `linear-gradient(135deg, ${theme.palette.primary.main}10, ${theme.palette.secondary.main}10)`,
            border: `1px solid ${theme.palette.primary.main}20`,
          }}
        >
          <CardContent sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="h4" component="h2" gutterBottom fontWeight="bold">
              Ready to Get Started?
            </Typography>
            <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
              Join thousands of developers and businesses using vagais.ai to build and deploy AI agents
            </Typography>
            <Box display="flex" gap={2} justifyContent="center">
              {!user ? (
                <>
                  <Button
                    variant="contained"
                    size="large"
                    onClick={() => navigate('/register')}
                  >
                    Create Account
                  </Button>
                  <Button
                    variant="outlined"
                    size="large"
                    onClick={() => navigate('/marketplace')}
                  >
                    Explore Marketplace
                  </Button>
                </>
              ) : (
                <>
                  <Button
                    variant="contained"
                    size="large"
                    onClick={() => navigate('/marketplace')}
                  >
                    Browse Marketplace
                  </Button>
                  <Button
                    variant="outlined"
                    size="large"
                    onClick={() => navigate('/agents')}
                  >
                    My Agents
                  </Button>
                </>
              )}
            </Box>
          </CardContent>
        </Card>
      </motion.div>
    </Container>
  );
};

export default DemoDashboard;

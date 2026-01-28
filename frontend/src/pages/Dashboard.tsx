import React, { useState, useEffect } from 'react';
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
  LinearProgress,
  IconButton,
  Badge,
  useTheme,
} from '@mui/material';
import {
  TrendingUp,
  TrendingDown,
  SmartToy,
  PlayArrow,
  Settings,
  Notifications,
  Add,
  Mic,
  Visibility,
  Download,
  Star,
  Timeline,
  Assessment,
  Speed,
  Security,
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import { useQuery } from '@tanstack/react-query';
import apiService from '../services/api';
import { Agent, User } from '../types/api';

const Dashboard: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);

  // Fetch dashboard data
  const { data: stats } = useQuery({
    queryKey: ['dashboardStats'],
    queryFn: () => apiService.getUserStats(),
    retry: false,

  });

  const { data: agents } = useQuery({
    queryKey: ['userAgents'],
    queryFn: () => apiService.getAgents(1, 10, { status: 'published' }),
    retry: false,
  });

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await apiService.getProfile();
        setUser(userData);
      } catch (error) {
        console.error('Failed to fetch user:', error);
        // Set mock user data for demo
        setUser({
          id: 'demo-user',
          email: 'demo@example.com',
          avatar: '',
          role: 'user',
          organization_id: 'demo-org',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        } as any);
      }
    };
    fetchUser();
  }, []);

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
          background: 'linear-gradient(135deg, rgba(152, 23, 126, 0.1) 0%, rgba(0, 212, 255, 0.1) 100%)',
          border: '1px solid rgba(152, 23, 126, 0.2)',
          borderRadius: 3,
          height: '100%',
        }}
      >
        <CardContent>
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Box
              sx={{
                p: 1,
                borderRadius: 2,
                background: `linear-gradient(135deg, ${theme.palette[colorKey].main}20, ${theme.palette[colorKey].light}20)`,
              }}
            >
              {icon}
            </Box>
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
          <Typography variant="h4" fontWeight="bold" gutterBottom>
            {value}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {title}
          </Typography>
        </CardContent>
      </Card>
    </motion.div>
    );
  };

  const AgentCard = ({ agent }: { agent: Agent }) => (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.3 }}
    >
      <Card
        sx={{
          background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
          borderRadius: 3,
          height: '100%',
          position: 'relative',
          overflow: 'hidden',
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
      >
        <CardContent>
          <Box display="flex" alignItems="center" mb={2}>
            <Avatar
              src={agent.icon}
              sx={{
                width: 48,
                height: 48,
                mr: 2,
                background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
              }}
            >
              <SmartToy />
            </Avatar>
            <Box flex={1}>
              <Typography variant="h6" fontWeight="bold">
                {agent.name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {agent.category}
              </Typography>
            </Box>
            <IconButton size="small">
              <PlayArrow />
            </IconButton>
          </Box>

          <Typography variant="body2" color="text.secondary" mb={2}>
            {agent.description}
          </Typography>

          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Box display="flex" alignItems="center">
              <Star sx={{ color: 'warning.main', fontSize: 16, mr: 0.5 }} />
              <Typography variant="body2">{agent.rating}</Typography>
            </Box>
            <Box display="flex" alignItems="center">
              <Visibility sx={{ fontSize: 16, mr: 0.5 }} />
              <Typography variant="body2">{agent.usage_count}</Typography>
            </Box>
            <Box display="flex" alignItems="center">
              <Download sx={{ fontSize: 16, mr: 0.5 }} />
              <Typography variant="body2">{agent.downloads}</Typography>
            </Box>
          </Box>

          <Box display="flex" gap={1} flexWrap="wrap">
            {(() => {
              let tags: string[] = [];
              if (agent.tags) {
                if (Array.isArray(agent.tags)) {
                  tags = agent.tags;
                } else if (typeof agent.tags === 'string') {
                  try {
                    // Handle base64 encoded JSON string
                    const decoded = atob(agent.tags);
                    const parsed = JSON.parse(decoded);
                    tags = Array.isArray(parsed) ? parsed : [];
                  } catch {
                    // If parsing fails, try to split by comma or treat as single tag
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

          <Box mt={2} display="flex" justifyContent="space-between" alignItems="center">
            <Typography variant="h6" color="primary">
              {agent.price === 0 ? 'Free' : `$${agent.price}`}
            </Typography>
            <Button
              variant="contained"
              size="small"
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
      {/* Header Section */}
      <Box display="flex" alignItems="center" justifyContent="space-between" mb={4}>
        <Box display="flex" alignItems="center" gap={2}>
          <Avatar
            sx={{
              width: 56,
              height: 56,
              background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
            }}
          >
            {user?.first_name?.[0] || 'U'}
          </Avatar>
          <Box>
            <Typography variant="h4" fontWeight="bold">
              Hey, {user?.first_name || 'User'}! 👋
            </Typography>
            <Typography variant="body1" color="text.secondary">
              Ready to explore AI agents? Just ask me anything!
            </Typography>
          </Box>
        </Box>
        <Box display="flex" alignItems="center" gap={2}>
          <IconButton
            sx={{
              background: 'rgba(255, 255, 255, 0.1)',
              border: '1px solid rgba(255, 255, 255, 0.2)',
            }}
          >
            <Mic />
          </IconButton>
          <Badge badgeContent={3} color="error">
            <IconButton>
              <Notifications />
            </IconButton>
          </Badge>
        </Box>
      </Box>

      {/* Stats Cards */}
      <Grid container spacing={3} mb={4}>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Total Agents"
            value={stats?.total_agents || 0}
            icon={<SmartToy color="primary" />}
            trend={12}
            color="primary"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Executions"
            value={stats?.total_executions || 0}
            icon={<PlayArrow color="secondary" />}
            trend={8}
            color="secondary"
          />
        </Grid>
      </Grid>

      {/* Main Content Grid */}
      <Grid container spacing={3}>
        {/* Left Column */}
        <Grid item xs={12} lg={8}>
          {/* Active Agents Section */}
          <Card
            sx={{
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
              borderRadius: 3,
              mb: 3,
            }}
          >
            <CardContent>
              <Box display="flex" alignItems="center" justifyContent="space-between" mb={3}>
                <Typography variant="h5" fontWeight="bold">
                  Your Active Agents
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

              <Grid container spacing={3}>
                {(agents as any)?.data?.slice(0, 6).map((agent: Agent) => (
                  <Grid item xs={12} sm={6} md={4} key={agent.id}>
                    <AgentCard agent={agent} />
                  </Grid>
                ))}
              </Grid>
            </CardContent>
          </Card>

          {/* Analytics Section */}
          <Card
            sx={{
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
              borderRadius: 3,
            }}
          >
            <CardContent>
              <Typography variant="h5" fontWeight="bold" mb={3}>
                Performance Analytics
              </Typography>
              
              <Grid container spacing={3}>
                <Grid item xs={12} md={6}>
                  <Box mb={3}>
                    <Typography variant="h6" mb={2}>
                      Usage Trends
                    </Typography>
                    <Box display="flex" alignItems="center" gap={2} mb={1}>
                      <Typography variant="body2">This Month</Typography>
                      <LinearProgress
                        variant="determinate"
                        value={75}
                        sx={{
                          flex: 1,
                          height: 8,
                          borderRadius: 4,
                          background: 'rgba(255, 255, 255, 0.1)',
                          '& .MuiLinearProgress-bar': {
                            background: `linear-gradient(90deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                          },
                        }}
                      />
                      <Typography variant="body2">75%</Typography>
                    </Box>
                  </Box>
                </Grid>
                
                <Grid item xs={12} md={6}>
                  <Box mb={3}>
                    <Typography variant="h6" mb={2}>
                      Revenue Growth
                    </Typography>
                    <Box display="flex" alignItems="center" gap={2}>
                      <Timeline color="primary" />
                      <Typography variant="h4" color="primary" fontWeight="bold">
                        +23.5%
                      </Typography>
                    </Box>
                  </Box>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Right Column */}
        <Grid item xs={12} lg={4}>
          {/* Quick Actions */}
          <Card
            sx={{
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
              borderRadius: 3,
              mb: 3,
            }}
          >
            <CardContent>
              <Typography variant="h6" fontWeight="bold" mb={3}>
                Quick Actions
              </Typography>
              
              <Box display="flex" flexDirection="column" gap={2}>
                <Button
                  variant="outlined"
                  fullWidth
                  startIcon={<Add />}
                  onClick={() => navigate('/agents/create')}
                  sx={{ borderRadius: 2, py: 1.5 }}
                >
                  Create New Agent
                </Button>
                <Button
                  variant="outlined"
                  fullWidth
                  startIcon={<Assessment />}
                  sx={{ borderRadius: 2, py: 1.5 }}
                >
                  View Analytics
                </Button>
                <Button
                  variant="outlined"
                  fullWidth
                  startIcon={<Settings />}
                  sx={{ borderRadius: 2, py: 1.5 }}
                >
                  Manage Integrations
                </Button>
              </Box>
            </CardContent>
          </Card>

          {/* System Status */}
          <Card
            sx={{
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
              borderRadius: 3,
              mb: 3,
            }}
          >
            <CardContent>
              <Typography variant="h6" fontWeight="bold" mb={3}>
                System Status
              </Typography>
              
              <Box display="flex" flexDirection="column" gap={2}>
                <Box display="flex" alignItems="center" justifyContent="space-between">
                  <Box display="flex" alignItems="center" gap={1}>
                    <Security color="success" />
                    <Typography variant="body2">System Security</Typography>
                  </Box>
                  <Chip label="Active" color="success" size="small" />
                </Box>
                
                <Box display="flex" alignItems="center" justifyContent="space-between">
                  <Box display="flex" alignItems="center" gap={1}>
                    <Speed color="primary" />
                    <Typography variant="body2">Performance</Typography>
                  </Box>
                  <Chip label="Optimal" color="primary" size="small" />
                </Box>
                
                <Box display="flex" alignItems="center" justifyContent="space-between">
                  <Box display="flex" alignItems="center" gap={1}>
                    <SmartToy color="secondary" />
                    <Typography variant="body2">AI Models</Typography>
                  </Box>
                  <Chip label="Online" color="secondary" size="small" />
                </Box>
              </Box>
            </CardContent>
          </Card>

          {/* Recent Activity */}
          <Card
            sx={{
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
              borderRadius: 3,
            }}
          >
            <CardContent>
              <Typography variant="h6" fontWeight="bold" mb={3}>
                Recent Activity
              </Typography>
              
              <Box display="flex" flexDirection="column" gap={2}>
                {[
                  { action: 'Agent executed', time: '2 min ago', color: 'success' as const },
                  { action: 'New subscription', time: '1 hour ago', color: 'primary' as const },
                  { action: 'Payment received', time: '3 hours ago', color: 'warning' as const },
                ].map((activity, index) => {
                  const colorKey = activity.color as 'primary' | 'secondary' | 'error' | 'warning' | 'info' | 'success';
                  return (
                  <Box key={index} display="flex" alignItems="center" gap={2}>
                    <Avatar
                      sx={{
                        width: 32,
                        height: 32,
                        background: `linear-gradient(135deg, ${theme.palette[colorKey].main}, ${theme.palette[colorKey].light})`,
                      }}
                    >
                      {activity.action[0]}
                    </Avatar>
                    <Box flex={1}>
                      <Typography variant="body2" fontWeight="medium">
                        {activity.action}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {activity.time}
                      </Typography>
                    </Box>
                  </Box>
                  );
                })}
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard;

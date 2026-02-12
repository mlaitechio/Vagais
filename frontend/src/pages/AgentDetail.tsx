import  {  useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Container,
  Card,
  CardContent,
  Typography,
  Button,
  Avatar,
  Chip,
  Grid,
  Rating,
  Divider,
  Tab,
  Tabs,
  CircularProgress,
  useTheme,
  IconButton,
} from '@mui/material';
import {
  SmartToy,
  PlayArrow,
  ArrowBack,

  Visibility,
  Download,
  Share,
  FavoriteBorder,
  Favorite,
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import { useQuery } from '@tanstack/react-query';
import apiService from '../services/api';
import { Agent } from '../types/api';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`agent-tabpanel-${index}`}
      aria-labelledby={`agent-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  );
}

const AgentDetail: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [tabValue, setTabValue] = useState(0);
  const [isFavorite, setIsFavorite] = useState(false);

  const { data: agent, isLoading } = useQuery({
    queryKey: ['agent', id],
    queryFn: () => id ? apiService.getAgent(id) : Promise.reject('No ID'),
    retry: false,
    enabled: !!id,
  });

  // Helper function to convert YouTube URL to embed format
  const getYouTubeEmbedUrl = (url: string): string => {
    if (!url) return '';
    
    // Already an embed URL
    if (url.includes('youtube.com/embed/')) {
      return url;
    }
    
    // Extract video ID from various YouTube URL formats
    let videoId = '';
    
    // Handle youtu.be format
    if (url.includes('youtu.be/')) {
      videoId = url.split('youtu.be/')[1].split('?')[0].split('&')[0];
    }
    // Handle youtube.com/watch?v= format
    else if (url.includes('youtube.com/watch?v=')) {
      videoId = url.split('watch?v=')[1].split('&')[0];
    }
    // Handle youtube.com/v/ format
    else if (url.includes('youtube.com/v/')) {
      videoId = url.split('youtube.com/v/')[1].split('?')[0].split('&')[0];
    }
    
    // Return embed URL if video ID was found, otherwise return original URL
    return videoId ? `https://www.youtube.com/embed/${videoId}` : url;
  };

  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="100vh">
        <CircularProgress />
      </Box>
    );
  }

  if (!agent) {
    return (
      <Box sx={{ p: 3, minHeight: '100vh', background: theme.palette.background.default }}>
        <Button startIcon={<ArrowBack />} onClick={() => navigate('/marketplace')}>
          Back to Agents
        </Button>
        <Box textAlign="center" py={8}>
          <Typography variant="h5" color="text.secondary" mb={1}>
            Agent not found
          </Typography>
        </Box>
      </Box>
    );
  }

  const agentData = agent as unknown as Agent;

  return (
    <Box sx={{ minHeight: '100vh', background: theme.palette.background.default }}>
      {/* Header with Back Button */}
      <Box sx={{ p: 3, borderBottom: '1px solid rgba(255, 255, 255, 0.1)' }}>
        <Button
          startIcon={<ArrowBack />}
          onClick={() => navigate('/marketplace')}
          sx={{ mb: 2 }}
        >
          Back to Agents
        </Button>
      </Box>

      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Grid container spacing={4}>
          {/* Left Column - Agent Info */}
          <Grid className="agd_left" item xs={12} md={4}>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5 }}
            >
              <Card
                sx={{
                  background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
                  border: '1px solid rgba(255, 255, 255, 0.1)',
                  borderRadius: 3,
                  position: 'sticky',
                  top: 100,
                }}
              >
                <CardContent>
                  {/* Agent Avatar */}
                  <Box textAlign="center" mb={3}>
                    <Avatar
                      src={agentData.icon}
                      sx={{
                        width: 120,
                        height: 120,
                        mx: 'auto',
                        mb: 2,
                        background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                      }}
                    >
                      <SmartToy sx={{ fontSize: 60 }} />
                    </Avatar>
                    <Typography variant="h5" fontWeight="bold" mb={1}>
                      {agentData.name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" mb={2}>
                      {agentData.category}
                    </Typography>
                    <Chip
                      label={agentData.status}
                      color={agentData.status === 'published' ? 'success' : 'warning'}
                      size="small"
                      sx={{ mb: 2 }}
                    />
                  </Box>

                  <Divider sx={{ my: 2, borderColor: 'rgba(255, 255, 255, 0.1)' }} />

                  {/* Rating */}
                  <Box mb={3}>
                    <Box display="flex" alignItems="center" justifyContent="center" gap={1} mb={1}>
                      <Rating
                        value={agentData.rating}
                        readOnly
                        precision={0.1}
                        size="small"
                      />
                      <Typography variant="body2" color="text.secondary">
                        {agentData.rating.toFixed(1)}
                      </Typography>
                    </Box>
                    <Typography variant="caption" color="text.secondary" align="center" display="block">
                      ({agentData.review_count} reviews)
                    </Typography>
                  </Box>

                  {/* Stats */}
                  <Grid container spacing={2} mb={3}>
                    <Grid item xs={6}>
                      <Box textAlign="center">
                        <Box display="flex" justifyContent="center" mb={1}>
                          <Visibility color="primary" />
                        </Box>
                        <Typography variant="h6" fontWeight="bold">
                          {agentData.usage_count}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          Uses
                        </Typography>
                      </Box>
                    </Grid>
                    <Grid item xs={6}>
                      <Box textAlign="center">
                        <Box display="flex" justifyContent="center" mb={1}>
                          <Download color="primary" />
                        </Box>
                        <Typography variant="h6" fontWeight="bold">
                          {agentData.downloads}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          Downloads
                        </Typography>
                      </Box>
                    </Grid>
                  </Grid>

                  <Divider sx={{ my: 2, borderColor: 'rgba(255, 255, 255, 0.1)' }} />

                  {/* Price */}
                  <Box textAlign="center" mb={3}>
                    <Typography variant="body2" color="text.secondary" mb={1}>
                      Pricing
                    </Typography>
                    <Typography variant="h4" color="primary" fontWeight="bold">
                      {agentData.price === 0 ? 'Free' : `$${agentData.price.toFixed(2)}`}
                    </Typography>
                  </Box>

                  {/* Action Buttons */}
                  <Box display="flex" flexDirection="column" gap={2}>
                    <Button
                      variant="contained"
                      fullWidth
                      size="large"
                      startIcon={<PlayArrow />}
                      onClick={() => navigate(`/chat/${agentData.id}`)}
                      sx={{
                        background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                        borderRadius: 2,
                        py: 1.5,
                      }}
                    >
                      Execute Agent
                    </Button>
                    <Box display="flex" gap={1}>
                      <IconButton
                        
                       
                        onClick={() => setIsFavorite(!isFavorite)}
                        sx={{
                          borderRadius: 2,
			  flex: 1,
                          border: '1px solid rgba(255, 255, 255, 0.2)',
                        }}
                      >
                        {isFavorite ? <Favorite color="error" /> : <FavoriteBorder />}
                      </IconButton>
                      <IconButton
      	
                        sx={{
                          borderRadius: 2,
     			flex: 1,	
     			  border: '1px solid rgba(255, 255, 255, 0.2)',
                        }}
                      >
                        <Share />
                      </IconButton>
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            </motion.div>
          </Grid>

          {/* Right Column - Details */}
          <Grid className="ags_right" item xs={12} md={8}>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.1 }}
            >
              <Card
                sx={{
                  background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
                  border: '1px solid rgba(255, 255, 255, 0.1)',
                  borderRadius: 3,
                  mb: 3,
                }}
              >
                <CardContent>
                  {/* Tabs */}
                  <Tabs
                    value={tabValue}
                    onChange={(_, newValue) => setTabValue(newValue)}
                    aria-label="agent details"
                    sx={{
                      borderBottom: '1px solid rgba(255, 255, 255, 0.1)',
                      mb: 3,
                    }}
                  >
                    <Tab label="Overview" />
                    <Tab label="How It Works" />
                    <Tab label="Documentation" />
                    <Tab label="Reviews" />
                  </Tabs>

                  {/* Overview Tab */}
                  <TabPanel value={tabValue} index={0}>
                    <Typography variant="h6" fontWeight="bold" mb={2}>
                      Description
                    </Typography>
                    <Typography variant="body2" color="text.secondary" paragraph>
                      {agentData.description}
                    </Typography>

                    {/* Tags */}
                    {agentData.tags && agentData.tags.length > 0 && (
                      <Box mt={3}>
                        <Typography variant="h6" fontWeight="bold" mb={2}>
                          Tags
                        </Typography>
                        <Box display="flex" gap={1} flexWrap="wrap">
                          {(() => {
                            let tags: string[] = [];
                            if (Array.isArray(agentData.tags)) {
                              tags = agentData.tags;
                            } else if (typeof agentData.tags === 'string') {
                              try {
                                const decoded = atob(agentData.tags);
                                const parsed = JSON.parse(decoded);
                                tags = Array.isArray(parsed) ? parsed : [];
                              } catch {
                                const tagStr = agentData.tags as string;
                                tags = tagStr.includes(',') ? tagStr.split(',').map((t: string) => t.trim()) : [tagStr];
                              }
                            }
                            return tags.map((tag, index) => (
                              <Chip
                                key={index}
                                label={tag}
                                variant="outlined"
                              />
                            ));
                          })()}
                        </Box>
                      </Box>
                    )}
                  </TabPanel>

                  {/* How It Works Tab */}
                  <TabPanel value={tabValue} index={1}>
                    {agentData.how_it_works ? (
                      <>
                        <Typography variant="h6" fontWeight="bold" mb={2}>
                          How It Works
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph>
                          {agentData.how_it_works}
                        </Typography>
                      </>
                    ) : (
                      <Typography variant="body2" color="text.secondary">
                        No detailed information available yet.
                      </Typography>
                    )}

                    {/* Video Section */}
                    {agentData.video_url && (
                      <Box mt={4}>
                        <Typography variant="h6" fontWeight="bold" mb={2}>
                          Demo Video
                        </Typography>
                        <Box
                          component="iframe"
                          src={getYouTubeEmbedUrl(agentData.video_url)}
                          width="100%"
                          height="400"
                          frameBorder="0"
                          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                          allowFullScreen
                          sx={{ borderRadius: 2, mt: 2 }}
                        />
                      </Box>
                    )}
                  </TabPanel>

                  {/* Documentation Tab */}
                  <TabPanel value={tabValue} index={2}>
                    {agentData.documentation ? (
                      <Box>
                        <Typography variant="h6" fontWeight="bold" mb={2}>
                          Documentation
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph>
                          {agentData.documentation}
                        </Typography>
                      </Box>
                    ) : (
                      <Typography variant="body2" color="text.secondary">
                        No documentation available.
                      </Typography>
                    )}
                  </TabPanel>

                  {/* Reviews Tab */}
                  <TabPanel value={tabValue} index={3}>
                    <Typography variant="h6" fontWeight="bold" mb={2}>
                      User Reviews
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {agentData.review_count} reviews from users
                    </Typography>
                  </TabPanel>
                </CardContent>
              </Card>
            </motion.div>

            {/* Creator Info */}
            <Card
              sx={{
                background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
                border: '1px solid rgba(255, 255, 255, 0.1)',
                borderRadius: 3,
              }}
            >
              <CardContent>
                <Typography variant="h6" fontWeight="bold" mb={2}>
                  Creator
                </Typography>
                <Box display="flex" alignItems="center" gap={2}>
                  <Avatar
                    src={agentData.creator?.avatar}
                    sx={{
                      width: 56,
                      height: 56,
                      background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                    }}
                  >
                    {agentData.creator?.first_name?.[0] || 'U'}
                  </Avatar>
                  <Box>
                    <Typography variant="subtitle1" fontWeight="bold">
                      {agentData.creator?.first_name} {agentData.creator?.last_name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {agentData.creator?.email}
                    </Typography>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Container>
    </Box>
  );
};

export default AgentDetail;

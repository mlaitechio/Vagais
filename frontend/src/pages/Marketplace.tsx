import React, { useState, useEffect } from 'react';
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
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Checkbox,
  FormControlLabel,
  Slider,
  IconButton,
  Badge,
  useTheme,
  InputAdornment,
  ToggleButton,
  ToggleButtonGroup,
} from '@mui/material';
import {
  Search,
  ViewModule,
  ViewList,
  Star,
  Download,
  Visibility,
  Favorite,
  FavoriteBorder,
  FilterList,
  Sort,
  SmartToy,
  AttachMoney,
  FreeBreakfast,
} from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import apiService from '../services/api';
import { Agent, SearchAgentsRequest } from '../types/api';
import { useCategories } from '../hooks/useCategories';

const Marketplace: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [selectedPricing, setSelectedPricing] = useState<string[]>([]);
  const [priceRange, setPriceRange] = useState<[number, number]>([0, 100]);
  const [sortBy, setSortBy] = useState<string>('rating');
  const [favorites, setFavorites] = useState<string[]>([]);

  // Search parameters
  const searchParams: SearchAgentsRequest = {
    query: searchQuery,
    category: selectedCategory || undefined,
    price_min: priceRange[0],
    price_max: priceRange[1],
    sort_by: sortBy as any,
    sort_order: 'desc',
    page: 1,
    limit: 20,
  };

  // Fetch agents
  const { data: agents, isLoading } = useQuery({
    queryKey: ['marketplace', searchParams],
    queryFn: () => apiService.searchAgents(searchParams),
    placeholderData: (previousData) => previousData,
    retry: false,
  });



  // Fetch categories using the custom hook
  const { data: categories, isLoading: categoriesLoading } = useCategories();

  // Ensure categories is always an array
  const categoriesList = categories && typeof categories === 'object' ? Object.keys(categories as Record<string, number>) : [];

  const toggleFavorite = (agentId: string) => {
    setFavorites(prev =>
      prev.includes(agentId)
        ? prev.filter(id => id !== agentId)
        : [...prev, agentId]
    );
  };

  const AgentCard = ({ agent }: { agent: Agent }) => (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      layout
    >
      <Card
        sx={{
          background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
          borderRadius: 3,
          height: '100%',
          position: 'relative',
          overflow: 'hidden',
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'translateY(-4px)',
            boxShadow: '0 8px 32px rgba(152, 23, 126, 0.3)',
            border: `1px solid ${theme.palette.primary.main}`,
          },
          '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            height: '3px',
            background: `linear-gradient(90deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
          },
        }}
      >
        <CardContent>
          <Box display="flex" alignItems="flex-start" mb={2}>
            <Avatar
              src={agent.icon}
              sx={{
                width: 64,
                height: 64,
                mr: 2,
                background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
              }}
            >
              <SmartToy />
            </Avatar>
            <Box flex={1}>
              <Box display="flex" alignItems="center" justifyContent="space-between">
                <Typography variant="h6" fontWeight="bold" noWrap>
                  {agent.name}
                </Typography>
                <IconButton
                  size="small"
                  onClick={() => toggleFavorite(agent.id)}
                  sx={{ color: favorites.includes(agent.id) ? 'error.main' : 'text.secondary' }}
                >
                  {favorites.includes(agent.id) ? <Favorite /> : <FavoriteBorder />}
                </IconButton>
              </Box>
              <Typography variant="body2" color="text.secondary" mb={1}>
                {agent.category}
              </Typography>
              <Box display="flex" alignItems="center" gap={1}>
                <Star sx={{ color: 'warning.main', fontSize: 16 }} />
                <Typography variant="body2" fontWeight="medium">
                  {agent.rating}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  ({agent.review_count} reviews)
                </Typography>
              </Box>
            </Box>
          </Box>

          <Typography variant="body2" color="text.secondary" mb={2} sx={{ lineHeight: 1.6 }}>
            {agent.description}
          </Typography>

          <Box display="flex" gap={1} flexWrap="wrap" mb={2}>
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
                    tags = agent.tags.includes(',') ? agent.tags.split(',').map(t => t.trim()) : [agent.tags];
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

          <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
            <Box display="flex" alignItems="center" gap={2}>
              <Box display="flex" alignItems="center">
                <Visibility sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
                <Typography variant="body2" color="text.secondary">
                  {agent.usage_count}
                </Typography>
              </Box>
              <Box display="flex" alignItems="center">
                <Download sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
                <Typography variant="body2" color="text.secondary">
                  {agent.downloads}
                </Typography>
              </Box>
            </Box>
            <Box display="flex" alignItems="center" gap={1}>
              {agent.pricing_model === 'free' ? (
                <FreeBreakfast color="success" />
              ) : agent.pricing_model === 'subscription' ? (
                <AttachMoney color="primary" />
              ) : (
                <AttachMoney color="warning" />
              )}
              <Typography variant="h6" color="primary" fontWeight="bold">
                {agent.price === 0 ? 'Free' : `$${agent.price}`}
              </Typography>
            </Box>
          </Box>

          <Box display="flex" gap={1}>
            <Button
              variant="contained"
              fullWidth
              onClick={() => navigate(`/chat/${agent.id}`)}
              sx={{
                background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                borderRadius: 2,
                '&:hover': {
                  background: `linear-gradient(135deg, ${theme.palette.primary.dark}, ${theme.palette.secondary.dark})`,
                },
              }}
            >
              Try Agent
            </Button>
            <Button
              variant="outlined"
              sx={{ borderRadius: 2, minWidth: 'auto' }}
            >
              <Download />
            </Button>
          </Box>
        </CardContent>
      </Card>
    </motion.div>
  );

  return (
    <Box sx={{ p: 3, minHeight: '100vh', background: theme.palette.background.default }}>
      {/* Header */}
      <Box textAlign="center" mb={4}>
        <Typography variant="h3" fontWeight="bold" mb={2}>
          AI Agents Directory
        </Typography>
        <Typography variant="h6" color="text.secondary">
          Your One-Stop Destination to Explore and Learn About Modern AI Agents
        </Typography>
      </Box>

      {/* Search Bar */}
      <Box mb={4}>
        <TextField
          fullWidth
          placeholder="Search for AI Agent here..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search />
              </InputAdornment>
            ),
            sx: {
              borderRadius: 3,
              background: 'rgba(255, 255, 255, 0.1)',
              border: '1px solid rgba(255, 255, 255, 0.2)',
              '& .MuiOutlinedInput-notchedOutline': {
                border: 'none',
              },
            },
          }}
        />
      </Box>

      <Grid container spacing={3}>
        {/* Left Sidebar - Filters */}
        <Grid item xs={12} md={3}>
          <Card
            sx={{
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
              borderRadius: 3,
              p: 2,
            }}
          >
            <Typography variant="h6" fontWeight="bold" mb={3}>
              Refine Search
            </Typography>

            {/* Categories */}
            <Box mb={3}>
              <Typography variant="subtitle2" fontWeight="bold" mb={2}>
                Categories
              </Typography>
              <TextField
                size="small"
                placeholder="Search categories"
                fullWidth
                sx={{ mb: 2 }}
              />
              <Box display="flex" flexDirection="column" gap={1}>
                {categoriesList.length > 0 ? (
                  categoriesList.map((category) => (
                    <FormControlLabel
                      key={category}
                      control={
                        <Checkbox
                          checked={selectedCategory === category}
                          onChange={(e) => setSelectedCategory(e.target.checked ? category : '')}
                          size="small"
                        />
                      }
                      label={`${category} (${(categories as Record<string, number>)?.[category] || 0})`}
                      sx={{ fontSize: '0.9rem' }}
                    />
                  ))
                ) : (
                  <Typography variant="body2" color="text.secondary">
                    No categories available
                  </Typography>
                )}
              </Box>
            </Box>

            {/* Pricing Models */}
            <Box mb={3}>
              <Typography variant="subtitle2" fontWeight="bold" mb={2}>
                Pricing Models
              </Typography>
              <Box display="flex" flexDirection="column" gap={1}>
                {['Free', 'Freemium', 'Paid'].map((model) => (
                  <FormControlLabel
                    key={model}
                    control={
                      <Checkbox
                        checked={selectedPricing.includes(model)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setSelectedPricing([...selectedPricing, model]);
                          } else {
                            setSelectedPricing(selectedPricing.filter(p => p !== model));
                          }
                        }}
                        size="small"
                      />
                    }
                    label={`${model} (${Math.floor(Math.random() * 15) + 1})`}
                    sx={{ fontSize: '0.9rem' }}
                  />
                ))}
              </Box>
            </Box>

            {/* Price Range */}
            <Box mb={3}>
              <Typography variant="subtitle2" fontWeight="bold" mb={2}>
                Price Range
              </Typography>
              <Slider
                value={priceRange}
                onChange={(_, value) => setPriceRange(value as [number, number])}
                valueLabelDisplay="auto"
                min={0}
                max={100}
                sx={{
                  color: theme.palette.primary.main,
                  '& .MuiSlider-thumb': {
                    background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                  },
                }}
              />
              <Box display="flex" justifyContent="space-between">
                <Typography variant="body2">${priceRange[0]}</Typography>
                <Typography variant="body2">${priceRange[1]}</Typography>
              </Box>
            </Box>

            {/* Sort Options */}
            <Box>
              <Typography variant="subtitle2" fontWeight="bold" mb={2}>
                Sort By
              </Typography>
              <FormControl fullWidth size="small">
                <Select
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value)}
                  sx={{ borderRadius: 2 }}
                >
                  <MenuItem value="rating">Rating</MenuItem>
                  <MenuItem value="downloads">Downloads</MenuItem>
                  <MenuItem value="created_at">Newest</MenuItem>
                  <MenuItem value="price">Price</MenuItem>
                </Select>
              </FormControl>
            </Box>
          </Card>
        </Grid>

        {/* Right Content - Agent Listings */}
        <Grid item xs={12} md={9}>
          {/* Results Header */}
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={3}>
            <Typography variant="h6" fontWeight="bold">
              Showing ({(agents as any)?.total || 0}) AI Agents
            </Typography>
            <Box display="flex" alignItems="center" gap={2}>
              <ToggleButtonGroup
                value={viewMode}
                exclusive
                onChange={(_, newMode) => newMode && setViewMode(newMode)}
                size="small"
              >
                <ToggleButton value="grid">
                  <ViewModule />
                </ToggleButton>
                <ToggleButton value="list">
                  <ViewList />
                </ToggleButton>
              </ToggleButtonGroup>
            </Box>
          </Box>

          {/* Agent Grid */}
          <AnimatePresence>
            <Grid container spacing={3}>
              {(agents as any)?.data?.map((agent: Agent) => (
                <Grid
                  item
                  xs={12}
                  sm={viewMode === 'grid' ? 6 : 12}
                  md={viewMode === 'grid' ? 4 : 12}
                  lg={viewMode === 'grid' ? 3 : 12}
                  key={agent.id}
                >
                  <AgentCard agent={agent} />
                </Grid>
              ))}
            </Grid>
          </AnimatePresence>

          {/* No Results */}
          {(agents as any)?.data?.length === 0 && !isLoading && (
            <Box textAlign="center" py={8}>
              <SmartToy sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary" mb={1}>
                No agents found
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Try adjusting your search criteria or browse all agents
              </Typography>
            </Box>
          )}

          {/* Loading State */}
          {isLoading && (
            <Box textAlign="center" py={8}>
              <Typography variant="h6" color="text.secondary">
                Loading agents...
              </Typography>
            </Box>
          )}
        </Grid>
      </Grid>
    </Box>
  );
};

export default Marketplace;

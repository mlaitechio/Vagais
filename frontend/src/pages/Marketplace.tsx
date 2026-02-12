import React, { useState } from "react";
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
  Select,
  MenuItem,

  
 
  IconButton,
  useTheme,
  InputAdornment,
  ToggleButton,
  ToggleButtonGroup,
  Stack,
} from "@mui/material";
import {
  Search,
  ViewModule,
  ViewList,
  Star,
  Download,
  Visibility,
  Favorite,
  FavoriteBorder,
  SmartToy,
  AttachMoney,
  FreeBreakfast,
} from "@mui/icons-material";
import { motion, AnimatePresence } from "framer-motion";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import apiService from "../services/api";
import { Agent, SearchAgentsRequest } from "../types/api";
import { useCategories } from "../hooks/useCategories";

const Marketplace: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCategory, setSelectedCategory] = useState<string>("");
  const [selectedPricing, setSelectedPricing] = useState<string[]>([]);
  // @ts-ignore
  const [priceRange,setPriceRange] = useState<[number, number]>([0,100 ]);
  const [sortBy, setSortBy] = useState<string>("rating");
  const [favorites, setFavorites] = useState<string[]>([]);

  // Search parameters
  const searchParams: SearchAgentsRequest = {
    query: searchQuery,
    category: selectedCategory || undefined,
    price_min: priceRange[0],
    price_max: priceRange[1],
    sort_by: sortBy as any,
    sort_order: "desc",
    page: 1,
    limit: 20,
  };

  // Fetch agents
  const { data: agents, isLoading } = useQuery({
    queryKey: ["marketplace", searchParams],
    queryFn: () => apiService.searchAgents(searchParams),
    placeholderData: (previousData) => previousData,
    retry: false,
  });

  // Fetch categories using the custom hook
  const { data: categories } = useCategories();

  // Ensure categories is always an array
  const categoriesList =
    categories && typeof categories === "object"
      ? Object.keys(categories as Record<string, number>)
      : [];

  const toggleFavorite = (agentId: string) => {
    setFavorites((prev) =>
      prev.includes(agentId)
        ? prev.filter((id) => id !== agentId)
        : [...prev, agentId],
    );
  };

  const AgentCard = ({ agent }: { agent: Agent }) => (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      layout
      style={{ display: "flex", flex: 1 }}
    >
      <Card
        sx={{
          background:
            "linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%)",
          border: "1px solid rgba(255, 255, 255, 0.1)",
          borderRadius: 3,
          height: "100%",
          position: "relative",
          overflow: "hidden",
          cursor: "pointer",
          transition: "all 0.3s ease",
          "&:hover": {
            transform: "translateY(-4px)",
            boxShadow: "0 8px 32px rgba(152, 23, 126, 0.3)",
            border: `1px solid ${theme.palette.primary.main}`,
          },
          "&::before": {
            content: '""',
            position: "absolute",
            top: 0,
            left: 0,
            right: 0,
            height: "3px",
            background: `linear-gradient(90deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
          },
        }}
        onClick={() => navigate(`/agents/${agent.id}`)}
      >
        <CardContent style={{ height: "100%" }}>
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              justifyContent: "space-between",
              height: "100%",
            }}
          >
            <div>
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
                  <Box
                    display="flex"
                    alignItems="center"
                    justifyContent="space-between"
                  >
                    <Typography variant="h6" fontWeight="bold">
                      {agent.name}
                    </Typography>
                    <IconButton
                      size="small"
                      onClick={(e) => {
                        e.stopPropagation();
                        toggleFavorite(agent.id);
                      }}
                      sx={{
                        color: favorites.includes(agent.id)
                          ? "error.main"
                          : "text.secondary",
                      }}
                    >
                      {favorites.includes(agent.id) ? (
                        <Favorite />
                      ) : (
                        <FavoriteBorder />
                      )}
                    </IconButton>
                  </Box>
                  <Typography variant="body2" color="text.secondary" mb={1}>
                    {agent.category}
                  </Typography>
                  <Box display="flex" alignItems="center" gap={1}>
                    <Star sx={{ color: "warning.main", fontSize: 16 }} />
                    <Typography variant="body2" fontWeight="medium">
                      {agent.rating}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      ({agent.review_count} reviews)
                    </Typography>
                  </Box>
                </Box>
              </Box>

              <Typography
                variant="body2"
                color="text.secondary"
                mb={2}
                sx={{ lineHeight: 1.6 }}
              >
                {agent.description}
              </Typography>

              <Box display="flex" gap={1} flexWrap="wrap" mb={2}>
                {(() => {
                  let tags: string[] = [];
                  if (agent.tags) {
                    if (Array.isArray(agent.tags)) {
                      tags = agent.tags;
                    } else if (typeof agent.tags === "string") {
                      try {
                        // Handle base64 encoded JSON string
                        const decoded = atob(agent.tags);
                        const parsed = JSON.parse(decoded);
                        tags = Array.isArray(parsed) ? parsed : [];
                      } catch {
                        // If parsing fails, try to split by comma or treat as single tag
                        const tagStr = agent.tags as any as string;
                        tags =
                          tagStr && tagStr.includes(",")
                            ? tagStr.split(",").map((t: string) => t.trim())
                            : [tagStr];
                      }
                    }
                  }
                  return tags
                    .slice(0, 3)
                    .map((tag, index) => (
                      <Chip
                        key={index}
                        label={tag}
                        size="small"
                        variant="outlined"
                        sx={{ fontSize: "0.7rem" }}
                      />
                    ));
                })()}
              </Box>
            </div>

            <div>
              <Box
                display="flex"
                alignItems="center"
                justifyContent="space-between"
                mb={2}
              >
                <Box display="flex" alignItems="center" gap={2}>
                  <Box display="flex" alignItems="center">
                    <Visibility
                      sx={{ fontSize: 16, mr: 0.5, color: "text.secondary" }}
                    />
                    <Typography variant="body2" color="text.secondary">
                      {agent.usage_count}
                    </Typography>
                  </Box>
                  <Box display="flex" alignItems="center">
                    <Download
                      sx={{ fontSize: 16, mr: 0.5, color: "text.secondary" }}
                    />
                    <Typography variant="body2" color="text.secondary">
                      {agent.downloads}
                    </Typography>
                  </Box>
                </Box>
                <Box display="flex" alignItems="center" gap={1}>
                  {agent.pricing_model === "free" ? (
                    <FreeBreakfast color="success" />
                  ) : agent.pricing_model === "subscription" ? (
                    <AttachMoney color="primary" />
                  ) : (
                    <AttachMoney color="warning" />
                  )}
                  <Typography variant="h6" color="primary" fontWeight="bold">
                    {agent.price === 0 ? "Free" : `$${agent.price}`}
                  </Typography>
                </Box>
              </Box>

              <Box display="flex" gap={1}>
                <Button
                  variant="contained"
                  fullWidth
                  onClick={(e) => {
                    e.stopPropagation();
                    navigate(`/chat/${agent.id}`);
                  }}
                  sx={{
                    background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                    borderRadius: 2,
                    "&:hover": {
                      background: `linear-gradient(135deg, ${theme.palette.primary.dark}, ${theme.palette.secondary.dark})`,
                    },
                  }}
                >
                  Try Agent
                </Button>
                <Button
                  variant="outlined"
                  onClick={(e) => e.stopPropagation()}
                  sx={{ borderRadius: 2, minWidth: "auto" }}
                >
                  <Download />
                </Button>
              </Box>
            </div>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );

  return (
    <Box
      sx={{
        p: 3,
        minHeight: "100vh",
        background: theme.palette.background.default,
      }}
    >
      <Box sx={{ maxWidth: 2200, margin: "50px auto 0" }}>
        {/* Header */}
        <Box
          textAlign="center"
          mb={4}
          mx="auto"
          sx={{
            maxWidth: { xs: "100%", md: 1200 },
            px: { xs: 2, md: 4 },
          }}
        >
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <Stack spacing={2} alignItems="center">
              <Typography
                variant="h5"
                sx={{
                  color: "text.disabled",
                  maxWidth: 600,
                  mx: "auto",
                }}
              >
                One-Click Agent Deployment Across Cloud Platforms
              </Typography>

              {/* CENTER: Title */}
              <Box textAlign="center">
                <Typography
                  variant="h1"
                  sx={{
                    fontWeight: 800,
                    fontSize: { xs: "2.5rem", md: "4rem" },
                    lineHeight: 1.1,
                    mb: 2,
                  }}
                >
                  AI Agent{" "}
                  <Box
                    component="span"
                    sx={{
                      background: `linear-gradient(135deg, ${theme.palette.primary.main} 0%, ${theme.palette.secondary.main} 100%)`,
                      backgroundClip: "text",
                      WebkitBackgroundClip: "text",
                      WebkitTextFillColor: "transparent",
                    }}
                  >
                    Marketplace
                  </Box>
                </Typography>

                {/* BOTTOM: Description */}
                <Typography
                  variant="h6"
                  sx={{
                    color: "text.secondary",
                    fontWeight: 400,
                    maxWidth: 700,
                    mx: "auto",
                    mb: 1,
                  }}
                >
                  Discover, deploy, and manage the most powerful AI agents
                </Typography>

                <Typography
                  variant="body2"
                  sx={{
                    color: "text.disabled",
                    maxWidth: 600,
                    mx: "auto",
                  }}
                >
                  Filter by capability, price, and ratings
                </Typography>
              </Box>
            </Stack>
          </motion.div>
        </Box>

        {/* Search + Filters */}
        <Grid className="adsearch" container spacing={2} alignItems="flex-end" sx={{ mb: 4 }}>
          {/* Sort By */}
          <Grid item xs={12} md={2}>
            <FormControl fullWidth size="small">
              <Typography variant="caption" fontWeight="bold" mb={0.5}>
                Sort By
              </Typography>
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
          </Grid>

          {/* Search */}
          <Grid item xs={12} md={6}>
            <TextField
              fullWidth
              size="small"
              placeholder=""
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
                  background: "rgba(255,255,255,0.08)",
                  border: "1px solid rgba(255,255,255,0.15)",
                  "& .MuiOutlinedInput-notchedOutline": { border: "none" },
                },
              }}
            />
          </Grid>

          {/* Category */}
          <Grid item xs={12} md={2}>
            <FormControl fullWidth size="small">
              <Typography variant="caption" fontWeight="bold" mb={0.5}>
                Category
              </Typography>
              <Select
                value={selectedCategory || ""}
                displayEmpty
                onChange={(e) => setSelectedCategory(e.target.value)}
                sx={{ borderRadius: 2 }}
              >
                <MenuItem value="">
                  <em>All</em>
                </MenuItem>
                {categoriesList.map((category) => (
                  <MenuItem key={category} value={category}>
                    {category} (
                    {(categories as Record<string, number>)?.[category] || 0})
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>

          {/* Pricing */}
          <Grid item xs={12} md={2}>
            <FormControl fullWidth size="small">
              <Typography variant="caption" fontWeight="bold" mb={0.5}>
                Pricing Model
              </Typography>
              <Select
                value={selectedPricing || ""}
                displayEmpty
                onChange={(e) =>
                  setSelectedPricing(
                    typeof e.target.value === "string"
                      ? e.target.value.split(",")
                      : e.target.value,
                  )
                }
                sx={{ borderRadius: 2 }}
              >
                <MenuItem value="">
                  <em>All</em>
                </MenuItem>
                {["Free", "Freemium", "Paid"].map((category) => (
                  <MenuItem key={category} value={category}>
                    {category}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
        </Grid>

        {/* Contents  */}
        <Grid className="ad_grid_box" container spacing={3}>
          <Grid item xs={12} md={12}>
            {/* Results Header */}
            <Box
              display="flex"
              alignItems="center"
              justifyContent="space-between"
              mb={3}
            >
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
                    sm={viewMode === "grid" ? 6 : 12}
                    md={viewMode === "grid" ? 3 : 12}
                    lg={viewMode === "grid" ? 3 : 12}
                    key={agent.id}
                    sx={{ display: "flex" }}
                  >
                    <AgentCard agent={agent} />
                  </Grid>
                ))}
              </Grid>
            </AnimatePresence>

            {/* No Results */}
            {(agents as any)?.data?.length === 0 && !isLoading && (
              <Box textAlign="center" py={8}>
                <SmartToy
                  sx={{ fontSize: 64, color: "text.secondary", mb: 2 }}
                />
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
    </Box>
  );
};

export default Marketplace;

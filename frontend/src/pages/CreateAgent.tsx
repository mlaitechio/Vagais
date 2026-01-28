import React, { useState } from 'react';
import {
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  Box,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  IconButton,
  Alert,
} from '@mui/material';
import { ArrowBack, Save, Add, Delete } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { apiService } from '../services/api';

const CreateAgent: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const [formData, setFormData] = useState({
    name: '',
    description: '',
    category: '',
    pricing_model: 'free' as 'free' | 'one_time' | 'subscription',
    price: 0,
    currency: 'USD',
    tags: [] as string[],
    is_public: true,
    llm_provider: 'openai',
    llm_model: 'gpt-4',
    embedding_provider: 'openai',
    embedding_model: 'text-embedding-ada-002',
    config: {
      max_tokens: 1000,
      temperature: 0.7,
      instructions: '',
    },
  });

  const [newTag, setNewTag] = useState('');

  const handleInputChange = (field: string, value: any) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const handleAddTag = () => {
    if (newTag.trim() && !formData.tags.includes(newTag.trim())) {
      setFormData(prev => ({
        ...prev,
        tags: [...prev.tags, newTag.trim()]
      }));
      setNewTag('');
    }
  };

  const handleRemoveTag = (tagToRemove: string) => {
    setFormData(prev => ({
      ...prev,
      tags: prev.tags.filter(tag => tag !== tagToRemove)
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await apiService.createAgent(formData);
      
      setSuccess(true);
      setTimeout(() => {
        navigate('/agents');
      }, 2000);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create agent');
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          <Alert severity="success" sx={{ mb: 3 }}>
            Agent created successfully! Redirecting to agents page...
          </Alert>
        </motion.div>
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <Paper
          elevation={3}
          sx={{
            p: 4,
            background: 'linear-gradient(135deg, rgba(26, 26, 26, 0.95) 0%, rgba(10, 10, 10, 0.95) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
          }}
        >
          {/* Header */}
          <Box display="flex" alignItems="center" gap={2} mb={4}>
            <IconButton onClick={() => navigate('/agents')} sx={{ color: 'text.primary' }}>
              <ArrowBack />
            </IconButton>
            <Typography variant="h4" fontWeight="bold" sx={{ 
              background: 'linear-gradient(135deg, #98177e, #ff6b6b)',
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
            }}>
              Create New Agent
            </Typography>
          </Box>

          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          <form onSubmit={handleSubmit}>
            <Grid container spacing={3}>
              {/* Basic Information */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom>
                  Basic Information
                </Typography>
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Agent Name"
                  value={formData.name}
                  onChange={(e) => handleInputChange('name', e.target.value)}
                  required
                  sx={{ mb: 2 }}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>Category</InputLabel>
                  <Select
                    value={formData.category}
                    onChange={(e) => handleInputChange('category', e.target.value)}
                    label="Category"
                  >
                    <MenuItem value="productivity">Productivity</MenuItem>
                    <MenuItem value="creative">Creative</MenuItem>
                    <MenuItem value="analytics">Analytics</MenuItem>
                    <MenuItem value="customer-service">Customer Service</MenuItem>
                    <MenuItem value="education">Education</MenuItem>
                    <MenuItem value="entertainment">Entertainment</MenuItem>
                    <MenuItem value="healthcare">Healthcare</MenuItem>
                    <MenuItem value="finance">Finance</MenuItem>
                    <MenuItem value="marketing">Marketing</MenuItem>
                    <MenuItem value="development">Development</MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Description"
                  value={formData.description}
                  onChange={(e) => handleInputChange('description', e.target.value)}
                  multiline
                  rows={3}
                  required
                  sx={{ mb: 2 }}
                />
              </Grid>

              {/* Pricing */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom>
                  Pricing
                </Typography>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>Pricing Model</InputLabel>
                  <Select
                    value={formData.pricing_model}
                    onChange={(e) => handleInputChange('pricing_model', e.target.value)}
                    label="Pricing Model"
                  >
                    <MenuItem value="free">Free</MenuItem>
                    <MenuItem value="one_time">One-time Payment</MenuItem>
                    <MenuItem value="subscription">Subscription</MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Price ($)"
                  type="number"
                  value={formData.price}
                  onChange={(e) => handleInputChange('price', parseFloat(e.target.value) || 0)}
                  disabled={formData.pricing_model === 'free'}
                  sx={{ mb: 2 }}
                />
              </Grid>

              {/* AI Configuration */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom>
                  AI Configuration
                </Typography>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>LLM Provider</InputLabel>
                  <Select
                    value={formData.llm_provider}
                    onChange={(e) => handleInputChange('llm_provider', e.target.value)}
                    label="LLM Provider"
                  >
                    <MenuItem value="openai">OpenAI</MenuItem>
                    <MenuItem value="anthropic">Anthropic</MenuItem>
                    <MenuItem value="google">Google</MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>LLM Model</InputLabel>
                  <Select
                    value={formData.llm_model}
                    onChange={(e) => handleInputChange('llm_model', e.target.value)}
                    label="LLM Model"
                  >
                    <MenuItem value="gpt-4">GPT-4</MenuItem>
                    <MenuItem value="gpt-3.5-turbo">GPT-3.5 Turbo</MenuItem>
                    <MenuItem value="claude-3">Claude 3</MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>Embedding Provider</InputLabel>
                  <Select
                    value={formData.embedding_provider}
                    onChange={(e) => handleInputChange('embedding_provider', e.target.value)}
                    label="Embedding Provider"
                  >
                    <MenuItem value="openai">OpenAI</MenuItem>
                    <MenuItem value="cohere">Cohere</MenuItem>
                    <MenuItem value="huggingface">Hugging Face</MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth sx={{ mb: 2 }}>
                  <InputLabel>Embedding Model</InputLabel>
                  <Select
                    value={formData.embedding_model}
                    onChange={(e) => handleInputChange('embedding_model', e.target.value)}
                    label="Embedding Model"
                  >
                    <MenuItem value="text-embedding-ada-002">text-embedding-ada-002</MenuItem>
                    <MenuItem value="text-embedding-3-small">text-embedding-3-small</MenuItem>
                    <MenuItem value="text-embedding-3-large">text-embedding-3-large</MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12} md={4}>
                <TextField
                  fullWidth
                  label="Max Tokens"
                  type="number"
                  value={formData.config.max_tokens}
                  onChange={(e) => handleInputChange('config', { ...formData.config, max_tokens: parseInt(e.target.value) || 1000 })}
                  sx={{ mb: 2 }}
                />
              </Grid>

              <Grid item xs={12} md={4}>
                <TextField
                  fullWidth
                  label="Temperature"
                  type="number"
                  value={formData.config.temperature}
                  onChange={(e) => handleInputChange('config', { ...formData.config, temperature: parseFloat(e.target.value) || 0.7 })}
                  inputProps={{ step: 0.1, min: 0, max: 2 }}
                  sx={{ mb: 2 }}
                />
              </Grid>

              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Instructions"
                  value={formData.config.instructions}
                  onChange={(e) => handleInputChange('config', { ...formData.config, instructions: e.target.value })}
                  multiline
                  rows={4}
                  placeholder="Describe how the agent should behave, what it should do, and any specific guidelines..."
                  sx={{ mb: 2 }}
                />
              </Grid>

              {/* Tags */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom>
                  Tags
                </Typography>
                <Box display="flex" gap={1} mb={2} flexWrap="wrap">
                  {formData.tags.map((tag) => (
                    <Chip
                      key={tag}
                      label={tag}
                      onDelete={() => handleRemoveTag(tag)}
                      deleteIcon={<Delete />}
                      color="primary"
                      variant="outlined"
                    />
                  ))}
                </Box>
                <Box display="flex" gap={1}>
                  <TextField
                    size="small"
                    placeholder="Add tag"
                    value={newTag}
                    onChange={(e) => setNewTag(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), handleAddTag())}
                  />
                  <Button
                    variant="outlined"
                    startIcon={<Add />}
                    onClick={handleAddTag}
                    disabled={!newTag.trim()}
                  >
                    Add
                  </Button>
                </Box>
              </Grid>

              {/* Submit Button */}
              <Grid item xs={12}>
                <Box display="flex" gap={2} justifyContent="flex-end">
                  <Button
                    variant="outlined"
                    onClick={() => navigate('/agents')}
                    disabled={loading}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    variant="contained"
                    startIcon={<Save />}
                    disabled={loading}
                    sx={{
                      background: 'linear-gradient(135deg, #98177e, #ff6b6b)',
                      '&:hover': {
                        background: 'linear-gradient(135deg, #7a125f, #e55a5a)',
                      },
                    }}
                  >
                    {loading ? 'Creating...' : 'Create Agent'}
                  </Button>
                </Box>
              </Grid>
            </Grid>
          </form>
        </Paper>
      </motion.div>
    </Container>
  );
};

export default CreateAgent;

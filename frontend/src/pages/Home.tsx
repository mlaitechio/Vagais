import React from 'react';
import {
  Box,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  Container,
} from '@mui/material';
import {
  Security,
  Speed,
  Psychology,
  Business,
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';

const Home: React.FC = () => {
  const navigate = useNavigate();

  const features = [
    {
      icon: <Psychology sx={{ fontSize: 40 }} />,
      title: 'Intelligent AI Agents',
      description: 'Advanced AI agents that understand context and learn from interactions',
      color: '#98177E',
    },
    {
      icon: <Speed sx={{ fontSize: 40 }} />,
      title: 'Lightning Fast',
      description: 'Execute complex tasks in seconds, not hours',
      color: '#00D4FF',
    },
    {
      icon: <Security sx={{ fontSize: 40 }} />,
      title: 'Enterprise Security',
      description: 'Bank-grade security with end-to-end encryption',
      color: '#00FF88',
    },
    {
      icon: <Business sx={{ fontSize: 40 }} />,
      title: 'Business Ready',
      description: 'Seamlessly integrate with your existing workflows',
      color: '#FFB800',
    },
  ];

  const aiBenefits = [
    {
      title: 'Automate Repetitive Tasks',
      description: 'Free your team from mundane work and focus on innovation',
      icon: '🤖',
    },
    {
      title: '24/7 Availability',
      description: 'AI agents work around the clock, never sleep, never tire',
      icon: '⏰',
    },
    {
      title: 'Scalable Solutions',
      description: 'Handle thousands of requests simultaneously',
      icon: '📈',
    },
    {
      title: 'Data-Driven Insights',
      description: 'Extract meaningful patterns and predictions from your data',
      icon: '📊',
    },
  ];

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.3,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.6,
      },
    },
  };

  return (
    <Box sx={{ minHeight: '100vh', position: 'relative' }}>
      {/* Hero Section */}
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          position: 'relative',
          overflow: 'hidden',
        }}
      >
        {/* Animated background */}
        <Box
          sx={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            zIndex: -1,
          }}
        >
          <motion.div
            animate={{
              scale: [1, 1.1, 1],
              rotate: [0, 5, 0],
            }}
            transition={{
              duration: 20,
              repeat: Infinity,
              ease: "easeInOut",
            }}
            style={{
              position: 'absolute',
              top: '10%',
              left: '5%',
              width: '400px',
              height: '400px',
              background: 'radial-gradient(circle, rgba(152, 23, 126, 0.1) 0%, transparent 70%)',
              borderRadius: '50%',
            }}
          />
          <motion.div
            animate={{
              scale: [1.1, 1, 1.1],
              rotate: [0, -5, 0],
            }}
            transition={{
              duration: 25,
              repeat: Infinity,
              ease: "easeInOut",
            }}
            style={{
              position: 'absolute',
              top: '60%',
              right: '10%',
              width: '300px',
              height: '300px',
              background: 'radial-gradient(circle, rgba(0, 212, 255, 0.1) 0%, transparent 70%)',
              borderRadius: '50%',
            }}
          />
        </Box>

        <Container maxWidth="lg">
          <Grid container spacing={4} alignItems="center">
            <Grid item xs={12} md={6}>
              <motion.div
                initial={{ opacity: 0, x: -50 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.8 }}
              >
                <Typography
                  variant="h1"
                  sx={{
                    fontSize: { xs: '2.5rem', md: '4rem', lg: '5rem' },
                    fontWeight: 700,
                    mb: 3,
                    lineHeight: 1.1,
                  }}
                >
                  The Future of
                  <br />
                  <Box
                    component="span"
                    sx={{
                      background: 'linear-gradient(135deg, #98177E 0%, #00D4FF 100%)',
                      backgroundClip: 'text',
                      WebkitBackgroundClip: 'text',
                      WebkitTextFillColor: 'transparent',
                    }}
                  >
                    AI Agents
                  </Box>
                </Typography>
                <Typography
                  variant="h5"
                  sx={{
                    color: 'text.secondary',
                    mb: 4,
                    lineHeight: 1.6,
                    fontSize: { xs: '1.1rem', md: '1.3rem' },
                  }}
                >
                  Discover, deploy, and manage intelligent AI agents that transform your business.
                  Built for 2026 and beyond.
                </Typography>
                <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                  <motion.div
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                  >
                    <Button
                      variant="contained"
                      size="large"
                      onClick={() => navigate('/marketplace')}
                      sx={{
                        px: 4,
                        py: 1.5,
                        fontSize: '1.1rem',
                        fontWeight: 600,
                      }}
                    >
                      Explore Agents
                    </Button>
                  </motion.div>
                  <motion.div
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                  >
                    <Button
                      variant="outlined"
                      size="large"
                      onClick={() => navigate('/register')}
                      sx={{
                        px: 4,
                        py: 1.5,
                        fontSize: '1.1rem',
                        fontWeight: 600,
                        borderColor: '#98177E',
                        color: '#98177E',
                        '&:hover': {
                          borderColor: '#B23A9A',
                          backgroundColor: 'rgba(152, 23, 126, 0.1)',
                        },
                      }}
                    >
                      Get Started
                    </Button>
                  </motion.div>
                </Box>
              </motion.div>
            </Grid>
            <Grid item xs={12} md={6}>
              <motion.div
                initial={{ opacity: 0, x: 50 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.8, delay: 0.2 }}
              >
                <Box
                  sx={{
                    position: 'relative',
                    height: '500px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                >
                  {/* AI Agent Visualization */}
                  <motion.div
                    animate={{
                      y: [0, -20, 0],
                      rotate: [0, 5, 0],
                    }}
                    transition={{
                      duration: 6,
                      repeat: Infinity,
                      ease: "easeInOut",
                    }}
                    style={{
                      position: 'relative',
                      width: '300px',
                      height: '300px',
                    }}
                  >
                    <Box
                      sx={{
                        width: '100%',
                        height: '100%',
                        borderRadius: '50%',
                        background: 'linear-gradient(135deg, #98177E 0%, #00D4FF 50%, #00FF88 100%)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        boxShadow: '0 20px 60px rgba(152, 23, 126, 0.3)',
                        position: 'relative',
                      }}
                    >
                      <Typography
                        variant="h2"
                        sx={{
                          color: 'white',
                          fontWeight: 700,
                          textAlign: 'center',
                        }}
                      >
                        AI
                      </Typography>
                    </Box>
                    
                    {/* Orbiting elements */}
                    {[...Array(6)].map((_, index) => (
                      <motion.div
                        key={index}
                        animate={{
                          rotate: [0, 360],
                        }}
                        transition={{
                          duration: 8 + index * 2,
                          repeat: Infinity,
                          ease: "linear",
                        }}
                        style={{
                          position: 'absolute',
                          top: '50%',
                          left: '50%',
                          width: '20px',
                          height: '20px',
                          marginLeft: '-10px',
                          marginTop: '-10px',
                          transform: `rotate(${index * 60}deg) translateY(-150px)`,
                        }}
                      >
                        <Box
                          sx={{
                            width: '20px',
                            height: '20px',
                            borderRadius: '50%',
                            background: '#00D4FF',
                            boxShadow: '0 0 20px rgba(0, 212, 255, 0.8)',
                          }}
                        />
                      </motion.div>
                    ))}
                  </motion.div>
                </Box>
              </motion.div>
            </Grid>
          </Grid>
        </Container>
      </Box>

      {/* Features Section */}
      <Box sx={{ py: 8, position: 'relative' }}>
        <Container maxWidth="lg">
          <motion.div
            variants={containerVariants}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
          >
            <motion.div variants={itemVariants}>
              <Typography
                variant="h2"
                align="center"
                sx={{ mb: 6, fontWeight: 600 }}
              >
                Why Choose vagais.ai?
              </Typography>
            </motion.div>

            <Grid container spacing={4}>
              {features.map((feature, index) => (
                <Grid item xs={12} sm={6} md={3} key={index}>
                  <motion.div variants={itemVariants}>
                    <Card
                      sx={{
                        height: '100%',
                        textAlign: 'center',
                        transition: 'all 0.3s ease',
                        '&:hover': {
                          transform: 'translateY(-10px)',
                          boxShadow: '0 20px 40px rgba(0, 0, 0, 0.3)',
                        },
                      }}
                    >
                      <CardContent sx={{ p: 4 }}>
                        <Box
                          sx={{
                            color: feature.color,
                            mb: 2,
                            display: 'flex',
                            justifyContent: 'center',
                          }}
                        >
                          {feature.icon}
                        </Box>
                        <Typography
                          variant="h5"
                          sx={{ mb: 2, fontWeight: 600 }}
                        >
                          {feature.title}
                        </Typography>
                        <Typography
                          variant="body1"
                          color="text.secondary"
                          sx={{ lineHeight: 1.6 }}
                        >
                          {feature.description}
                        </Typography>
                      </CardContent>
                    </Card>
                  </motion.div>
                </Grid>
              ))}
            </Grid>
          </motion.div>
        </Container>
      </Box>

      {/* AI Benefits Section */}
      <Box sx={{ py: 8, backgroundColor: 'rgba(152, 23, 126, 0.05)' }}>
        <Container maxWidth="lg">
          <motion.div
            variants={containerVariants}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
          >
            <motion.div variants={itemVariants}>
              <Typography
                variant="h2"
                align="center"
                sx={{ mb: 6, fontWeight: 600 }}
              >
                How AI Transforms Business
              </Typography>
            </motion.div>

            <Grid container spacing={4}>
              {aiBenefits.map((benefit, index) => (
                <Grid item xs={12} sm={6} md={3} key={index}>
                  <motion.div variants={itemVariants}>
                    <Card
                      sx={{
                        height: '100%',
                        textAlign: 'center',
                        transition: 'all 0.3s ease',
                        '&:hover': {
                          transform: 'translateY(-10px)',
                          boxShadow: '0 20px 40px rgba(0, 0, 0, 0.3)',
                        },
                      }}
                    >
                      <CardContent sx={{ p: 4 }}>
                        <Typography
                          variant="h1"
                          sx={{ mb: 2, fontSize: '3rem' }}
                        >
                          {benefit.icon}
                        </Typography>
                        <Typography
                          variant="h6"
                          sx={{ mb: 2, fontWeight: 600 }}
                        >
                          {benefit.title}
                        </Typography>
                        <Typography
                          variant="body1"
                          color="text.secondary"
                          sx={{ lineHeight: 1.6 }}
                        >
                          {benefit.description}
                        </Typography>
                      </CardContent>
                    </Card>
                  </motion.div>
                </Grid>
              ))}
            </Grid>
          </motion.div>
        </Container>
      </Box>

      {/* CTA Section */}
      <Box sx={{ py: 8 }}>
        <Container maxWidth="md">
          <motion.div
            initial={{ opacity: 0, y: 50 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
          >
            <Card
              sx={{
                textAlign: 'center',
                p: 6,
                background: 'linear-gradient(135deg, rgba(152, 23, 126, 0.1) 0%, rgba(0, 212, 255, 0.1) 100%)',
                border: '1px solid rgba(152, 23, 126, 0.3)',
              }}
            >
              <Typography
                variant="h3"
                sx={{ mb: 3, fontWeight: 600 }}
              >
                Ready to Transform Your Business?
              </Typography>
              <Typography
                variant="h6"
                color="text.secondary"
                sx={{ mb: 4, lineHeight: 1.6 }}
              >
                Join thousands of businesses already using AI agents to automate, 
                optimize, and innovate their operations.
              </Typography>
              <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center', flexWrap: 'wrap' }}>
                <motion.div
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                >
                  <Button
                    variant="contained"
                    size="large"
                    onClick={() => navigate('/register')}
                    sx={{
                      px: 4,
                      py: 1.5,
                      fontSize: '1.1rem',
                      fontWeight: 600,
                    }}
                  >
                    Start Free Trial
                  </Button>
                </motion.div>
                <motion.div
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                >
                  <Button
                    variant="outlined"
                    size="large"
                    onClick={() => navigate('/marketplace')}
                    sx={{
                      px: 4,
                      py: 1.5,
                      fontSize: '1.1rem',
                      fontWeight: 600,
                      borderColor: '#98177E',
                      color: '#98177E',
                      '&:hover': {
                        borderColor: '#B23A9A',
                        backgroundColor: 'rgba(152, 23, 126, 0.1)',
                      },
                    }}
                  >
                    Browse Marketplace
                  </Button>
                </motion.div>
              </Box>
            </Card>
          </motion.div>
        </Container>
      </Box>
    </Box>
  );
};

export default Home; 
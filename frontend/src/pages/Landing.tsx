import React from 'react';
import { Box, Container, Typography, Button, Grid, Card, CardContent, CardMedia } from '@mui/material';
import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';

const Landing: React.FC = () => {
  const navigate = useNavigate();

  const features = [
    {
      title: 'AI Agent Marketplace',
      description: 'Discover and deploy powerful AI agents for your business needs',
      icon: '🤖',
    },
    {
      title: 'Easy Integration',
      description: 'Seamlessly integrate AI agents into your existing workflows',
      icon: '🔗',
    },
    {
      title: 'Advanced Analytics',
      description: 'Track performance and optimize your AI agents with detailed analytics',
      icon: '📊',
    },
    {
      title: 'Secure & Reliable',
      description: 'Enterprise-grade security and reliability for your AI operations',
      icon: '🔒',
    },
  ];

  return (
    <Box sx={{ minHeight: '100vh', background: 'linear-gradient(135deg, #0A0A0A 0%, #1A1A1A 100%)' }}>
      <Container maxWidth="lg">
        {/* Hero Section */}
        <Box sx={{ py: 8, textAlign: 'center' }}>
          <motion.div
            initial={{ opacity: 0, y: 50 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
          >
                         <Typography variant="h1" sx={{ mb: 3 }}>
               Welcome to vagais.ai
             </Typography>
             <Typography variant="h4" sx={{ mb: 4, color: 'text.secondary' }}>
               The Future of AI Agent Development and Deployment
             </Typography>
                         <Box sx={{ mt: 4 }}>
               <Button
                 variant="contained"
                 size="large"
                 onClick={() => navigate('/demo-marketplace')}
                 sx={{ mr: 2, mb: 2 }}
               >
                 Explore Marketplace
               </Button>
               <Button
                 variant="outlined"
                 size="large"
                 onClick={() => navigate('/demo-dashboard')}
                 sx={{ mr: 2, mb: 2 }}
               >
                 View Dashboard
               </Button>
               <Button
                 variant="outlined"
                 size="large"
                 onClick={() => navigate('/login')}
                 sx={{ mb: 2 }}
               >
                 Get Started
               </Button>
             </Box>
          </motion.div>
        </Box>

        {/* Features Section */}
        <Box sx={{ py: 8 }}>
                     <Typography variant="h2" sx={{ textAlign: 'center', mb: 6 }}>
             Why Choose vagais.ai?
           </Typography>
          <Grid container spacing={4}>
            {features.map((feature, index) => (
              <Grid item xs={12} sm={6} md={3} key={index}>
                <motion.div
                  initial={{ opacity: 0, y: 50 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.8, delay: index * 0.2 }}
                >
                  <Card
                    sx={{
                      height: '100%',
                      background: 'linear-gradient(135deg, #1A1A1A 0%, #2A2A2A 100%)',
                      border: '1px solid rgba(255, 255, 255, 0.1)',
                    }}
                  >
                    <CardContent sx={{ textAlign: 'center', py: 4 }}>
                      <Typography variant="h1" sx={{ mb: 2, fontSize: '3rem' }}>
                        {feature.icon}
                      </Typography>
                      <Typography variant="h5" sx={{ mb: 2, fontWeight: 600 }}>
                        {feature.title}
                      </Typography>
                      <Typography variant="body1" color="text.secondary">
                        {feature.description}
                      </Typography>
                    </CardContent>
                  </Card>
                </motion.div>
              </Grid>
            ))}
          </Grid>
        </Box>

        {/* CTA Section */}
        <Box sx={{ py: 8, textAlign: 'center' }}>
          <motion.div
            initial={{ opacity: 0, y: 50 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.5 }}
          >
            <Typography variant="h3" sx={{ mb: 3 }}>
              Ready to Transform Your Business?
            </Typography>
                         <Typography variant="h6" sx={{ mb: 4, color: 'text.secondary' }}>
               Join thousands of businesses already using vagais.ai
             </Typography>
                         <Button
               variant="contained"
               size="large"
               onClick={() => navigate('/demo-dashboard')}
               sx={{ mr: 2 }}
             >
               Try Dashboard
             </Button>
             <Button
               variant="outlined"
               size="large"
               onClick={() => navigate('/register')}
               sx={{ mr: 2 }}
             >
               Sign Up Now
             </Button>
          </motion.div>
        </Box>
      </Container>
    </Box>
  );
};

export default Landing;

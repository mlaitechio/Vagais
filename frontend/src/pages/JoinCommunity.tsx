import React from "react";
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  Avatar,
  useTheme,
  Button,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
} from "@mui/material";
import { motion } from "framer-motion";
import {
  Groups,
  CheckCircle,
  TrendingUp,
  Security,
  Speed,
  SmartToy,
  ConnectWithoutContact,
} from "@mui/icons-material";

const JoinCommunity: React.FC = () => {
  const theme = useTheme();

  const benefits = [
    {
      title: "Verified BFSI Network",
      description: "Access verified BFSI agents and trusted financial partners",
      icon: CheckCircle,
      color: "success",
    },
    {
      title: "Early Access Updates",
      description: "Get early updates on new AI-powered tools and features",
      icon: TrendingUp,
      color: "info",
    },
    {
      title: "Industry Insights",
      description: "Stay informed with compliance updates and market trends",
      icon: Security,
      color: "warning",
    },
    {
      title: "Professional Network",
      description: "Expand your network within banking, finance, and insurance",
      icon: Groups,
      color: "primary",
    },
  ];

  const features = [
    {
      icon: SmartToy,
      title: "AI-Powered Agents",
      description: "All agents powered by Claude Sonnet 4.6 for intelligent interactions",
    },
    {
      icon: Security,
      title: "Secure Platform",
      description: "Bank-level security and compliance-ready infrastructure",
    },
    {
      icon: Speed,
      title: "High Performance",
      description: "Lightning-fast deployment and real-time processing",
    },
    {
      icon: ConnectWithoutContact,
      title: "Seamless Integration",
      description: "Connect with institutions and professionals effortlessly",
    },
  ];

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.2,
        delayChildren: 0.3,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.8, ease: "easeOut" },
    },
  };

  return (
    <Box sx={{ minHeight: "100vh", overflow: "hidden", position: "relative" }}>
      {/* Animated background elements */}
      <Box
        sx={{
          position: "fixed",
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          zIndex: -1,
          background: `linear-gradient(135deg, #0a0e27 0%, #0f1438 50%, #111533 100%)`,
        }}
      >
        <motion.div
          animate={{
            scale: [1, 1.1, 1],
            opacity: [0.3, 0.5, 0.3],
          }}
          transition={{ duration: 8, repeat: Infinity }}
          style={{
            position: "absolute",
            width: "500px",
            height: "500px",
            background:
              "radial-gradient(circle, rgba(168, 85, 247, 0.2) 0%, transparent 70%)",
            top: "-100px",
            left: "10%",
            borderRadius: "50%",
          }}
        />
        <motion.div
          animate={{
            scale: [1, 1.15, 1],
            opacity: [0.2, 0.4, 0.2],
          }}
          transition={{ duration: 10, repeat: Infinity, delay: 1 }}
          style={{
            position: "absolute",
            width: "600px",
            height: "600px",
            background:
              "radial-gradient(circle, rgba(6, 182, 212, 0.15) 0%, transparent 70%)",
            bottom: "-150px",
            right: "5%",
            borderRadius: "50%",
          }}
        />
      </Box>

      <Container maxWidth="lg">
        {/* Hero Section */}
        <Box
          sx={{
            py: 12,
            textAlign: "center",
            position: "relative",
            zIndex: 1,
          }}
        >
          <motion.div
            initial={{ opacity: 0, y: 50 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, ease: "easeOut" }}
          >
            <Typography
              variant="h1"
              sx={{
                mb: 3,
                fontSize: { xs: "2.5rem", sm: "3rem", md: "4rem" },
                fontWeight: 800,
                lineHeight: 1.1,
              }}
            >
              Join Our{" "}
              <Box
                component="span"
                sx={{
                  background: `linear-gradient(135deg, #a855f7 0%, #06b6d4 100%)`,
                  backgroundClip: "text",
                  WebkitBackgroundClip: "text",
                  WebkitTextFillColor: "transparent",
                }}
              >
                Community
              </Box>
            </Typography>

            <Typography
              variant="h5"
              sx={{
                mb: 4,
                color: "text.secondary",
                fontWeight: 400,
                fontSize: { xs: "1rem", md: "1.25rem" },
                maxWidth: 800,
                mx: "auto",
              }}
            >
              Be part of a growing ecosystem built exclusively for the BFSI sector. 
              Connect with industry professionals, discover new opportunities, and collaborate with trusted 
              financial agents and institutions.
            </Typography>

            <Typography
              variant="h6"
              sx={{
                mb: 2,
                color: "text.secondary",
                fontWeight: 400,
                maxWidth: 700,
                mx: "auto",
              }}
            >
              All our AI agents are powered by Claude Sonnet 4.6, ensuring intelligent, 
              secure, and high-performance interactions across the platform.
            </Typography>
          </motion.div>
        </Box>

        {/* Why Join Section */}
        <Box sx={{ py: 8, position: "relative", zIndex: 1 }}>
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
          >
            <Typography
              variant="h2"
              sx={{
                textAlign: "center",
                mb: 4,
                fontWeight: 700,
                fontSize: { xs: "2rem", md: "2.5rem" },
              }}
            >
              Why Join?
            </Typography>

            <Box
              sx={{
                maxWidth: 1000,
                mx: "auto",
                background: `linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(6, 182, 212, 0.05) 100%)`,
                border: `1px solid ${theme.palette.primary.main}30`,
                borderRadius: 3,
                p: 4,
                backdropFilter: "blur(10px)",
              }}
            >
              <List>
                {[
                  "Access verified BFSI agents and partners",
                  "Get early updates on new AI-powered tools",
                  "Stay informed with industry insights and compliance updates",
                  "Expand your professional network within banking, finance, and insurance",
                ].map((benefit, index) => (
                  <motion.div
                    key={index}
                    initial={{ opacity: 0, x: -20 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.6, delay: index * 0.1 }}
                    viewport={{ once: true }}
                  >
                    <ListItem
                      sx={{
                        py: 2,
                        px: 2,
                        borderRadius: 2,
                        mb: 1,
                        background: "rgba(255, 255, 255, 0.02)",
                        border: "1px solid rgba(255, 255, 255, 0.05)",
                        transition: "all 0.3s ease",
                        "&:hover": {
                          background: "rgba(168, 85, 247, 0.1)",
                          border: "1px solid rgba(168, 85, 247, 0.2)",
                          transform: "translateX(8px)",
                        },
                      }}
                    >
                      <ListItemIcon>
                        <CheckCircle
                          sx={{
                            color: "primary.main",
                            fontSize: 24,
                          }}
                        />
                      </ListItemIcon>
                      <ListItemText
                        primary={benefit}
                        primaryTypographyProps={{
                          sx: {
                            color: "text.primary",
                            fontWeight: 500,
                            fontSize: "1.1rem",
                          },
                        }}
                      />
                    </ListItem>
                  </motion.div>
                ))}
              </List>
            </Box>
          </motion.div>
        </Box>

        {/* Benefits Grid */}
        <Box sx={{ py: 8, position: "relative", zIndex: 1 }}>
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
          >
            <Typography
              variant="h2"
              sx={{
                textAlign: "center",
                mb: 6,
                fontWeight: 700,
                fontSize: { xs: "2rem", md: "2.5rem" },
              }}
            >
              Community Benefits
            </Typography>

            <motion.div
              variants={containerVariants}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true }}
            >
              <Grid container spacing={4}>
                {benefits.map((benefit, index) => {
                  const IconComponent = benefit.icon;
                  const colorKey = benefit.color as
                    | "primary"
                    | "info"
                    | "success"
                    | "warning";
                  return (
                    <Grid item xs={12} sm={6} md={3} key={index}>
                      <motion.div variants={itemVariants}>
                        <Card
                          sx={{
                            height: "100%",
                            position: "relative",
                            overflow: "hidden",
                            "&::before": {
                              content: '""',
                              position: "absolute",
                              inset: 0,
                              background: `linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(6, 182, 212, 0.05) 100%)`,
                              opacity: 0,
                              transition: "opacity 0.3s ease",
                            },
                            "&:hover::before": {
                              opacity: 1,
                            },
                            "&:hover": {
                              transform: "translateY(-8px)",
                            },
                          }}
                        >
                          <CardContent sx={{ py: 4 }}>
                            <motion.div
                              whileHover={{ rotate: 10, scale: 1.1 }}
                              transition={{ type: "spring", stiffness: 200 }}
                            >
                              <Box
                                sx={{
                                  display: "flex",
                                  alignItems: "center",
                                  justifyContent: "center",
                                  width: 60,
                                  height: 60,
                                  borderRadius: "12px",
                                  background: `linear-gradient(135deg, ${theme.palette[colorKey].main}20 0%, ${theme.palette[colorKey].main}10 100%)`,
                                  mb: 3,
                                  mx: "auto",
                                }}
                              >
                                <IconComponent
                                  sx={{
                                    fontSize: 32,
                                    color: `${colorKey}.main`,
                                  }}
                                />
                              </Box>
                            </motion.div>

                            <Typography
                              variant="h6"
                              sx={{ mb: 2, fontWeight: 700, textAlign: "center" }}
                            >
                              {benefit.title}
                            </Typography>
                            <Typography
                              variant="body2"
                              sx={{
                                color: "text.secondary",
                                textAlign: "center",
                                lineHeight: 1.7,
                              }}
                            >
                              {benefit.description}
                            </Typography>
                          </CardContent>
                        </Card>
                      </motion.div>
                    </Grid>
                  );
                })}
              </Grid>
            </motion.div>
          </motion.div>
        </Box>

        {/* Features Section */}
        <Box sx={{ py: 8, position: "relative", zIndex: 1 }}>
          <motion.div
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
          >
            <Typography
              variant="h2"
              sx={{
                textAlign: "center",
                mb: 6,
                fontWeight: 700,
                fontSize: { xs: "2rem", md: "2.5rem" },
              }}
            >
              Platform Features
            </Typography>

            <Grid container spacing={4}>
              {features.map((feature, index) => (
                <Grid item xs={12} sm={6} md={3} key={index}>
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    whileInView={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.6, delay: index * 0.1 }}
                    viewport={{ once: true }}
                  >
                    <Card
                      sx={{
                        height: "100%",
                        background:
                          "linear-gradient(135deg, rgba(255, 255, 255, 0.05) 0%, rgba(255, 255, 255, 0.02) 100%)",
                        border: "1px solid rgba(255, 255, 255, 0.1)",
                        borderRadius: 3,
                        transition: "all 0.3s ease",
                        "&:hover": {
                          transform: "translateY(-4px)",
                          boxShadow: "0 8px 32px rgba(152, 23, 126, 0.3)",
                          border: "1px solid rgba(152, 23, 126, 0.2)",
                        },
                      }}
                    >
                      <CardContent sx={{ p: 3 }}>
                        <Box display="flex" alignItems="center" mb={2}>
                          <Avatar
                            sx={{
                              width: 48,
                              height: 48,
                              mr: 2,
                              background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                            }}
                          >
                            <feature.icon sx={{ color: "white", fontSize: 24 }} />
                          </Avatar>
                          <Box>
                            <Typography variant="h6" fontWeight="bold">
                              {feature.title}
                            </Typography>
                          </Box>
                        </Box>
                        <Typography
                          variant="body2"
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
        </Box>

        {/* CTA Section */}
        <Box sx={{ py: 8, textAlign: "center", mb: 8 }}>
          <motion.div
            initial={{ opacity: 0, y: 50 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            viewport={{ once: true }}
          >
            <Typography
              variant="h2"
              sx={{
                mb: 3,
                fontWeight: 700,
                fontSize: { xs: "1.8rem", md: "2.5rem" },
              }}
            >
              Ready to Join?
            </Typography>
            <Typography
              variant="h6"
              sx={{
                mb: 4,
                color: "text.secondary",
                fontWeight: 400,
                maxWidth: 600,
                mx: "auto",
              }}
            >
              Become part of the exclusive BFSI community and transform your financial operations 
              with cutting-edge AI agents.
            </Typography>

            <Box
              sx={{
                display: "flex",
                gap: 2,
                justifyContent: "center",
                flexWrap: "wrap",
              }}
            >
              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Button
                  variant="contained"
                  size="large"
                  sx={{
                    px: 4,
                    py: 1.8,
                    fontSize: "1.1rem",
                    fontWeight: 700,
                    background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                  }}
                >
                  Join Community Now
                </Button>
              </motion.div>
              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Button
                  variant="outlined"
                  size="large"
                  sx={{
                    px: 4,
                    py: 1.8,
                    fontSize: "1.1rem",
                    fontWeight: 700,
                  }}
                >
                  Learn More
                </Button>
              </motion.div>
            </Box>
          </motion.div>
        </Box>
      </Container>
    </Box>
  );
};

export default JoinCommunity;

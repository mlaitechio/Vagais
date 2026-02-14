import React from "react";
import {
  Box,
  Container,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  useTheme,
  Stack,
} from "@mui/material";
import { motion } from "framer-motion";
import { useNavigate } from "react-router-dom";
import {
  FlashOn,
  SmartToy,
  Speed,
  Shield,
  ArrowRight,
} from "@mui/icons-material";
import AzureCloud from "../assets/azure_cloud.png";
import Footer from "../components/footert";
import Copilot from "../assets/copilot.png";
import Section from "../assets/ailogo.png";
const Landing: React.FC = () => {
  const navigate = useNavigate();
  const theme = useTheme();

  const features = [
    {
      title: "Neural Marketplace",
      description:
        "Build & Deploy cutting-edge AI agents powered by latest SLMs/LLMs",
      icon: SmartToy,
      color: "primary",
    },
    {
      title: "Lightning Fast",
      description:
        "Build, Execute & Deploye agents in milliseconds with optimized infrastructure and parallel processing",
      icon: FlashOn,
      color: "info",
    },
    {
      title: "Real-time Analytics",
      description:
        "Monitor agent performance with advanced metrics, usage patterns, and insights",
      icon: Speed,
      color: "success",
    },
    {
      title: "Enterprise Security",
      description:
        "Bank-level encryption, access controls, and compliance with industry standards",
      icon: Shield,
      color: "warning",
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
            style={{
              height: "65dvh",
              display: "flex",
              flexDirection: "column",
              justifyContent: "center",
            }}
          >
            {/* Badge */}
            <motion.div
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              {/* TOP: Powered by badges */}
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  alignItems: "center",
                  gap: 1,
                  marginBottom: "30px",
                }}
              >
                <Stack spacing={1} alignItems="center">
                  <Box className="iconBox"
                    sx={{
                      backgroundColor: "white",
                      p: "4px 12px",
                      borderRadius: 2,
                      display: "flex",
                      gap: 1,
                      alignItems: "center",
                      boxShadow: 1,
                    }}
                  >
                    <img src={AzureCloud} alt="Azure" />
                    <img src={Copilot} alt="Copilot" />
                  </Box>
                  <Typography
                    variant="caption"
                    sx={{
                      color: "text.secondary",
                      textTransform: "uppercase",
                      letterSpacing: 1,
                      fontWeight: 600,
                      fontSize: "16px",
                    }}
                  >
                    Powered by Azure AI Foundry and Microsoft CoPilot Studio
                  </Typography>
                  
                </Stack>
              </Box>
            </motion.div>

            {/* Main Heading */}
            <Typography className="max-heading"
              variant="h1"
              sx={{
                mb: 3,
                fontSize: { xs: "2.5rem", sm: "3rem", md: "4rem" },
                fontWeight: 800,
                lineHeight: 1.1,
              }}
            >
              The Future of{" "}
              <Box
                component="span"
                sx={{
                  background: `linear-gradient(135deg, #a855f7 0%, #06b6d4 100%)`,
                  backgroundClip: "text",
                  WebkitBackgroundClip: "text",
                  WebkitTextFillColor: "transparent",
                }}
              >
                Autonomous AI
              </Box>{" "}
              is Here
            </Typography>

            {/* Subheading */}
            <Typography className="deas"
              variant="h5"
              sx={{
                mb: 4,
                color: "text.secondary",
                fontWeight: 400,
                fontSize: { xs: "1rem", md: "1.25rem" },
                maxWidth: 600,
                mx: "auto",
              }}
            >
              Deploy, execute, and scale powerful AI agents in seconds.
              No coding required.
            </Typography>

            {/* CTA Buttons */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, delay: 0.3 }}
            >
              <Box
                sx={{
                  display: "flex",
                  gap: 2,
                  justifyContent: "center",
                  flexWrap: "wrap",
                  mt: 4,
                }}
              >
                <motion.div
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                >
                  <Button
                    variant="contained"
                    size="large"
                    endIcon={<ArrowRight />}
                    onClick={() => navigate("/marketplace")}
                    sx={{
                      px: 4,
                      py: 1.8,
                      fontSize: "1.1rem",
                      fontWeight: 700,
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
                    onClick={() => navigate("/login")}
                    sx={{
                      px: 4,
                      py: 1.8,
                      fontSize: "1.1rem",
                      fontWeight: 700,
                    }}
                  >
                    Sign In
                  </Button>
                </motion.div>
              </Box>
            </motion.div>
          </motion.div>
        </Box>
        {/* Section Divider Wrapper (Parent) */}
        <Box
          sx={{
            width: "100vw",
            position: "relative",
            left: "50%",
            right: "50%",
            ml: "-50vw",
            mr: "-50vw",
            my: 0,
            overflow: "hidden",
          }}
        >
          {/* Section Divider (Child) */}
          <Box className="logo_container">
            <img src={Section}/>
            <Box className="mynone"
              sx={{
                width: "100%",
                height: { xs: 370, md: 370 },
                backgroundImage: ``,
                backgroundSize: "cover",
                backgroundPosition: "center",
                backgroundRepeat: "no-repeat",

                /* Optional polish */
                maskImage:
                  "linear-gradient(to right, transparent, black 10%, black 90%, transparent)",
                WebkitMaskImage:
                  "linear-gradient(to right, transparent, black 10%, black 90%, transparent)",
              }}/>
            </Box>
            
        </Box>
        {/* Features Section */}
        <Box sx={{ py: 12, position: "relative", zIndex: 1 }}>
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
              Why Choose{" "}
              <Box component="span" sx={{ color: "primary.main" }}>
                merv.one
              </Box>
            </Typography>
            <Typography
              variant="h6"
              sx={{
                textAlign: "center",
                mb: 8,
                color: "text.secondary",
                fontWeight: 400,
              }}
            >
              Built for modern businesses that demand speed, reliability, and
              innovation
            </Typography>
          </motion.div>

          <motion.div
            variants={containerVariants}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
          >
            <Grid container spacing={3}>
              {features.map((feature, index) => {
                const IconComponent = feature.icon;
                const colorKey = feature.color as
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
                            {feature.title}
                          </Typography>
                          <Typography
                            variant="body2"
                            sx={{
                              color: "text.secondary",
                              textAlign: "center",
                              lineHeight: 1.7,
                            }}
                          >
                            {feature.description}
                          </Typography>
                        </CardContent>
                      </Card>
                    </motion.div>
                  </Grid>
                );
              })}
            </Grid>
          </motion.div>
        </Box>

        {/* Stats Section */}
        <motion.div
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
        >
          <Box className="max-round"
            sx={{
              py: 8,
              px: 4,
              borderRadius: 3,
              background: `linear-gradient(135deg, rgba(168, 85, 247, 0.15) 0%, rgba(6, 182, 212, 0.1) 100%)`,
              border: `1px solid ${theme.palette.primary.main}30`,
              backdropFilter: "blur(10px)",
              mb: 8,
              textAlign: "center",
            }}
          >
            <Grid container spacing={4}>
              {[
                { number: "100+", label: "Active Agents" },
                { number: "100+", label: "Executions" },
                { number: "99.99%", label: "Uptime" },
                { number: "30+", label: "Large Corporates" },
              ].map((stat, index) => (
                <Grid item xs={12} sm={6} md={3} key={index}>
                  <motion.div
                    initial={{ opacity: 0, scale: 0.5 }}
                    whileInView={{ opacity: 1, scale: 1 }}
                    transition={{ duration: 0.6, delay: index * 0.1 }}
                    viewport={{ once: true }}
                  >
                    <Box>
                      <Typography
                        variant="h2"
                        sx={{ fontWeight: 800, color: "primary.main", mb: 1 }}
                      >
                        {stat.number}
                      </Typography>
                      <Typography
                        variant="body2"
                        sx={{ color: "text.secondary", fontWeight: 500 }}
                      >
                        {stat.label}
                      </Typography>
                    </Box>
                  </motion.div>
                </Grid>
              ))}
            </Grid>
          </Box>
        </motion.div>

        {/* Final CTA Section */}
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
              Ready to Automate?
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
              Join thousands of teams building the future with AI agents
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
                  onClick={() => navigate("/register")}
                  sx={{
                    px: 4,
                    py: 1.8,
                    fontSize: "1.1rem",
                    fontWeight: 700,
                  }}
                >
                  Start to Build Your First Agent
                </Button>
              </motion.div>
              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Button
                  variant="outlined"
                  size="large"
                  onClick={() => navigate("/marketplace")}
                  sx={{
                    px: 4,
                    py: 1.8,
                    fontSize: "1.1rem",
                    fontWeight: 700,
                  }}
                >
                  Browse Marketplace
                </Button>
              </motion.div>
            </Box>
          </motion.div>
          
        </Box>
      </Container>
      <Box>
        <Footer></Footer>
      </Box>
    </Box>
  );
};

export default Landing;

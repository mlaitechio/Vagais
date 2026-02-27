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
            py: { xs: 2, md: 4 },
            textAlign: "center",
            position: "relative",
            zIndex: 1,
            minHeight: { xs: "auto", md: "90vh" },
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <Container maxWidth="md">
            <motion.div
              initial={{ opacity: 0, y: 50 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, ease: "easeOut" }}
              style={{
                minHeight: "70vh",
                display: "flex",
                flexDirection: "column",
                justifyContent: "center",
                alignItems: "center",
                textAlign: "center",
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
                    marginBottom: "20px",
                    width: "100%",
                  }}
                >
                  <Stack spacing={1} alignItems="center">
                    <Box className="iconBox"
                      sx={{
                        backgroundColor: "white",
                        p: "8px 20px",
                        borderRadius: 2,
                        display: "flex",
                        gap: 1.5,
                        alignItems: "center",
                        boxShadow: 2,
                      }}
                    >
                      <img src={AzureCloud} alt="Azure" style={{ height: "32px" }} />
                      <img src={Copilot} alt="Copilot" style={{ height: "32px" }} />
                    </Box>
                    <Typography
                      variant="caption"
                      sx={{
                        color: "text.secondary",
                        textTransform: "uppercase",
                        letterSpacing: 1,
                        fontWeight: 600,
                        fontSize: { xs: "14px", md: "16px" },
                        textAlign: "center",
                      }}
                    >
                      Powered by Azure AI Foundry and Microsoft CoPilot Studio
                    </Typography>
                    
                  </Stack>
                </Box>
              </motion.div>

              {/* Main Heading */}
              <motion.div
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.8, delay: 0.3 }}
              >
                <Typography className="max-heading"
                  variant="h1"
                  sx={{
                    mb: 4,
                    fontSize: { xs: "2.5rem", sm: "3rem", md: "3.5rem", lg: "4rem" },
                    fontWeight: 900,
                    lineHeight: 1.1,
                    textAlign: "center",
                    width: "100%",
                  }}
                >
                  <motion.div
                    initial={{ opacity: 0, x: -50 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.6, delay: 0.4 }}
                  >
                    <Box
                      component="span"
                      sx={{
                        fontSize: { xs: "1.2rem", sm: "1.5rem", md: "1.7rem", lg: "2rem" },
                        fontWeight: 600,
                        display: "block",
                        mb: 1,
                        textAlign: "center",
                      }}
                    >
                      The Future of
                    </Box>
                  </motion.div>
                  <motion.div
                    initial={{ opacity: 0, scale: 0.8 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ duration: 0.8, delay: 0.6 }}
                    whileHover={{ scale: 1.05 }}
                  >
                    <Box
                      component="span"
                      sx={{
                        background: `linear-gradient(135deg, #a855f7 0%, #06b6d4 100%)`,
                        backgroundClip: "text",
                        WebkitBackgroundClip: "text",
                        WebkitTextFillColor: "transparent",
                        fontSize: { xs: "2.5rem", sm: "3rem", md: "3.5rem", lg: "4rem" },
                        fontWeight: 900,
                        display: "block",
                        lineHeight: 1,
                        mb: 1,
                        textAlign: "center",
                      }}
                    >
                      Autonomous AI is Here
                    </Box>
                  </motion.div>
                </Typography>
              </motion.div>

              {/* Subheading */}
              <Typography className="deas"
                variant="h5"
                sx={{
                  mb: 5,
                  color: "text.secondary",
                  fontWeight: 400,
                  fontSize: { xs: "1.1rem", md: "1.3rem", lg: "1.5rem" },
                  maxWidth: { xs: "100%", md: "700px", lg: "800px" },
                  mx: "auto",
                  lineHeight: 1.6,
                  textAlign: "center",
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
                    gap: 3,
                    justifyContent: "center",
                    flexWrap: "wrap",
                    mt: 4,
                    width: "100%",
                    alignItems: "center",
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
                        px: { xs: 3, md: 4 },
                        py: { xs: 1.5, md: 1.8 },
                        fontSize: { xs: "1rem", md: "1.1rem" },
                        fontWeight: 700,
                        borderRadius: 2,
                        textTransform: "none",
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
                        px: { xs: 3, md: 4 },
                        py: { xs: 1.5, md: 1.8 },
                        fontSize: { xs: "1rem", md: "1.1rem" },
                        fontWeight: 700,
                        borderRadius: 2,
                        textTransform: "none",
                      }}
                    >
                      Sign In
                    </Button>
                  </motion.div>
                </Box>
              </motion.div>
            </motion.div>
          </Container>
        </Box>
        {/* Section Divider Wrapper (Parent) */}
        <Box
          sx={{
            width: "100%",
            position: "relative",
            my: 6,
            overflow: "hidden",
          }}
        >
          {/* Section Divider (Child) */}
          <Box className="logo_container">
            <img src={Section}/>
            <Box className="mynone"
              sx={{
                width: "100%",
                height: { xs: 150, md: 180 },
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
        <Box sx={{ py: { xs: 8, md: 12 }, position: "relative", zIndex: 1 }}>
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
                fontSize: { xs: "1.8rem", sm: "2rem", md: "2.5rem", lg: "3rem" },
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
                fontSize: { xs: "1rem", md: "1.1rem", lg: "1.2rem" },
                maxWidth: { xs: "100%", md: "80%" },
                mx: "auto",
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
              py: { xs: 6, md: 8 },
              px: { xs: 2, md: 4 },
              borderRadius: 3,
              background: `linear-gradient(135deg, rgba(168, 85, 247, 0.15) 0%, rgba(6, 182, 212, 0.1) 100%)`,
              border: `1px solid ${theme.palette.primary.main}30`,
              backdropFilter: "blur(10px)",
              mb: { xs: 6, md: 8 },
              textAlign: "center",
            }}
          >
            <Grid container spacing={{ xs: 2, md: 4 }}>
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
                        sx={{ 
                          fontWeight: 800, 
                          color: "primary.main", 
                          mb: 1,
                          fontSize: { xs: "1.8rem", md: "2.5rem", lg: "3rem" }
                        }}
                      >
                        {stat.number}
                      </Typography>
                      <Typography
                        variant="body2"
                        sx={{ 
                          color: "text.secondary", 
                          fontWeight: 500,
                          fontSize: { xs: "0.9rem", md: "1rem" }
                        }}
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
        <Box sx={{ py: { xs: 6, md: 8 }, textAlign: "center", mb: { xs: 4, md: 8 } }}>
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
                fontSize: { xs: "1.8rem", md: "2.5rem", lg: "3rem" },
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
                maxWidth: { xs: "100%", md: "600px" },
                mx: "auto",
                fontSize: { xs: "1rem", md: "1.1rem" },
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
    </Box>
  );
};

export default Landing;

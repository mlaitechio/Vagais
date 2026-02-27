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
  Stack,
} from "@mui/material";
import { motion } from "framer-motion";
import {
  SmartToy,
  FlashOn,
  Shield,
  Speed,
  Groups,
  Lightbulb,
} from "@mui/icons-material";

const AboutUs: React.FC = () => {
  const theme = useTheme();

  const values = [
    {
      title: "Innovation",
      description: "Pushing the boundaries of AI technology to create transformative solutions",
      icon: Lightbulb,
      color: "primary",
    },
    {
      title: "Excellence",
      description: "Delivering world-class AI agents with unmatched performance and reliability",
      icon: FlashOn,
      color: "info",
    },
    {
      title: "Security",
      description: "Enterprise-grade security with bank-level encryption and compliance",
      icon: Shield,
      color: "success",
    },
    {
      title: "Speed",
      description: "Lightning-fast deployment and execution with optimized infrastructure",
      icon: Speed,
      color: "warning",
    },
  ];

  const stats = [
    { number: "100+", label: "AI Agents Deployed" },
    { number: "99.99%", label: "Uptime SLA" },
    { number: "30+", label: "Enterprise Clients" },
    { number: "24/7", label: "Support Available" },
  ];

  const team = [
    {
      name: "AI Innovation Team",
      role: "Research & Development",
      description: "Leading AI researchers and engineers pushing autonomous agents forward",
    },
    {
      name: "Security Experts",
      role: "Trust & Safety",
      description: "Dedicated team ensuring enterprise-grade security and compliance standards",
    },
    {
      name: "Customer Success",
      role: "Client Support",
      description: "24/7 support team helping clients maximize AI agent investments",
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
    <Box sx={{ minHeight: "100vh", overflow: "hidden", position: "relative", display: "flex", flexDirection: "column", marginBottom: "50px" }}>
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

      <Container maxWidth="lg" sx={{ flex: 1, pb: 2 }}>
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
              About{" "}
              <Box
                component="span"
                sx={{
                  background: `linear-gradient(135deg, #a855f7 0%, #06b6d4 100%)`,
                  backgroundClip: "text",
                  WebkitBackgroundClip: "text",
                  WebkitTextFillColor: "transparent",
                }}
              >
                merv.one
              </Box>
            </Typography>

            <Typography
              variant="h5"
              sx={{
                mb: 4,
                color: "text.secondary",
                fontWeight: 400,
                fontSize: { xs: "1rem", md: "1.25rem" },
                maxWidth: 700,
                mx: "auto",
              }}
            >
              Empowering businesses with cutting-edge AI agents that transform how work gets done. 
              We're building the future of autonomous AI, one agent at a time.
            </Typography>
          </motion.div>
        </Box>

        {/* Stats Section */}
        <motion.div
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          transition={{ duration: 0.8 }}
          viewport={{ once: true }}
        >
          <Box
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
              {stats.map((stat, index) => (
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

        {/* Our Story Section */}
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
              Our Story
            </Typography>
            <Typography
              variant="h6"
              sx={{
                textAlign: "center",
                mb: 4,
                color: "text.secondary",
                fontWeight: 400,
                maxWidth: 800,
                mx: "auto",
                lineHeight: 1.7,
              }}
            >
              We are a part of MLAI Digital, founded in 2017, and we are a B2B-based company. 
              merv.one emerged from our vision to make AI agents accessible to every business. 
              Our team of AI researchers, engineers, and industry experts came together to solve the complex challenge 
              of deploying and managing AI agents at scale. Today, we're proud to power some of the most innovative 
              companies in fintech, banking, and financial services.
            </Typography>
          </motion.div>
        </Box>

        {/* Platform Expertise Section */}
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
              Built for BFSI Excellence
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
              <Typography
                variant="h6"
                sx={{
                  color: "text.secondary",
                  fontWeight: 400,
                  lineHeight: 1.8,
                  mb: 3,
                }}
              >
                Our platform is built on Microsoft Azure Platform by a dedicated team of developers with deep DOMAIN Knowledge in BFSI (Banking, Financial Services, and Insurance) sector. They have engineered a secure, scalable, and intelligent agent marketplace tailored specifically for financial institutions and service providers.
              </Typography>
              <Typography
                variant="h6"
                sx={{
                  color: "text.secondary",
                  fontWeight: 400,
                  lineHeight: 1.8,
                  mb: 3,
                }}
              >
                All our AI agents are powered by AZURE AI Foundry and realising up-to-date LLMs like GPT 5.2, Claude Sonnet 4.6, delivering advanced intelligence.
              </Typography>
              <Typography
                variant="h6"
                sx={{
                  color: "text.secondary",
                  fontWeight: 400,
                  lineHeight: 1.8,
                  mb: 3,
                }}
              >
                Leveraging modern technologies and industry best practices, they have created a robust ecosystem that connects BFSI agents which helps in:
              </Typography>
              <Box sx={{ mb: 3 }}>
                {[
                  "FINOPS",
                  "Sales Enhancements",
                  "Customer Onboarding and Enhancing Customer Experience",
                  "Security and Governance",
                ].map((useCase, index) => (
                  <Box
                    key={index}
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      mb: 2,
                      px: 2,
                      py: 1,
                      background: "rgba(255, 255, 255, 0.03)",
                      borderRadius: 2,
                      border: "1px solid rgba(255, 255, 255, 0.08)",
                      transition: "all 0.3s ease",
                      "&:hover": {
                        background: "rgba(168, 85, 247, 0.1)",
                        border: "1px solid rgba(168, 85, 247, 0.2)",
                        transform: "translateX(8px)",
                      },
                    }}
                  >
                    <Box
                      sx={{
                        width: 8,
                        height: 8,
                        borderRadius: "50%",
                        background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                        mr: 2,
                      }}
                    />
                    <Typography
                      variant="body1"
                      sx={{
                        color: "text.primary",
                        fontWeight: 500,
                        fontSize: "1rem",
                      }}
                    >
                      {useCase}
                    </Typography>
                  </Box>
                ))}
              </Box>
              <Typography
                variant="h6"
                sx={{
                  color: "text.secondary",
                  fontWeight: 400,
                  lineHeight: 1.8,
                }}
              >
                These Agents further make use of MCP thereby empowering the BFSI ecosystem to grow and innovate confidently.
              </Typography>
            </Box>
          </motion.div>
        </Box>

        {/* Core Values Section */}
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
                mb: 8,
                fontWeight: 700,
                fontSize: { xs: "2rem", md: "2.5rem" },
              }}
            >
              Our Core Values
            </Typography>

            <motion.div
              variants={containerVariants}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true }}
            >
              <Grid container spacing={3}>
                {values.map((value, index) => {
                  const IconComponent = value.icon;
                  const colorKey = value.color as
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
                              {value.title}
                            </Typography>
                            <Typography
                              variant="body2"
                              sx={{
                                color: "text.secondary",
                                textAlign: "center",
                                lineHeight: 1.7,
                              }}
                            >
                              {value.description}
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

        {/* Team Section */}
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
                mb: 8,
                fontWeight: 700,
                fontSize: { xs: "2rem", md: "2.5rem" },
              }}
            >
              Our Team
            </Typography>

            <Grid container spacing={4}>
              {team.map((member, index) => (
                <Grid item xs={12} md={4} key={index}>
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
                            <Groups />
                          </Avatar>
                          <Box>
                            <Typography variant="h6" fontWeight="bold">
                              {member.name}
                            </Typography>
                            <Typography variant="body2" color="primary">
                              {member.role}
                            </Typography>
                          </Box>
                        </Box>
                        <Typography
                          variant="body2"
                          color="text.secondary"
                          sx={{ lineHeight: 1.6 }}
                        >
                          {member.description}
                        </Typography>
                      </CardContent>
                    </Card>
                  </motion.div>
                </Grid>
              ))}
            </Grid>
          </motion.div>
        </Box>

        {/* Mission Section */}
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
              Our Mission
            </Typography>
            <Typography
              variant="h6"
              sx={{
                mb: 4,
                color: "text.secondary",
                fontWeight: 400,
                maxWidth: 800,
                mx: "auto",
                lineHeight: 1.7,
              }}
            >
              To democratize AI technology by making powerful AI agents accessible, deployable, and manageable 
              for businesses of all sizes. We believe the future of work is autonomous, and we're building 
              the tools to make that future a reality today.
            </Typography>
          </motion.div>
        </Box>
      </Container>
    </Box>
  );
};

export default AboutUs;

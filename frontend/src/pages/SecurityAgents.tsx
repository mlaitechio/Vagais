import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  useTheme,
} from "@mui/material";
import { Theme } from "@mui/material/styles";
import { motion } from "framer-motion";
import SecurityIcon from "@mui/icons-material/Security";
import { SecurityOutlined } from "@mui/icons-material";

const containerVariants = {
  hidden: {},
  visible: {
    transition: {
      staggerChildren: 0.2,
    },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.6 },
  },
};

const SecurityAgents = () => {
  const theme = useTheme();

  const features = [
    {
      icon: SecurityOutlined,
      title: "Access Review Agent",
      description:
        "Empower your reviewers to make fast, accurate access decisions. The Access Review Agent delivers insights and recommendations so reviewers can complete their work through a simple conversation, right inside Microsoft Teams.",
      color: "primary",
    },
    {
      icon: SecurityIcon,
      title: "Phishing Triage Agent in Microsoft Defender",
      description:
        "Designed to scale security teams' response in triaging and classifying user-submitted phishing incidents, allowing organizations to improve their efficiency by reducing manual effort and streamlining their phishing response.",
      color: "info",
    },
    {
      icon: SecurityIcon,
      title: "Threat Intelligence Briefing Agent in Security Copilot",
      description:
        "Automatically curates relevant and timely threat intelligence based on an organization's unique attributes and threat exposure.",
      color: "warning",
    },
    {
      icon: SecurityOutlined,
      title: "Vulnerability Remediation Agent in Microsoft Intune",
      description:
        "Identify top vulnerabilities, understand their impact, and get step-by-step remediation guidance to fix vulnerabilities using Intune capabilities.",
      color: "success",
    },
     {
      icon: SecurityOutlined,
      title: "Vulnerability Remediation Agent",
      description:
        "Identify top vulnerabilities, understand their impact, and get step-by-step remediation guidance to fix vulnerabilities using Intune capabilities.",
      color: "success",
    },
  ];

  return (
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
            fontSize: {
              xs: "1.8rem",
              sm: "2rem",
              md: "2.5rem",
              lg: "3rem",
            },
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

      {/* 🔥 alignItems="stretch" added */}
      <motion.div
        variants={containerVariants}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true }}
      >
        <Grid container spacing={3} alignItems="stretch">
          {features.map((feature, index) => {
            const IconComponent = feature.icon;
            // const colorKey = feature.color;
            const colorKey: keyof Theme["palette"] = "primary";

            return (
              <Grid item xs={12} sm={6} md={3} key={index}>
                <motion.div variants={itemVariants} style={{ height: "100%" }}>
                  <Card
                    sx={{
                      height: "100%", // 🔥 important
                      display: "flex",
                      flexDirection: "column",
                      position: "relative",
                      overflow: "hidden",
                      transition: "all 0.3s ease",
                      "&::before": {
                        content: '""',
                        position: "absolute",
                        inset: 0,
                        background:
                          "linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(6, 182, 212, 0.05) 100%)",
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
                    <CardContent
                      sx={{
                        py: 4,
                        display: "flex",
                        flexDirection: "column",
                        height: "100%",
                      }}
                    >
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
                            background: `linear-gradient(
      135deg, 
      ${theme.palette[colorKey].main}20 0%, 
      ${theme.palette[colorKey].main}10 100%
    )`,
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

                      {/* 🔥 FIXED HEIGHT DESCRIPTION */}
                      <Typography
                        variant="body2"
                        sx={{
                          color: "text.secondary",
                          textAlign: "center",
                          lineHeight: 1.7,
                          minHeight: "130px", // 👈 this keeps cards equal
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
  );
};

export default SecurityAgents;

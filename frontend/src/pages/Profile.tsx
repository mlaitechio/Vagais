import React from "react";
import { Box, Card, CardContent, Typography } from "@mui/material";
import ConstructionIcon from "@mui/icons-material/Construction";

const ProfileComingSoon: React.FC = () => {
  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        background: "linear-gradient(135deg, #667eea, #764ba2)",
        p: 2,
      }}
    >
      <Card
        elevation={6}
        sx={{
          maxWidth: 420,
          width: "100%",
          textAlign: "center",
          borderRadius: 3,
        }}
      >
        <CardContent sx={{ p: 4 }}>
          <ConstructionIcon
            color="primary"
            sx={{ fontSize: 50, mb: 2 }}
          />

          <Typography variant="h4" fontWeight="bold" gutterBottom>
            Working On This
          </Typography>

          <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
            This profile page is currently under development.
          </Typography>

          <Typography
            variant="body1"
            sx={{
              fontWeight: 600,
              color: "primary.main",
            }}
          >
            Coming Soon with Amazing Features ✨
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
};

export default ProfileComingSoon;

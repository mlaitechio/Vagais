import { createTheme, ThemeOptions } from '@mui/material/styles';

// Modern AI-Centric Color Palette
const aiColors = {
  // Primary - Electric purple/magenta
  primary: {
    50: '#faf5ff',
    100: '#f3e8ff',
    200: '#e9d5ff',
    300: '#d8b4fe',
    400: '#c084fc',
    500: '#a855f7',
    600: '#9333ea',
    700: '#7e22ce',
    800: '#6b21a8',
    900: '#581c87',
  },
  // Secondary - Cyan/Electric blue
  secondary: {
    50: '#ecf0ff',
    100: '#dde9ff',
    200: '#c7d2ff',
    300: '#a5b4ff',
    400: '#7c8aff',
    500: '#5b6bff',
    600: '#4f46e5',
    700: '#4338ca',
    800: '#342e8f',
    900: '#27205f',
  },
  // Accent - Neon cyan
  accent: {
    50: '#f0f9ff',
    100: '#e0f2fe',
    200: '#bae6fd',
    300: '#7dd3fc',
    400: '#38bdf8',
    500: '#06b6d4',
    600: '#0891b2',
    700: '#0e7490',
    800: '#155e75',
    900: '#164e63',
  },
  // Background - Deep space
  background: {
    default: '#0a0e27', // Deep space blue
    paper: '#111533',   // Darker overlay
    elevated: '#1a1f3a', // Elevated surfaces
  },
  // Text colors - High contrast
  text: {
    primary: '#f0f4f8',   // Almost white
    secondary: '#a1aac4', // Cool gray
    disabled: '#6b7280',  // Muted
  },
  // Neon status colors
  success: '#10b981',    // Neon green
  warning: '#f59e0b',    // Neon amber
  error: '#ef4444',      // Neon red
  info: '#06b6d4',       // Neon cyan
};

const aiTheme: ThemeOptions = {
  palette: {
    mode: 'dark',
    primary: {
      main: aiColors.primary[500],
      light: aiColors.primary[400],
      dark: aiColors.primary[600],
      contrastText: '#ffffff',
    },
    secondary: {
      main: aiColors.secondary[600],
      light: aiColors.secondary[400],
      dark: aiColors.secondary[700],
      contrastText: '#ffffff',
    },
    background: {
      default: aiColors.background.default,
      paper: aiColors.background.paper,
    },
    text: {
      primary: aiColors.text.primary,
      secondary: aiColors.text.secondary,
    },
    success: {
      main: aiColors.success,
      light: '#34d399',
      dark: '#059669',
    },
    warning: {
      main: aiColors.warning,
      light: '#fbbf24',
      dark: '#d97706',
    },
    error: {
      main: aiColors.error,
      light: '#f87171',
      dark: '#dc2626',
    },
    info: {
      main: aiColors.info,
      light: '#22d3ee',
      dark: '#0891b2',
    },
    divider: 'rgba(255, 255, 255, 0.08)',
  },
  typography: {
    fontFamily: '"Geist Sans", "Inter", "Roboto", "-apple-system", "BlinkMacSystemFont", "Segoe UI", sans-serif',
    h1: {
      fontSize: '3.5rem',
      fontWeight: 800,
      lineHeight: 1.1,
      letterSpacing: '-0.03em',
      background: `linear-gradient(135deg, ${aiColors.primary[500]} 0%, ${aiColors.accent[500]} 100%)`,
      backgroundClip: 'text',
      WebkitBackgroundClip: 'text',
      WebkitTextFillColor: 'transparent',
    },
    h2: {
      fontSize: '2.5rem',
      fontWeight: 700,
      lineHeight: 1.2,
      letterSpacing: '-0.02em',
    },
    h3: {
      fontSize: '2rem',
      fontWeight: 700,
      lineHeight: 1.3,
      letterSpacing: '-0.01em',
    },
    h4: {
      fontSize: '1.5rem',
      fontWeight: 600,
      lineHeight: 1.4,
    },
    h5: {
      fontSize: '1.25rem',
      fontWeight: 600,
      lineHeight: 1.4,
    },
    h6: {
      fontSize: '1rem',
      fontWeight: 600,
      lineHeight: 1.5,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    body1: {
      fontSize: '1rem',
      lineHeight: 1.6,
      letterSpacing: '0.01em',
    },
    body2: {
      fontSize: '0.875rem',
      lineHeight: 1.6,
      letterSpacing: '0.01em',
    },
    button: {
      fontWeight: 700,
      textTransform: 'none',
      letterSpacing: '0.03em',
      fontSize: '0.95rem',
    },
  },
  shape: {
    borderRadius: 16,
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          background: `linear-gradient(135deg, ${aiColors.background.default} 0%, #0f1438 50%, ${aiColors.background.paper} 100%)`,
          minHeight: '100vh',
          position: 'relative',
          backgroundAttachment: 'fixed',
        },
        '*': {
          scrollbarWidth: 'thin',
          scrollbarColor: `${aiColors.primary[500]} ${aiColors.background.paper}`,
        },
        '*::-webkit-scrollbar': {
          width: '8px',
        },
        '*::-webkit-scrollbar-track': {
          background: aiColors.background.paper,
        },
        '*::-webkit-scrollbar-thumb': {
          background: aiColors.primary[500],
          borderRadius: '4px',
        },
        '*::-webkit-scrollbar-thumb:hover': {
          background: aiColors.primary[400],
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 10,
          padding: '12px 28px',
          fontSize: '0.95rem',
          fontWeight: 700,
          textTransform: 'none',
          boxShadow: 'none',
          position: 'relative',
          overflow: 'hidden',
          transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
          '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: '-100%',
            width: '100%',
            height: '100%',
            background: 'rgba(255, 255, 255, 0.1)',
            transition: 'left 0.3s ease',
            zIndex: 0,
          },
          '&:hover::before': {
            left: '100%',
          },
          '& span': {
            position: 'relative',
            zIndex: 1,
          },
        },
        contained: {
          background: `linear-gradient(135deg, ${aiColors.primary[500]} 0%, ${aiColors.primary[600]} 100%)`,
          color: '#ffffff',
          boxShadow: `0 8px 32px rgba(168, 85, 247, 0.3)`,
          '&:hover': {
            boxShadow: `0 16px 48px rgba(168, 85, 247, 0.4)`,
            transform: 'translateY(-2px)',
          },
          '&:active': {
            transform: 'translateY(0)',
          },
        },
        outlined: {
          borderColor: aiColors.primary[500],
          color: aiColors.primary[400],
          border: `2px solid ${aiColors.primary[500]}`,
          '&:hover': {
            borderColor: aiColors.primary[400],
            background: 'rgba(168, 85, 247, 0.1)',
            boxShadow: `0 8px 24px rgba(168, 85, 247, 0.2)`,
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          background: `linear-gradient(135deg, rgba(17, 21, 51, 0.8) 0%, rgba(26, 31, 58, 0.8) 100%)`,
          border: `1px solid rgba(168, 85, 247, 0.2)`,
          backdropFilter: 'blur(20px)',
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.1)',
          position: 'relative',
          overflow: 'hidden',
          transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
          '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            height: '1px',
            background: `linear-gradient(90deg, transparent, rgba(168, 85, 247, 0.5), transparent)`,
          },
          '&:hover': {
            border: `1px solid rgba(168, 85, 247, 0.4)`,
            boxShadow: '0 16px 48px rgba(168, 85, 247, 0.2), inset 0 1px 0 rgba(255, 255, 255, 0.1)',
          },
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.05) 0%, rgba(168, 85, 247, 0.05) 100%)',
            borderRadius: 12,
            border: `2px solid rgba(168, 85, 247, 0.2)`,
            transition: 'all 0.3s ease',
            '& fieldset': {
              border: 'none',
            },
            '&:hover': {
              border: `2px solid rgba(168, 85, 247, 0.4)`,
              background: 'linear-gradient(135deg, rgba(255, 255, 255, 0.08) 0%, rgba(168, 85, 247, 0.08) 100%)',
            },
            '&.Mui-focused': {
              border: `2px solid ${aiColors.primary[500]}`,
              boxShadow: `0 0 20px rgba(168, 85, 247, 0.3)`,
            },
          },
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          background: `linear-gradient(90deg, rgba(10, 14, 39, 0.9) 0%, rgba(17, 21, 51, 0.9) 100%)`,
          border: `1px solid rgba(168, 85, 247, 0.2)`,
          backdropFilter: 'blur(20px)',
          boxShadow: '0 4px 24px rgba(0, 0, 0, 0.3)',
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          background: `linear-gradient(135deg, rgba(10, 14, 39, 0.95) 0%, rgba(17, 21, 51, 0.95) 100%)`,
          border: `1px solid rgba(168, 85, 247, 0.2)`,
          backdropFilter: 'blur(10px)',
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: 10,
          margin: '4px 8px',
          transition: 'all 0.2s ease',
          '&:hover': {
            backgroundColor: 'rgba(168, 85, 247, 0.15)',
            boxShadow: 'inset 0 0 20px rgba(168, 85, 247, 0.1)',
          },
          '&.Mui-selected': {
            backgroundColor: 'rgba(168, 85, 247, 0.25)',
            borderLeft: `3px solid ${aiColors.primary[500]}`,
            paddingLeft: 'calc(16px - 3px)',
            '&:hover': {
              backgroundColor: 'rgba(168, 85, 247, 0.35)',
            },
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          background: `linear-gradient(135deg, ${aiColors.primary[600]} 0%, ${aiColors.primary[700]} 100%)`,
          color: '#ffffff',
          fontWeight: 600,
          border: `1px solid rgba(255, 255, 255, 0.2)`,
          '&:hover': {
            boxShadow: `0 4px 12px rgba(168, 85, 247, 0.3)`,
          },
        },
      },
    },
    MuiAlert: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          border: '1px solid rgba(255, 255, 255, 0.1)',
          backdropFilter: 'blur(10px)',
        },
        standardSuccess: {
          background: 'linear-gradient(135deg, rgba(16, 185, 129, 0.15) 0%, rgba(16, 185, 129, 0.05) 100%)',
          color: aiColors.success,
          borderColor: aiColors.success,
        },
        standardError: {
          background: 'linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(239, 68, 68, 0.05) 100%)',
          color: aiColors.error,
          borderColor: aiColors.error,
        },
        standardWarning: {
          background: 'linear-gradient(135deg, rgba(245, 158, 11, 0.15) 0%, rgba(245, 158, 11, 0.05) 100%)',
          color: aiColors.warning,
          borderColor: aiColors.warning,
        },
        standardInfo: {
          background: 'linear-gradient(135deg, rgba(6, 182, 212, 0.15) 0%, rgba(6, 182, 212, 0.05) 100%)',
          color: aiColors.info,
          borderColor: aiColors.info,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
      },
    },
  },
};

export const theme = createTheme(aiTheme);
export default theme;

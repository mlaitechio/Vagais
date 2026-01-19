import { createTheme, ThemeOptions } from '@mui/material/styles';

// Night Owl color palette
const nightOwlColors = {
  // Primary colors - Deep purple/indigo
  primary: {
    50: '#f3e5f5',
    100: '#e1bee7',
    200: '#ce93d8',
    300: '#ba68c8',
    400: '#ab47bc',
    500: '#9c27b0',
    600: '#8e24aa',
    700: '#7b1fa2',
    800: '#6a1b9a',
    900: '#4a148c',
  },
  // Secondary colors - Bright cyan/teal
  secondary: {
    50: '#e0f2f1',
    100: '#b2dfdb',
    200: '#80cbc4',
    300: '#4db6ac',
    400: '#26a69a',
    500: '#009688',
    600: '#00897b',
    700: '#00796b',
    800: '#00695c',
    900: '#004d40',
  },
  // Accent colors - Orange/amber for highlights
  accent: {
    50: '#fff8e1',
    100: '#ffecb3',
    200: '#ffe082',
    300: '#ffd54f',
    400: '#ffca28',
    500: '#ffc107',
    600: '#ffb300',
    700: '#ffa000',
    800: '#ff8f00',
    900: '#ff6f00',
  },
  // Background colors - Dark theme
  background: {
    default: '#011627', // Very dark blue-gray
    paper: '#0d1117',   // Slightly lighter dark
    elevated: '#161b22', // Elevated surfaces
  },
  // Text colors
  text: {
    primary: '#d6deeb',   // Light gray-blue
    secondary: '#a2aabc', // Medium gray
    disabled: '#6c757d',  // Muted gray
  },
  // Status colors
  success: '#7c3aed',    // Purple
  warning: '#f59e0b',    // Amber
  error: '#ef4444',      // Red
  info: '#3b82f6',       // Blue
};

const nightOwlTheme: ThemeOptions = {
  palette: {
    mode: 'dark',
    primary: {
      main: nightOwlColors.primary[500],
      light: nightOwlColors.primary[300],
      dark: nightOwlColors.primary[700],
      contrastText: '#ffffff',
    },
    secondary: {
      main: nightOwlColors.secondary[500],
      light: nightOwlColors.secondary[300],
      dark: nightOwlColors.secondary[700],
      contrastText: '#ffffff',
    },
    background: {
      default: nightOwlColors.background.default,
      paper: nightOwlColors.background.paper,
    },
    text: {
      primary: nightOwlColors.text.primary,
      secondary: nightOwlColors.text.secondary,
    },
    success: {
      main: nightOwlColors.success,
      light: '#a855f7',
      dark: '#5b21b6',
    },
    warning: {
      main: nightOwlColors.warning,
      light: '#fbbf24',
      dark: '#d97706',
    },
    error: {
      main: nightOwlColors.error,
      light: '#f87171',
      dark: '#dc2626',
    },
    info: {
      main: nightOwlColors.info,
      light: '#60a5fa',
      dark: '#2563eb',
    },
    divider: 'rgba(255, 255, 255, 0.12)',
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontSize: '2.5rem',
      fontWeight: 700,
      lineHeight: 1.2,
      letterSpacing: '-0.02em',
    },
    h2: {
      fontSize: '2rem',
      fontWeight: 600,
      lineHeight: 1.3,
      letterSpacing: '-0.01em',
    },
    h3: {
      fontSize: '1.75rem',
      fontWeight: 600,
      lineHeight: 1.3,
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
      fontSize: '1.125rem',
      fontWeight: 600,
      lineHeight: 1.4,
    },
    body1: {
      fontSize: '1rem',
      lineHeight: 1.6,
    },
    body2: {
      fontSize: '0.875rem',
      lineHeight: 1.6,
    },
    button: {
      fontWeight: 600,
      textTransform: 'none',
      letterSpacing: '0.02em',
    },
  },
  shape: {
    borderRadius: 12,
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          background: `linear-gradient(135deg, ${nightOwlColors.background.default} 0%, ${nightOwlColors.background.paper} 100%)`,
          minHeight: '100vh',
        },
        '*': {
          scrollbarWidth: 'thin',
          scrollbarColor: `${nightOwlColors.primary[500]} ${nightOwlColors.background.paper}`,
        },
        '*::-webkit-scrollbar': {
          width: '8px',
        },
        '*::-webkit-scrollbar-track': {
          background: nightOwlColors.background.paper,
        },
        '*::-webkit-scrollbar-thumb': {
          background: nightOwlColors.primary[500],
          borderRadius: '4px',
        },
        '*::-webkit-scrollbar-thumb:hover': {
          background: nightOwlColors.primary[400],
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          padding: '10px 24px',
          fontSize: '0.95rem',
          fontWeight: 600,
          textTransform: 'none',
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0 4px 12px rgba(156, 39, 176, 0.3)',
          },
        },
        contained: {
          background: `linear-gradient(135deg, ${nightOwlColors.primary[500]} 0%, ${nightOwlColors.primary[600]} 100%)`,
          '&:hover': {
            background: `linear-gradient(135deg, ${nightOwlColors.primary[400]} 0%, ${nightOwlColors.primary[500]} 100%)`,
          },
        },
        outlined: {
          borderColor: nightOwlColors.primary[500],
          color: nightOwlColors.primary[500],
          '&:hover': {
            borderColor: nightOwlColors.primary[400],
            backgroundColor: 'rgba(156, 39, 176, 0.08)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          background: `linear-gradient(135deg, ${nightOwlColors.background.paper} 0%, ${nightOwlColors.background.elevated} 100%)`,
          border: `1px solid rgba(255, 255, 255, 0.1)`,
          backdropFilter: 'blur(10px)',
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.3)',
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            background: 'rgba(255, 255, 255, 0.05)',
            borderRadius: 8,
            '& fieldset': {
              borderColor: 'rgba(255, 255, 255, 0.2)',
            },
            '&:hover fieldset': {
              borderColor: nightOwlColors.primary[500],
            },
            '&.Mui-focused fieldset': {
              borderColor: nightOwlColors.primary[500],
            },
          },
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          background: `linear-gradient(135deg, ${nightOwlColors.background.paper} 0%, ${nightOwlColors.background.elevated} 100%)`,
          borderBottom: `1px solid rgba(255, 255, 255, 0.1)`,
          backdropFilter: 'blur(10px)',
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          background: `linear-gradient(135deg, ${nightOwlColors.background.paper} 0%, ${nightOwlColors.background.elevated} 100%)`,
          borderRight: `1px solid rgba(255, 255, 255, 0.1)`,
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          margin: '4px 8px',
          '&:hover': {
            backgroundColor: 'rgba(156, 39, 176, 0.1)',
          },
          '&.Mui-selected': {
            backgroundColor: 'rgba(156, 39, 176, 0.2)',
            '&:hover': {
              backgroundColor: 'rgba(156, 39, 176, 0.3)',
            },
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          background: `linear-gradient(135deg, ${nightOwlColors.primary[500]} 0%, ${nightOwlColors.primary[600]} 100%)`,
          color: '#ffffff',
          fontWeight: 600,
        },
      },
    },
    MuiAlert: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          border: '1px solid rgba(255, 255, 255, 0.1)',
        },
        standardSuccess: {
          backgroundColor: 'rgba(124, 58, 237, 0.1)',
          color: nightOwlColors.success,
          borderColor: nightOwlColors.success,
        },
        standardError: {
          backgroundColor: 'rgba(239, 68, 68, 0.1)',
          color: nightOwlColors.error,
          borderColor: nightOwlColors.error,
        },
        standardWarning: {
          backgroundColor: 'rgba(245, 158, 11, 0.1)',
          color: nightOwlColors.warning,
          borderColor: nightOwlColors.warning,
        },
        standardInfo: {
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          color: nightOwlColors.info,
          borderColor: nightOwlColors.info,
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

export const theme = createTheme(nightOwlTheme);
export default theme;

import React, { useState } from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Avatar,
  IconButton,
  Badge,
  Box,
  Menu,
  MenuItem,
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Dashboard,
  Store,
  SmartToy,
  Person,
  Settings,
  Notifications,
  Search,
  Add,
  Logout,
  AccountCircle,
  Assessment,
  Payment,
  Security,
  Help,
  Close,
} from '@mui/icons-material';
import { motion } from 'framer-motion';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Navbar: React.FC = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleProfileMenuClose = () => {
    setAnchorEl(null);
  };

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const navigationItems = [
    { label: 'Dashboard', path: '/dashboard', icon: <Dashboard /> },
    { label: 'Marketplace', path: '/marketplace', icon: <Store /> },
    { label: 'My Agents', path: '/agents', icon: <SmartToy /> },
    ...(user?.role === 'admin' ? [{ label: 'Analytics', path: '/analytics', icon: <Assessment /> }] : []),
    { label: 'Billing', path: '/billing', icon: <Payment /> },
    { label: 'Settings', path: '/settings', icon: <Settings /> },
  ];

  const isActive = (path: string) => location.pathname === path;

  const Logo = () => (
    <Box display="flex" alignItems="center" gap={1}>
      <Box
        sx={{
          width: 40,
          height: 40,
          borderRadius: '50%',
          background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          border: '2px solid rgba(255, 255, 255, 0.2)',
        }}
      >
        <SmartToy sx={{ color: 'white', fontSize: 20 }} />
      </Box>
      <Typography
        variant="h6"
        fontWeight="bold"
        sx={{
          background: 'linear-gradient(135deg, #98177E 0%, #00D4FF 100%)',
          backgroundClip: 'text',
          WebkitBackgroundClip: 'text',
          WebkitTextFillColor: 'transparent',
        }}
      >
        vagais.ai
      </Typography>
    </Box>
  );

  const DesktopNavigation = () => (
    <Box display="flex" alignItems="center" gap={2} sx={{ ml: 4 }}>
      {navigationItems.map((item) => (
        <Button
          key={item.path}
          startIcon={item.icon}
          onClick={() => navigate(item.path)}
          sx={{
            color: isActive(item.path) ? 'primary.main' : 'text.primary',
            background: isActive(item.path) ? 'rgba(152, 23, 126, 0.1)' : 'transparent',
            borderRadius: 2,
            px: 2,
            py: 1,
            '&:hover': {
              background: 'rgba(152, 23, 126, 0.1)',
            },
          }}
        >
          {item.label}
        </Button>
      ))}
    </Box>
  );

  const MobileDrawer = () => (
    <Drawer
      variant="temporary"
      anchor="left"
      open={mobileOpen}
      onClose={handleDrawerToggle}
      ModalProps={{
        keepMounted: true,
      }}
      sx={{
        '& .MuiDrawer-paper': {
          width: 280,
          background: 'linear-gradient(135deg, rgba(26, 26, 26, 0.95) 0%, rgba(10, 10, 10, 0.95) 100%)',
          backdropFilter: 'blur(10px)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
        },
      }}
    >
      <Box sx={{ p: 2 }}>
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={3}>
          <Logo />
          <IconButton onClick={handleDrawerToggle} sx={{ color: 'text.primary' }}>
            <Close />
          </IconButton>
        </Box>
        
        <List>
          {navigationItems.map((item) => (
            <ListItem
              key={item.path}
              button
              onClick={() => {
                navigate(item.path);
                handleDrawerToggle();
              }}
              sx={{
                borderRadius: 2,
                mb: 1,
                background: isActive(item.path) ? 'rgba(152, 23, 126, 0.2)' : 'transparent',
                '&:hover': {
                  background: 'rgba(152, 23, 126, 0.1)',
                },
              }}
            >
              <ListItemIcon sx={{ color: isActive(item.path) ? 'primary.main' : 'text.secondary' }}>
                {item.icon}
              </ListItemIcon>
              <ListItemText
                primary={item.label}
                sx={{
                  color: isActive(item.path) ? 'primary.main' : 'text.primary',
                  fontWeight: isActive(item.path) ? 'bold' : 'normal',
                }}
              />
            </ListItem>
          ))}
        </List>

        <Divider sx={{ my: 2, borderColor: 'rgba(255, 255, 255, 0.1)' }} />

        <Box sx={{ p: 2 }}>
          <Box display="flex" alignItems="center" gap={2} mb={2}>
            <Avatar
              src={user?.avatar}
              sx={{
                width: 48,
                height: 48,
                background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
              }}
            >
              {user?.first_name?.[0] || 'U'}
            </Avatar>
            <Box>
              <Typography variant="subtitle1" fontWeight="bold">
                {user?.first_name} {user?.last_name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {user?.email}
              </Typography>
            </Box>
          </Box>

          <Button
            fullWidth
            variant="outlined"
            startIcon={<Logout />}
            onClick={handleLogout}
            sx={{ borderRadius: 2 }}
          >
            Logout
          </Button>
        </Box>
      </Box>
    </Drawer>
  );

  return (
    <>
      <AppBar
        position="fixed"
        sx={{
          background: 'linear-gradient(135deg, rgba(26, 26, 26, 0.95) 0%, rgba(10, 10, 10, 0.95) 100%)',
          backdropFilter: 'blur(10px)',
          borderBottom: '1px solid rgba(255, 255, 255, 0.1)',
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.3)',
        }}
      >
        <Toolbar>
          {/* Logo */}
          <Box
            component={motion.div}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            sx={{ cursor: 'pointer' }}
            onClick={() => navigate('/dashboard')}
          >
            <Logo />
          </Box>

          {/* Desktop Navigation */}
          {!isMobile && <DesktopNavigation />}

          {/* Spacer */}
          <Box sx={{ flexGrow: 1 }} />

          {/* Right Side Actions */}
          <Box display="flex" alignItems="center" gap={1}>
            {/* Search */}
            <IconButton
              sx={{
                color: 'text.secondary',
                '&:hover': { color: 'primary.main' },
              }}
            >
              <Search />
            </IconButton>

            {/* Notifications */}
            <Badge badgeContent={3} color="error">
              <IconButton
                sx={{
                  color: 'text.secondary',
                  '&:hover': { color: 'primary.main' },
                }}
              >
                <Notifications />
              </IconButton>
            </Badge>

            {/* Create Agent - Only for Admin */}
            {user?.role === 'admin' && (
              <Button
                variant="contained"
                startIcon={<Add />}
                onClick={() => navigate('/agents/create')}
                sx={{
                  background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                  borderRadius: 2,
                  px: 2,
                  mr: 2,
                  '&:hover': {
                    background: `linear-gradient(135deg, ${theme.palette.primary.dark}, ${theme.palette.secondary.dark})`,
                  },
                }}
              >
                Create Agent
              </Button>
            )}

            {/* Mobile Menu Button */}
            {isMobile ? (
              <IconButton
                color="inherit"
                aria-label="open drawer"
                edge="start"
                onClick={handleDrawerToggle}
              >
                <MenuIcon />
              </IconButton>
            ) : (
              /* Desktop Profile Menu */
              <Box display="flex" alignItems="center" gap={2}>
                <Box textAlign="right" sx={{ display: { xs: 'none', sm: 'block' } }}>
                  <Typography variant="body2" fontWeight="medium">
                    {user?.first_name} {user?.last_name}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    {user?.role}
                  </Typography>
                </Box>
                <Avatar
                  src={user?.avatar}
                  onClick={handleProfileMenuOpen}
                  sx={{
                    cursor: 'pointer',
                    width: 40,
                    height: 40,
                    background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
                    '&:hover': {
                      transform: 'scale(1.1)',
                    },
                  }}
                >
                  {user?.first_name?.[0] || 'U'}
                </Avatar>
              </Box>
            )}
          </Box>
        </Toolbar>
      </AppBar>

      {/* Mobile Drawer */}
      <MobileDrawer />

      {/* Profile Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleProfileMenuClose}
        PaperProps={{
          sx: {
            background: 'linear-gradient(135deg, rgba(26, 26, 26, 0.95) 0%, rgba(10, 10, 10, 0.95) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: 2,
            mt: 1,
          },
        }}
      >
        <MenuItem onClick={() => { navigate('/profile'); handleProfileMenuClose(); }}>
          <ListItemIcon>
            <AccountCircle />
          </ListItemIcon>
          Profile
        </MenuItem>
        <MenuItem onClick={() => { navigate('/settings'); handleProfileMenuClose(); }}>
          <ListItemIcon>
            <Settings />
          </ListItemIcon>
          Settings
        </MenuItem>
        <MenuItem onClick={() => { navigate('/security'); handleProfileMenuClose(); }}>
          <ListItemIcon>
            <Security />
          </ListItemIcon>
          Security
        </MenuItem>
        <MenuItem onClick={() => { navigate('/help'); handleProfileMenuClose(); }}>
          <ListItemIcon>
            <Help />
          </ListItemIcon>
          Help & Support
        </MenuItem>
        <Divider sx={{ my: 1, borderColor: 'rgba(255, 255, 255, 0.1)' }} />
        <MenuItem onClick={handleLogout}>
          <ListItemIcon>
            <Logout />
          </ListItemIcon>
          Logout
        </MenuItem>
      </Menu>

      {/* Toolbar spacer */}
      <Toolbar />
    </>
  );
};

export default Navbar; 
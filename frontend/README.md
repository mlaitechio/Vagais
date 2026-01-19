# AI-gency Frontend

A modern, responsive React application for the AI Agent Marketplace platform.

## Features

- 🎨 **Modern UI Design**: Inspired by financial dashboards and AI agency aesthetics
- 📱 **Responsive Design**: Works seamlessly on desktop, tablet, and mobile
- 🔐 **Authentication**: Secure login/register with JWT tokens
- 🏪 **Marketplace**: Browse and search AI agents with advanced filtering
- 📊 **Dashboard**: Comprehensive analytics and agent management
- ⚡ **Performance**: Optimized with React Query and Framer Motion
- 🎯 **TypeScript**: Full type safety throughout the application

## Tech Stack

- **React 18** with TypeScript
- **Material-UI (MUI)** for component library
- **Framer Motion** for animations
- **React Query** for data fetching and caching
- **React Router** for navigation
- **Axios** for API communication

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Backend API server running (see backend README)

### Installation

1. Clone the repository
2. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

3. Install dependencies:
   ```bash
   npm install
   ```

4. Create environment file:
   ```bash
   cp .env.example .env
   ```

5. Update the environment variables in `.env`:
   ```
   VITE_API_URL=http://localhost:8080/api/v1
   ```

6. Start the development server:
   ```bash
   npm run dev
   ```

7. Open [http://localhost:5173](http://localhost:5173) in your browser

## Project Structure

```
frontend/
├── src/
│   ├── components/          # Reusable UI components
│   ├── contexts/           # React contexts (Auth, Theme)
│   ├── pages/              # Page components
│   ├── services/           # API services
│   ├── styles/             # Global styles
│   ├── types/              # TypeScript type definitions
│   ├── App.tsx            # Main app component
│   └── main.tsx           # Entry point
├── public/                # Static assets
├── package.json           # Dependencies and scripts
└── README.md             # This file
```

## Key Components

### Pages
- **Dashboard**: Main dashboard with analytics and agent overview
- **Marketplace**: AI agent browsing with search and filters
- **Login/Register**: Authentication pages with modern design
- **Agent Detail**: Detailed view of individual agents

### Components
- **Navbar**: Responsive navigation with user menu
- **LoadingScreen**: Animated loading component
- **AgentCard**: Reusable agent display component

### Services
- **api.ts**: Centralized API service with axios
- **AuthContext**: Authentication state management

## Design System

### Color Palette
- **Primary**: `#98177E` (Purple)
- **Secondary**: `#00D4FF` (Cyan)
- **Success**: `#00FF88` (Green)
- **Warning**: `#FFB800` (Yellow)
- **Error**: `#FF4757` (Red)

### Typography
- **Font**: Inter (Google Fonts)
- **Weights**: 400, 500, 600, 700

### Animations
- **Framer Motion**: Smooth page transitions and micro-interactions
- **Gradients**: Modern gradient backgrounds and accents
- **Glassmorphism**: Translucent cards with backdrop blur

## API Integration

The frontend integrates with the Go backend API through:

- **RESTful endpoints** for all CRUD operations
- **JWT authentication** with automatic token refresh
- **Real-time updates** using React Query
- **Error handling** with user-friendly messages

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Environment Variables

Create a `.env` file in the frontend directory:

```env
# API Configuration
VITE_API_URL=http://localhost:8080/api/v1

# App Configuration
VITE_APP_NAME=AI-gency
VITE_APP_VERSION=1.0.0

# Feature Flags
VITE_ENABLE_ANALYTICS=true
VITE_ENABLE_NOTIFICATIONS=true
```

## Contributing

1. Follow the existing code style and patterns
2. Add TypeScript types for new components
3. Test on multiple screen sizes
4. Ensure accessibility standards are met

## Deployment

### Build for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.

### Environment Setup

Make sure to set the correct `VITE_API_URL` for your production environment.

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## License

This project is part of the AI-gency platform.

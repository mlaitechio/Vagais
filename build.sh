#!/bin/bash

# Build script for production deployment on single port
set -e

echo "🚀 Building application for single-port deployment..."

# Build frontend
echo "📦 Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Build backend
echo "🔨 Building backend..."
cd backend
go build -o ../bin/server .
cd ..

echo "✅ Build complete!"
echo ""
echo "📁 Files:"
echo "  - Backend binary: ./bin/server"
echo "  - Frontend static files: ./backend/dist/"
echo ""
echo "🌐 To run in production:"
echo "  cd backend"
echo "  PORT=8080 ENVIRONMENT=production ./server"
echo ""
echo "  The server will serve both API (on /api/*) and frontend (on /*) on the same port."

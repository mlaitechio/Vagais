#!/bin/bash

# Deployment script for production
set -e

echo "🚀 Starting production server..."

# Check if build exists
if [ ! -f "backend/dist/index.html" ]; then
    echo "❌ Frontend build not found. Run ./build.sh first."
    exit 1
fi

if [ ! -f "bin/server" ]; then
    echo "❌ Backend binary not found. Run ./build.sh first."
    exit 1
fi

# Set environment variables
export ENVIRONMENT=production
export PORT=${PORT:-8080}

# Navigate to backend directory (so relative paths work for static files)
cd backend

echo "🌐 Server starting on port $PORT"
echo "📍 Access application at: http://localhost:$PORT"
echo ""

# Run the server
../bin/server

#!/bin/bash

# CatiSip Startup Script
# This script helps you start the backend and frontend services

set -e

echo "====================================="
echo "   CatiSip - SIP System Startup     "
echo "====================================="
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
echo "Checking prerequisites..."

if ! command_exists go; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.21 or higher.${NC}"
    exit 1
fi

if ! command_exists node; then
    echo -e "${RED}Error: Node.js is not installed. Please install Node.js 16 or higher.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ All prerequisites met${NC}"
echo ""

# Start backend
echo "Starting backend server..."
cd backend

if [ ! -f "catiSip" ]; then
    echo "Building backend..."
    go build -o catiSip
fi

# Set default environment variables if not set
export SIP_HOST=${SIP_HOST:-localhost}
export SIP_PORT=${SIP_PORT:-5060}
export SIP_USERNAME=${SIP_USERNAME:-1000}
export SIP_PASSWORD=${SIP_PASSWORD:-1234}
export SIP_DOMAIN=${SIP_DOMAIN:-localhost}
export PORT=${PORT:-8080}

echo "Starting backend on port $PORT..."
./catiSip > ../backend.log 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > ../backend.pid

echo -e "${GREEN}✓ Backend started (PID: $BACKEND_PID)${NC}"
cd ..

# Wait for backend to be ready
echo "Waiting for backend to be ready..."
sleep 2
if curl -s http://localhost:8080/api/health > /dev/null; then
    echo -e "${GREEN}✓ Backend is healthy${NC}"
else
    echo -e "${YELLOW}⚠ Backend health check failed, but continuing...${NC}"
fi
echo ""

# Start frontend
echo "Starting frontend server..."
cd frontend

if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi

echo "Starting frontend development server..."
BROWSER=none npm start > ../frontend.log 2>&1 &
FRONTEND_PID=$!
echo $FRONTEND_PID > ../frontend.pid

echo -e "${GREEN}✓ Frontend started (PID: $FRONTEND_PID)${NC}"
cd ..

echo ""
echo "====================================="
echo -e "${GREEN}CatiSip is now running!${NC}"
echo "====================================="
echo ""
echo "Services:"
echo "  • Frontend: http://localhost:3000"
echo "  • Backend API: http://localhost:8080"
echo ""
echo "Logs:"
echo "  • Backend: tail -f backend.log"
echo "  • Frontend: tail -f frontend.log"
echo ""
echo "To stop all services, run:"
echo "  ./stop.sh"
echo ""
echo "Press Ctrl+C to view logs (services will continue running in background)"
echo ""

# Show logs
tail -f backend.log frontend.log

#!/bin/bash

# CatiSip Stop Script
# This script stops the backend and frontend services

echo "Stopping CatiSip services..."

# Stop backend
if [ -f backend.pid ]; then
    BACKEND_PID=$(cat backend.pid)
    if kill -0 $BACKEND_PID 2>/dev/null; then
        echo "Stopping backend (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        rm backend.pid
    else
        echo "Backend is not running"
        rm backend.pid
    fi
else
    echo "No backend PID file found"
fi

# Stop frontend
if [ -f frontend.pid ]; then
    FRONTEND_PID=$(cat frontend.pid)
    if kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "Stopping frontend (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
        rm frontend.pid
    else
        echo "Frontend is not running"
        rm frontend.pid
    fi
else
    echo "No frontend PID file found"
fi

# Clean up any remaining node processes
pkill -f "react-scripts start" 2>/dev/null || true

echo "CatiSip services stopped"

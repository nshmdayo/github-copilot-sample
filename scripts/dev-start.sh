#!/bin/bash

echo "🚀 Starting development environment..."

# Start with Docker Compose
docker-compose up -d postgres

echo "Waiting for PostgreSQL to start..."
sleep 10

# Start backend (background)
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!
cd ..

echo "Waiting for backend to start..."
sleep 5

# Start frontend
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo "✅ Development environment started!"
echo "📱 Frontend: http://localhost:3000"
echo "🔧 Backend: http://localhost:8080"
echo "📚 API Documentation: http://localhost:8080/swagger/index.html"
echo ""
echo "Press Ctrl+C to stop"

# Trap to terminate processes
trap "kill $BACKEND_PID $FRONTEND_PID; docker-compose down" EXIT

wait

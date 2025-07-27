#!/bin/bash

echo "🔨 Starting full build..."

# Backend build
echo "🔧 Building backend..."
cd backend
go build -o bin/server ./cmd/server
cd ..

# Frontend build
echo "📱 Building frontend..."
cd frontend
npm run build
cd ..

echo "✅ Full build completed!"

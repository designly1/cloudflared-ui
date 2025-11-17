#!/bin/bash

# Cloudflared GUI Setup Script
# This script sets up the development environment

set -e

echo "🚀 Setting up Cloudflared GUI..."
echo ""

# Check for required tools
echo "📋 Checking prerequisites..."

if ! command -v node &>/dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 20+ first."
    exit 1
fi

if ! command -v npm &>/dev/null; then
    echo "❌ npm is not installed. Please install npm first."
    exit 1
fi

if ! command -v go &>/dev/null; then
    echo "❌ Go is not installed. Please install Go 1.23+ first."
    exit 1
fi

echo "✅ All prerequisites found"
echo ""

# Install root dependencies
echo "📦 Installing root dependencies..."
npm install

# Install backend dependencies
echo "📦 Installing backend dependencies..."
cd apps/backend
go mod download
echo "✅ Backend dependencies installed"
cd ../..

# Install frontend dependencies
echo "📦 Installing frontend dependencies..."
cd apps/dashboard
npm install
echo "✅ Frontend dependencies installed"
cd ../..

# Install shared packages dependencies
echo "📦 Installing shared packages dependencies..."
cd packages/types
npm install
cd ../..

cd packages/ui
npm install
cd ../..

echo ""
echo "✅ Setup complete!"
echo ""
echo "📝 Next steps:"
echo ""
echo "  Development:"
echo "    npm run dev          # Run all services"
echo ""
echo "  Backend only:"
echo "    cd apps/backend && go run ./cmd/server"
echo ""
echo "  Frontend only:"
echo "    cd apps/dashboard && npm run dev"
echo ""
echo "  Build for production:"
echo "    npm run build"
echo ""
echo "  See README.md for deployment instructions."
echo ""

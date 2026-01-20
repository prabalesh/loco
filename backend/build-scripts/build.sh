#!/bin/bash

set -e  # Exit immediately if any command fails

echo "🔨 Building worker binary..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./builds/loco-worker cmd/worker/main.go

echo "🔨 Building server binary..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./builds/loco cmd/server/main.go

echo "✅ Build successful!"

echo "📤 Transferring binaries to server..."
scp ./builds/loco ./builds/loco-worker neo@37.27.6.129:/home/neo

echo "✅ Transfer complete!"
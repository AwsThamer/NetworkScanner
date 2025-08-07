#!/bin/bash

echo "Building Network Scanner for Linux and Windows..."

# Build for Linux
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o bin/network-scanner-linux main.go

# Build for Windows
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o bin/network-scanner-windows.exe main.go

echo "Build complete!"
echo "Linux binary: bin/network-scanner-linux"
echo "Windows binary: bin/network-scanner-windows.exe"

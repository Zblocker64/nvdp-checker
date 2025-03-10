#!/bin/bash

# Create output directory
mkdir -p bin

# Build for Linux (most common for servers)
GOOS=linux GOARCH=amd64 go build -o bin/nvdp-checker-linux-amd64 gpu.go
echo "Built for Linux (amd64)"

# Build for Linux ARM64 (for ARM-based servers)
GOOS=linux GOARCH=arm64 go build -o bin/nvdp-checker-linux-arm64 gpu.go
echo "Built for Linux (arm64)"

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o bin/nvdp-checker-darwin-amd64 gpu.go
echo "Built for macOS (amd64)"

# Build for macOS ARM (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o bin/nvdp-checker-darwin-arm64 gpu.go
echo "Built for macOS (arm64)"

echo "Build complete! Binaries are in the bin directory" 
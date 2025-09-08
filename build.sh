#!/bin/bash

# SWIFT Text Editor Build Script
# Streamlined Workflow, Increased Focus Typography

echo "🚀 Building SWIFT Text Editor..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -f swift swift.exe swift-editor swft

# Download dependencies
echo "📦 Downloading dependencies..."
go mod tidy

# Build for current platform
echo "🔨 Building SWIFT..."
go build -o swft .

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "📋 Installation options:"
    echo "   1. Local use: ./swft filename.txt"
    echo "   2. Global install: sudo cp swft /usr/local/bin/"
    echo "   3. User install: cp swft ~/bin/ (add ~/bin to PATH)"
    echo ""
    echo "🎯 Usage:"
    echo "   swft filename.txt    # Edit a file"
    echo "   swft                 # Start with welcome screen"
    echo ""
    echo "⌨️  Quick start:"
    echo "   - Press 'g' for help and commands"
    echo "   - All commands are intuitive and easy to remember!"
else
    echo "❌ Build failed!"
    exit 1
fi

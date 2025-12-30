#!/bin/bash
# Install TestU01 on macOS or Linux

set -e

echo "═══════════════════════════════════════════════════════════"
echo "  TestU01 Installation Script"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Detect OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
    echo "Detected: macOS"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
    echo "Detected: Linux"
else
    echo "Error: Unsupported OS: $OSTYPE"
    exit 1
fi

echo ""

# Install dependencies
echo "Installing dependencies..."
if [ "$OS" = "macos" ]; then
    if ! command -v brew &> /dev/null; then
        echo "Error: Homebrew not found. Please install Homebrew first."
        echo "Visit: https://brew.sh"
        exit 1
    fi
    brew install gsl
elif [ "$OS" = "linux" ]; then
    sudo apt-get update
    sudo apt-get install -y build-essential libgsl-dev wget unzip
fi

echo "✓ Dependencies installed"
echo ""

# Check if already installed
if [ -f /usr/local/lib/libtestu01.a ]; then
    echo "TestU01 appears to be already installed at /usr/local"
    read -p "Reinstall? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Installation cancelled"
        exit 0
    fi
fi

# Download TestU01
echo "Downloading TestU01..."
cd /tmp
rm -rf TestU01-1.2.3 TestU01.zip

if command -v curl &> /dev/null; then
    curl -L -o TestU01.zip https://simul.iro.umontreal.ca/testu01/TestU01.zip
else
    wget https://simul.iro.umontreal.ca/testu01/TestU01.zip
fi

# Verify download
if [ ! -f TestU01.zip ]; then
    echo "✗ Download failed"
    exit 1
fi

echo "✓ Downloaded"
echo ""

# Extract
echo "Extracting..."
unzip -q TestU01.zip
cd TestU01-1.2.3
echo "✓ Extracted"
echo ""

# Build and install
echo "Building TestU01 (this may take several minutes)..."
./configure --prefix=/usr/local
make -j$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 2)
echo "✓ Built"
echo ""

echo "Installing TestU01 (requires sudo)..."
sudo make install
echo "✓ Installed"
echo ""

# Verify installation
echo "Verifying installation..."
if [ -f /usr/local/lib/libtestu01.a ] && [ -f /usr/local/include/TestU01.h ]; then
    echo "✓ TestU01 successfully installed"
    echo ""
    echo "Library: /usr/local/lib/libtestu01.a"
    echo "Headers: /usr/local/include/*.h"
else
    echo "✗ Installation verification failed"
    echo "Checking what was installed:"
    ls -l /usr/local/lib/libtestu01* 2>/dev/null || echo "  No libraries found"
    ls -l /usr/local/include/TestU01.h 2>/dev/null || echo "  No TestU01.h found"
    exit 1
fi

# Cleanup
echo ""
echo "Cleaning up..."
cd /tmp
rm -rf TestU01-1.2.3 TestU01.zip
echo "✓ Cleanup complete"

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Installation Complete!"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Next steps:"
echo "  1. cd testu01"
echo "  2. make all              # Build test programs"
echo "  3. make smallcrush       # Run quick test (~2 min)"
echo ""

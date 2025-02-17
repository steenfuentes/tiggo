#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Installing tiggo..."

if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/doc/install"
    exit 1
fi
go build -o tiggo

INSTALL_DIR="/usr/local/bin"
if [[ ":$PATH:" != *":/usr/local/bin:"* ]]; then
    # If /usr/local/bin is not in PATH, use $HOME/.local/bin
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"

    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo "export PATH=\$PATH:$HOME/.local/bin" >> "$HOME/.bashrc"
        echo "export PATH=\$PATH:$HOME/.local/bin" >> "$HOME/.zshrc" 2>/dev/null || true
    fi
fi

echo "Installing to $INSTALL_DIR..."
echo "Requesting sudo permission to copy tiggo executable to $INSTALL_DIR..."
sudo cp tiggo "$INSTALL_DIR/tiggo" 2>/dev/null || cp tiggo "$INSTALL_DIR/tiggo"

# Clean up after ourselves
rm tiggo

echo "tiggo is installed in $INSTALL_DIR"
echo -e "${GREEN}Installation complete!${NC}"
echo "You may need to restart your terminal or run 'source ~/.bashrc (or ~/.zshrc)' to use tiggo"
echo "Use tiggo by running 'tiggo' in your terminal. Use 'tiggo --help' for more information."
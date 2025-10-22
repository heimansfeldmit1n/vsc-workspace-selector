#!/bin/bash

# Install latest Go if not present (assumes Linux amd64)
if ! command -v go &> /dev/null; then
    echo "Installing latest Go..."
    wget -q "https://dl.google.com/go/$(curl -s https://go.dev/VERSION?m=text).linux-amd64.tar.gz"
    sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz
    rm go*.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo "Go installed. Restart terminal or source ~/.bashrc for PATH changes."
else
    echo "Go is already installed (version: $(go version))."
fi

# Build contents in src directory and copy to global bin
if [ -d "src" ]; then
    cd src
    go mod tidy  # Ensure dependencies
    BINARY_NAME="vsc"  # Name binary after src dir (or customize here)
    go build -o "$BINARY_NAME" .   # Build with explicit name for easy handling
    if [ -f "$BINARY_NAME" ]; then
        sudo mv "$BINARY_NAME" /usr/local/bin/
        echo "Build complete. Binary '$BINARY_NAME' copied to /usr/local/bin/ for global access."
        echo "You can now run it from anywhere: $BINARY_NAME"
    else
        echo "Error: Build failed - binary not found."
        exit 1
    fi
    cd ..  # Return to original dir
else
    echo "Error: src directory not found."
    exit 1
fi
# Copy config to original user's home directory (even when run as sudo)
if [ -n "$SUDO_USER" ] && [ "$SUDO_USER" != "$USER" ]; then
    USER_HOME=$(getent passwd "$SUDO_USER" | cut -d: -f6)
    echo "Copying config file to original user's home directory: $USER_HOME"
    cp vsc-selector.config.toml "$USER_HOME/"
else
    echo "Copying config file to home directory"
    cp vsc-selector.config.toml ~/
fi
echo "Setup done"
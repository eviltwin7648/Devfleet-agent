#!/usr/bin/env bash
set -e

REPO="eviltwin7648/Devfleet-agent"
VERSION="agent"
TOKEN=$1

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# normalize arch
if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; fi
if [ "$ARCH" = "aarch64" ]; then ARCH="arm64"; fi

BINARY="devfleet-agent-$OS-$ARCH"
URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY"

echo "Downloading $BINARY..."

curl -L "$URL" -o devfleet-agent
chmod +x devfleet-agent

echo "Installing..."
sudo mv devfleet-agent /usr/local/bin/

echo "Starting agent..."
devfleet-agent start --token="$TOKEN"

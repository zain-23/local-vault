#!/bin/sh
set -e

# LocalVault installer
# Usage: curl -fsSL https://raw.githubusercontent.com/zain-23/local-vault/main/install.sh | sh

REPO="zain-23/local-vault"
BINARY="lv"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  *)      echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
  x86_64)          ARCH="amd64" ;;
  arm64|aarch64)   ARCH="arm64" ;;
  *)               echo "❌ Unsupported arch: $ARCH"; exit 1 ;;
esac

# Get latest version
echo "🔍 Fetching latest version..."
VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' \
  | cut -d'"' -f4)

if [ -z "$VERSION" ]; then
  echo "❌ Could not fetch latest version"
  exit 1
fi

echo "📦 Installing LocalVault $VERSION ($OS/$ARCH)..."

# Download and extract
URL="https://github.com/$REPO/releases/download/$VERSION/lv_${OS}_${ARCH}.tar.gz"
curl -fsSL "$URL" | tar -xz -C /tmp lv

# Install binary
sudo mv /tmp/lv "$INSTALL_DIR/$BINARY"
sudo chmod +x "$INSTALL_DIR/$BINARY"

echo ""
echo "✅ LocalVault $VERSION installed successfully"
echo ""
echo "Get started:"
echo "  cd your-project"
echo "  lv init"
echo "  lv unlock"
echo "  lv add DATABASE_URL=postgres://..."
echo "  lv inject -- npm run dev"
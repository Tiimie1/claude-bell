#!/bin/bash
set -euo pipefail

REPO="Tiimie1/claude-bell"
INSTALL_DIR="/usr/local/bin"

# macOS only
if [ "$(uname -s)" != "Darwin" ]; then
    echo "Error: claude-bell currently only supports macOS."
    exit 1
fi

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    arm64)   ARCH="arm64" ;;
    *)
        echo "Error: unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Get latest release tag
LATEST="$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)"
if [ -z "$LATEST" ]; then
    echo "Error: could not determine latest release."
    exit 1
fi

TARBALL="claude-bell_darwin_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARBALL}"

echo "Downloading claude-bell ${LATEST} for darwin/${ARCH}..."
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

curl -sSL "$URL" -o "${TMPDIR}/${TARBALL}"
tar -xzf "${TMPDIR}/${TARBALL}" -C "$TMPDIR"

echo "Installing to ${INSTALL_DIR}/claude-bell..."
sudo install -m 755 "${TMPDIR}/claude-bell" "${INSTALL_DIR}/claude-bell"

echo "Done! Run 'claude-bell setup' to get started."

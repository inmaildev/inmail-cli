#!/usr/bin/env bash
set -euo pipefail

REPO="inmaildev/inmail-cli"
BIN_NAME="inmail"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

detect_os() {
  case "$(uname -s)" in
    Darwin) echo "darwin" ;;
    Linux)  echo "linux"  ;;
    *)      echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
  esac
}

OS=$(detect_os)
ARCH=$(detect_arch)

echo "→ Detecting latest release..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "✗ Could not determine latest release. Check https://github.com/${REPO}/releases" >&2
  exit 1
fi

VERSION="${LATEST#v}"
ASSET="${BIN_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${ASSET}"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "→ Downloading ${BIN_NAME} ${LATEST} (${OS}/${ARCH})..."
curl -fsSL "$URL" -o "${TMP}/${ASSET}"

echo "→ Extracting..."
tar -xzf "${TMP}/${ASSET}" -C "$TMP"

echo "→ Installing to ${INSTALL_DIR}/${BIN_NAME}..."
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
else
  sudo mv "${TMP}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
fi

chmod +x "${INSTALL_DIR}/${BIN_NAME}"

echo "✓ ${BIN_NAME} ${LATEST} installed successfully"
echo ""
echo "  Get started:"
echo "    inmail configure"
echo "    inmail --help"

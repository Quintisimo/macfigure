#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="macfigure"
INSTALL_DIR="/usr/local/bin"

ARCH=$(uname -m)
case "${ARCH}" in
  arm64|aarch64) PKL_ARCH="aarch64" ;;
  x86_64|amd64) PKL_ARCH="amd64" ;;
  *) echo "Unsupported architecture: ${ARCH}" >&2; exit 1 ;;
esac

TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT

if ! command -v brew >/dev/null 2>&1; then
  echo "Installing Homebrew..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

  if [[ "${PKL_ARCH}" == "aarch64" ]]; then
    BREW_PREFIX="/opt/homebrew"
  else
    BREW_PREFIX="/usr/local"
  fi
  eval "$(${BREW_PREFIX}/bin/brew shellenv)"
fi

gh_release_url() {
  echo "https://github.com/$1/releases/latest/download/$2"
}

echo "Downloading pkl..."
curl -fsSL "$(gh_release_url "apple/pkl" "pkl-macos-${PKL_ARCH}")" -o "${TMP_DIR}/pkl"

echo "Installing pkl to ${INSTALL_DIR}..."
sudo mv "${TMP_DIR}/pkl" "${INSTALL_DIR}/pkl"
sudo chmod +x "${INSTALL_DIR}/pkl"

echo "Downloading ${PROJECT_NAME}..."
curl -fsSL "$(gh_release_url "quintisimo/macfigure" "$PROJECT_NAME")" -o "${TMP_DIR}/${PROJECT_NAME}"

echo "Installing ${PROJECT_NAME} to ${INSTALL_DIR}..."
sudo mv "${TMP_DIR}/${PROJECT_NAME}" "${INSTALL_DIR}/${PROJECT_NAME}"
sudo chmod +x "${INSTALL_DIR}/${PROJECT_NAME}"

echo "${PROJECT_NAME} installed successfully."

#!/usr/bin/env bash

set -euo pipefail
# set -x
# exec 2>&1

SCRIPT_DIRECTORY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIRECTORY/.." && pwd)"

set -o allexport
source "$PROJECT_ROOT/.env"
set +o allexport

ARCH="$(uname -m)"
RELEASE_URL="https://github.com/firecracker-microvm/firecracker/releases"
LATEST=$(basename "$(curl -fsSLI -o /dev/null -w '%{url_effective}' "${RELEASE_URL}/latest")")

curl -L "${RELEASE_URL}/download/${LATEST}/firecracker-${LATEST}-${ARCH}.tgz" | tar -xz -C "$PROJECT_ROOT"

mkdir -p "${FIRECRACKER_BINARY_PATH}"

mv "$PROJECT_ROOT/release-${LATEST}-${ARCH}/firecracker-${LATEST}-${ARCH}" "${FIRECRACKER_BINARY_PATH}/firecracker"

sudo rm -rf "${PROJECT_ROOT}/release-${LATEST}-${ARCH}"

echo "Firecracker ${LATEST} installed at ${FIRECRACKER_BINARY_PATH}/firecracker"
#!/usr/bin/env bash
set -e
#there has to be a better way to do this

SCRIPT_DIRECTORY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIRECTORY/.." && pwd)"

set -o allexport
source "$PROJECT_ROOT/.env"
set +o allexport

rm -f "${KERNEL_PATH}"/vmlinux-*
mkdir -p "${KERNEL_PATH}"
ARCH="$(uname -m)"
S3="https://s3.amazonaws.com/spec.ccfc.min"

# Fetch latest CI artifact build. `sort` + `tail` returns artifacts with latest YYYYMMDD.
CI_ARTIFACTS_PREFIX=$(curl -fsSL "$S3?list-type=2&prefix=firecracker-ci/&delimiter=/" \
    | grep -oP "(?<=<Prefix>)firecracker-ci/[0-9]{8}-[^/]+/(?=</Prefix>)" \
    | sort \
    | tail -1)

latest_kernel_key=$(curl -fsSL "$S3?list-type=2&prefix=${CI_ARTIFACTS_PREFIX}${ARCH}/vmlinux-" \
      | grep -oP "(?<=<Key>)(${CI_ARTIFACTS_PREFIX}${ARCH}/vmlinux-[0-9]+\.[0-9]+\.[0-9]{1,3})(?=</Key>)" \
      | sort -V \
      | tail -1)

# Download a linux kernel binary
wget -O "${KERNEL_PATH}"/vmlinux "$S3/${latest_kernel_key}"

echo "VMLINUX Version : ${latest_kernel_key}"
echo
echo "The following files were downloaded and set up:"
KERNEL=$(ls "${KERNEL_PATH}"/vmlinux | tail -1)
[ -f $KERNEL ] && echo "Kernel: $KERNEL" || echo "ERROR: Kernel $KERNEL does not exist"

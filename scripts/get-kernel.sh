#!/bin/bash
# Fetch a Firecracker CI guest kernel (vmlinux) for the host architecture.
# Kernel only: no rootfs, no SSH keys, nothing built. Run just in case the
# working kernel ever needs replacing; the known-good vmlinux-6.1.186 stays
# untouched (downloads land as versioned files, never overwriting it).
#
# Usage:
#   ./scripts/get-kernel.sh
#   KERNEL_LINE=6.1 ./scripts/get-kernel.sh          # default: 6.x LTS line
#   KERNEL_LINE=6.18 ./scripts/get-kernel.sh         # explicit line
#   OUT_DIR=/tmp/kernels ./scripts/get-kernel.sh     # custom destination
#
# Flow (same as the upstream getting-started guide): list the CI bucket for
# the newest dated build prefix -> list vmlinux images under it -> pick the
# highest patch on the requested line -> resume-safe download -> verify.

# -e: abort on first failure (a partial kernel must never look complete).
# -u: abort on unset variables. -o pipefail: grep finding nothing fails loudly.
set -euo pipefail

# Host arch selects the bucket path (x86_64 or aarch64).
ARCH="$(uname -m)"
# Public CI artifact bucket (same one the Firecracker docs use).
S3="https://s3.amazonaws.com/spec.ccfc.min"
# Kernel line to fetch: 6.1 is the LTS line proven with Firecracker v1.17 in
# this repo (virtio attaches, root mounts). Override per usage above.
KERNEL_LINE="${KERNEL_LINE:-6.1}"
# Destination dir; defaults to the existing kernel dir next to this repo.
OUT_DIR="${OUT_DIR:-$HOME/Desktop/linux}"

# Newest CI build: `sort` + `tail -1` over dated prefixes (YYYYMMDD-...).
CI_ARTIFACTS_PREFIX=$(curl -fsSL "$S3?list-type=2&prefix=firecracker-ci/&delimiter=/" \
    | grep -oP "(?<=<Prefix>)firecracker-ci/[0-9]{8}-[^/]+/(?=</Prefix>)" \
    | sort \
    | tail -1)
echo "CI build: $CI_ARTIFACTS_PREFIX"

# Highest patch release on the requested line, e.g. vmlinux-6.1.186.
latest_kernel_key=$(curl -fsSL "$S3?list-type=2&prefix=${CI_ARTIFACTS_PREFIX}${ARCH}/vmlinux-" \
    | grep -oP "(?<=<Key>)(${CI_ARTIFACTS_PREFIX}${ARCH}/vmlinux-${KERNEL_LINE}\.[0-9]{1,3})(?=</Key>)" \
    | sort -V \
    | tail -1)

# Fail loudly if the line has no build (e.g. typo in KERNEL_LINE) instead of
# downloading the bucket root or an empty URL.
if [ -z "$latest_kernel_key" ]; then
    echo "ERROR: no vmlinux-${KERNEL_LINE}.x found under $CI_ARTIFACTS_PREFIX$ARCH/" >&2
    exit 1
fi
echo "Kernel: $latest_kernel_key"

# Versioned filename (vmlinux-6.1.186): never clobbers the working kernel.
OUT="$OUT_DIR/$(basename "$latest_kernel_key")"
mkdir -p "$OUT_DIR"

# Resume-safe download (-C -): re-runs continue where a dropped connection
# stopped instead of restarting ~45MB (this link is slow).
curl -fL -C - -o "$OUT" "$S3/$latest_kernel_key"

# Verify: non-empty ELF kernel image, not an error page or truncated file.
[ -s "$OUT" ] || { echo "ERROR: downloaded file is empty: $OUT" >&2; exit 1; }
file "$OUT" | grep -q ELF || { echo "ERROR: not an ELF kernel image: $OUT" >&2; exit 1; }
ls -lh "$OUT"
echo "Saved kernel to $OUT"

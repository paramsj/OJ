#!/bin/bash
# Build the Firecracker guest rootfs image.
# Run from anywhere: ./scripts/build-rootfs.sh  (needs sudo for mount/umount,
# and /tmp/trials-agent must already exist — see step 0 in the header below).
#
# 0. Prerequisite (run once, rebuild only if trials/ source changed):
#      cd ~/Desktop/trials
#      CGO_ENABLED=0 go build -tags netgo -ldflags '-extldflags "-static"' -o /tmp/trials-agent .
#    This produces the static trials :1323 agent binary that the guest will run.
# 1. This script: dd zeroed file -> mkfs.ext4 -> mount it -> run an Alpine
#    container with the image bind-mounted at /my-rootfs so setup-guest.sh
#    (fed via stdin, like the reference flow) installs and configures the
#    whole system straight into the image -> unmount -> fsck verify.
# 2. Boot it (gg server, unmodified) with kernel ~/Desktop/linux/vmlinux-6.1.186:
#      curl -X POST localhost:1323/vm/start \
#        -d '{"root_image_path":"/home/jaip/Desktop/trials/images/rootfs-base.ext4","kernel_path":"/home/jaip/Desktop/linux/vmlinux-6.1.186"}'

# -e: abort on first failing command (a half-built image must never be used).
# -x: echo each command (build log shows exactly what ran).
set -xe

# Resolve the repo root from this script's own location, so the script works
# no matter which directory you invoke it from (the reference relied on $PWD).
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# Final artifact: the bootable ext4 image, kept inside the repo.
IMG="$ROOT/images/rootfs-base.ext4"

# Create a 512MB zeroed file: the raw block device the ext4 will live on.
# (~300MB expected used after slimming; bump count to 1024 for more margin.)
dd if=/dev/zero of="$IMG" bs=1M count=512
# Format it as ext4 (the filesystem the guest kernel mounts as root).
mkfs.ext4 "$IMG"
# Scratch mount point where the host will access the image contents.
mkdir -p /tmp/my-rootfs
# Mount requires root: attach the image file as a filesystem at /tmp/my-rootfs.
sudo mount "$IMG" /tmp/my-rootfs

# Launch a throwaway Alpine container that writes the guest OS into the mount:
# -i keeps stdin open so setup-guest.sh can be piped in as the shell script.
# --rm deletes the container afterwards (the artifact is the image, not it).
# -v /tmp/my-rootfs:/my-rootfs exposes the mounted image inside at /my-rootfs.
# -v /tmp/trials-agent:/usr/local/bin/agent:ro overlays your prebuilt static
#   binary onto that container path, so the tar copy-loop below carries it
#   into the image as /usr/local/bin/agent (same trick as the reference).
# -v .../openrc-agent:/etc/init.d/agent:ro same trick for the OpenRC service
#   file: must already sit at /etc/init.d/agent before `rc-update add agent`.
# alpine:3.20 pinned (not `latest`) so gcc version and layout are reproducible.
# `sh <file` feeds the guest setup script to the container's shell via stdin.
docker run -i --rm \
    -v /tmp/my-rootfs:/my-rootfs \
    -v /tmp/trials-agent:/usr/local/bin/agent:ro \
    -v "$ROOT/scripts/openrc-agent:/etc/init.d/agent:ro" \
    alpine:3.20 sh <"$ROOT/scripts/setup-guest.sh"

# Detach the image (needs root); must happen after the container exited,
# otherwise the mount is busy and umount fails.
sudo umount /tmp/my-rootfs
# Read-only consistency check: proves the image is a valid, mountable ext4
# before anyone tries to boot it.
e2fsck -fn "$IMG"

# rootfs available under `images/rootfs-base.ext4`

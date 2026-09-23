#!/bin/sh
# Guest-side setup. Runs INSIDE the Alpine container (paths like /bin, /etc
# mean the container's Alpine root, NOT your laptop — running this on the host
# would pack Ubuntu system dirs into the image and destroy it).
# It configures Alpine + toolchain + agent, then copies the whole system into
# the mounted ext4 image at /my-rootfs (see the tar loop near the end).

# -e: abort on first failure (never copy a half-configured system).
# -x: echo each command into the build log.
set -xe

# Install the init system that will start the agent at every boot.
apk add --no-cache openrc
# Serial-console login tools (agetty) and misc base utilities.
apk add --no-cache util-linux
# The C/C++ toolchain the agent shells out to (g++ compiles submissions).
# Deliberately NO python/go: alpine+gcc is all this C++-only agent needs
# (python2 doesn't even exist in Alpine 3.20 repos and would abort the build).
apk add --no-cache gcc g++ libc-dev

# Create the serial-console login service: a microVM has no screen/keyboard,
# its only console is the serial port (ttyS0, per kernel console=ttyS0), so a
# getty there is the only way to see and reach the guest console.
ln -s agetty /etc/init.d/agetty.ttyS0
# Allow root logins on that serial console (appended; duplicates harmless).
echo ttyS0 >>/etc/securetty
# Start the serial getty automatically in the default runlevel.
rc-update add agetty.ttyS0 default

# DNS: the container's resolv.conf holds Docker's embedded DNS (127.0.0.11),
# which must NOT leak into the image, so overwrite (not append) with public
# resolvers BEFORE the copy loop below carries /etc into the image.
printf 'nameserver 1.1.1.1\nnameserver 8.8.8.8\n' >/etc/resolv.conf
# Guest hostname, read by the hostname OpenRC service at boot.
printf 'fc-guest\n' >/etc/hostname

# Belt-and-braces serial console: if Alpine's default inittab already enables
# a getty on ttyS0 the grep matches and nothing happens; otherwise append it.
# Without this line a silent guest is un-debuggable (no console output).
grep -q '^ttyS0::respawn' /etc/inittab || echo 'ttyS0::respawn:/sbin/getty -L 115200 ttyS0 vt100' >>/etc/inittab

# Mount the kernel virtual filesystems at boot: OpenRC mounts devtmpfs/proc/
# sysfs onto /dev /proc /sys, so those runlevel entries must exist.
rc-update add devfs boot
rc-update add procfs boot
rc-update add sysfs boot

# Start the agent in the default runlevel (normal services live in `default`;
# `boot` is for early system setup). The service file itself arrived via the
# -v bind-mount, so it is already at /etc/init.d/agent when this runs.
rc-update add agent default

# Transplant the configured system into the image: for each system dir, tar
# the container path to stdout and untar it under /my-rootfs. Only these five
# dirs hold real files; dev/proc/sys/run are kernel/runtime pseudo-filesystems
# (copying the container's /proc would dump garbage process entries).
for d in bin etc lib root sbin usr; do tar c "/$d" | tar x -C /my-rootfs; done
# Scaffold the runtime dirs as EMPTY mount points/scratch: OpenRC mounts
# devtmpfs/proc/sysfs onto dev/proc/sys at boot, and run/var/tmp hold PID
# files, spool and temp data (incl. the agent's /tmp/<id> submission files).
for dir in dev proc run sys var tmp; do mkdir /my-rootfs/${dir}; done

# Standard /tmp semantics, load-bearing for the agent: anyone may create
# submission files, but only the owner may delete them (sticky bit).
chmod 1777 /my-rootfs/tmp

#!/usr/bin/env bash
set -e

# im assuming this much is enuf for gcc + jvm + python + openrc and more stuff
# if you want more just increase the count

# snapshotting in the future

SCRIPT_DIRECTORY="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIRECTORY/.." && pwd)"

set -o allexport
source "$PROJECT_ROOT/.env"
set +o allexport

mkdir -p "${ROOTFS_PATH}"
 
ROOTFS="${ROOTFS_PATH}/rootfs.ext4"
dd if=/dev/zero of="$ROOTFS" bs=1M count=1000
mkfs.ext4 "$ROOTFS"

# mount the fs to the machines tmp, anything i put in there is with me
mkdir -p /tmp/my-rootfs
sudo mount "$ROOTFS" /tmp/my-rootfs

echo "Entering the docker container now"

sudo docker run --rm \
-v "${PROJECT_ROOT}/build/agent:/usr/bin/agent" \
-v "${SCRIPT_DIRECTORY}/openrc-setup.sh:/etc/init.d/agent" \
-v /tmp/my-rootfs:/my-rootfs alpine:latest sh -c '

    echo "installing packages"
    apk add --no-cache openrc util-linux gcc libc-dev g++ agetty
    echo "all packages installed"

    echo "Permission Setup for Binaries"
    chmod +x /usr/bin/agent /etc/init.d/agent
    
    # Set up a login terminal on the serial console (ttyS0) -> uncomment to debug it hijacks my terminals:
    # ln -s agetty /etc/init.d/agetty.ttyS0
    # echo ttyS0 > /etc/securetty
    # rc-update add agetty.ttyS0 default

    echo "Adding Services top OpenRC"
    # echo "root:root" | chpasswd

    # Make sure special file systems are mounted on boot:
    rc-update add devfs boot
    rc-update add procfs boot
    rc-update add sysfs boot
    rc-update add agent boot

    echo "Copyfiles into the rootfs"
    # Then, copy the newly configured system to the rootfs image:
    for d in bin etc lib root sbin usr; do tar c "/$d" | tar x -C /my-rootfs; done

    # The above command may trigger the following message:
    # tar: Removing leading "/" from member names
    # However, this is just a warning, so you should be able to
    # proceed with the setup process.

    for dir in dev proc run sys var tmp; do mkdir /my-rootfs/${dir}; done

    chmod 1777 /my-rootfs/tmp
    echo "All done"
    exit
'
sudo umount /tmp/my-rootfs
rm -rf /tmp/my-rootfs

echo "now you can hope and pray that the rootfs was created"

# okay, in case the script fails, YOU NEED to UMOUNT /tmp/root-fs (add)
#sudo umount /tmp/root-fs -> run this in case it says rootfs is already mounted

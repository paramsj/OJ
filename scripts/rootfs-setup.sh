#!/usr/bin/env bash
set -e

# im assuming this much is enuf for gcc + jvm + python + openrc and more stuff
# if you want more just increase the count

# snapshotting in the future
dd if=/dev/zero of=../rootfs.ext4 bs=1M count=1000
mkfs.ext4 ../rootfs.ext4

# mount the fs to the machines tmp, anything i put in there is with me
mkdir -p /tmp/my-rootfs
sudo mount ../rootfs.ext4 /tmp/my-rootfs

echo "im here"

docker run --rm \
-v "$(pwd)/../build/agent:/usr/bin/agent" \
-v "$(pwd)/openrc-setup.sh:/etc/init.d/agent" \
-v /tmp/my-rootfs:/my-rootfs alpine:latest sh -c '

    apk add openrc
    apk add util-linux
    apk add --no-cache gcc libc-dev
    apk add --no-cache g++

    chmod +x /usr/bin/agent /etc/init.d/agent

    # to do here is to add gcc and other stuff you may require
    
    # Set up a login terminal on the serial console (ttyS0):
    ln -s agetty /etc/init.d/agetty.ttyS0
    echo ttyS0 > /etc/securetty
    rc-update add agetty.ttyS0 default

    # Make sure special file systems are mounted on boot:
    rc-update add devfs boot
    rc-update add procfs boot
    rc-update add sysfs boot
    rc-update add agent boot

    # Then, copy the newly configured system to the rootfs image:
    for d in bin etc lib root sbin usr; do tar c "/$d" | tar x -C /my-rootfs; done

    # The above command may trigger the following message:
    # tar: Removing leading "/" from member names
    # However, this is just a warning, so you should be able to
    # proceed with the setup process.

    for dir in dev proc run sys var tmp; do mkdir /my-rootfs/${dir}; done

    chmod 1777 /my-rootfs/tmp

    # All done, exit docker shell.
    exit
'
sudo umount /tmp/my-rootfs
rm -rf /tmp/my-rootfs

echo "now you can hope and pray that the rootfs was created"

# okay, in case the script fails, YOU NEED to UMOUNT /tmp/root-fs (add)
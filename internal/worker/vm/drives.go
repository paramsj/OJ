package vm

import (
	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

func DriveConfig(VMID string) models.Drive {
	return models.Drive{
		IsReadOnly:   firecracker.Bool(false),
		IsRootDevice: firecracker.Bool(true),
		// VeryLargeUUID-rootfs.ext4
		PathOnHost: firecracker.String("/tmp/rootfs-" + VMID + ".ext4"),
		DriveID:    firecracker.String("rootfs"),
	}
}

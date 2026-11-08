package vm

import (
	"os"
	"path/filepath"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

func GetFirecrackerConfig(VMID string) (firecracker.Config, error) {
	sockPath, _ := getSocketPath(VMID)
	return firecracker.Config{
		SocketPath:      sockPath,
		KernelImagePath: "/home/jaip/Desktop/agentV2/vmlinux-6.18.56",
		Drives: []models.Drive{
			DriveConfig(VMID),
		},
		NetworkInterfaces: NetworkConfig(),
		MachineCfg:        MachineConfig(),
	}, nil
}

func getSocketPath(VMID string) (string, error) {
	vmDir := filepath.Join("/tmp", "firecracker", VMID)
	if err := os.MkdirAll(vmDir, 0700); err != nil {
		return "", err
	}
	socketPath := filepath.Join(vmDir, "firecracker.sock")
	return socketPath, nil
}

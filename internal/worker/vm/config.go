package vm

import (
	"os"
	"path/filepath"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

func GetFirecrackerConfig(VMID string) (firecracker.Config, error) {
	sockPath, _ := getSocketPath(VMID)
	pwd, err := os.Getwd()
	if err != nil {
		return firecracker.Config{}, err
	}
	_ = os.MkdirAll(filepath.Join(pwd, "logs", VMID), 0755)
	return firecracker.Config{
		SocketPath: sockPath,
		// KERNEL_PATH/vmlinux
		KernelImagePath: filepath.Join(os.Getenv("KERNEL_PATH"), "vmlinux"),
		Drives: []models.Drive{
			DriveConfig(VMID),
		},
		NetworkInterfaces: NetworkConfig(),
		MachineCfg:        MachineConfig(),
		LogPath:           filepath.Join(pwd, "logs", VMID, "firecracker.log"),
		LogLevel:          "Debug",
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

package vm

import (
	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

func MachineConfig() models.MachineConfiguration {
	return models.MachineConfiguration{
		MemSizeMib: firecracker.Int64(256),
		VcpuCount:  firecracker.Int64(1),
	}
}

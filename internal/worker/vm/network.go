package vm

import (
	"github.com/firecracker-microvm/firecracker-go-sdk"
)

func NetworkConfig() firecracker.NetworkInterfaces {
	return firecracker.NetworkInterfaces{
		{
			CNIConfiguration: &firecracker.CNIConfiguration{
				BinPath:     []string{"/opt/cni/bin"},
				ConfDir:     "/etc/cni/conf.d",
				NetworkName: "fcnet",
				IfName:      "veth0",
			},
		},
	}
}

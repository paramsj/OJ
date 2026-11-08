package controllers

import (
	"agentV2/internal/worker/vm"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"uuid"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/labstack/echo/v5"
)

// type SubmissionRequest struct {
// 	Language    string  `json:"language"`
// 	Code        string  `json:"code"`
// 	Input       string  `json:"input"`
// 	TimeAllowed float64 `json:"timeAllowed"`
// }

// type SubmissionResponse struct {
// 	Verdict   string  `json:"verdict"`
// 	Stdout    string  `json:"stdout"`
// 	Stderr    string  `json:"stderr"`
// 	TimeTaken float64 `json:"timetaken"`
// 	MemUsage  int64   `json:"memusage"`
// }

func Submit(c *echo.Context) error {
	// var req SubmissionRequest
	// if err := c.Bind(&req); err != nil {
	// 	return c.String(http.StatusBadRequest, "this was a bad request, please try again")
	// }

	VMID := uuid.NewV4()

	_, err := copyFS("/home/jaip/Desktop/agentV2/rootfs.ext4", "/tmp/rootfs-"+VMID.String()+".ext4")

	if err != nil {
		return c.String(http.StatusNotFound, "We are here"+err.Error())
	}

	cfg, err := vm.GetFirecrackerConfig(VMID.String())
	if err != nil {
		return c.String(http.StatusNotFound, "After config"+err.Error())
	}
	ctx := context.Background()
	m, err := firecracker.NewMachine(ctx, cfg)

	if err != nil {
		panic(fmt.Errorf("failed to create new machine: %v", err))
	}
	defer os.Remove(cfg.SocketPath)
	if err := m.Start(ctx); err != nil {
		panic(fmt.Errorf("failed to initialize machine: %v", err))
	}
	defer m.StopVMM()
	ipAddr := m.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP
	return c.String(http.StatusOK, ipAddr.String())
}

func copyFS(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err

}

package controllers

import (
	"agentV2/internal/worker/vm"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"uuid"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/labstack/echo/v5"
)

type SubmissionRequest struct {
	Language    string  `json:"language"`
	Code        string  `json:"code"`
	Input       string  `json:"input"`
	TimeAllowed float64 `json:"timeAllowed"`
}

type SubmissionResponse struct {
	Verdict   string  `json:"verdict"`
	Stdout    string  `json:"stdout"`
	Stderr    string  `json:"stderr"`
	TimeTaken float64 `json:"timetaken"`
	MemUsage  int64   `json:"memusage"`
}

func Start(c *echo.Context) error {
	// var req SubmissionRequest
	// if err := c.Bind(&req); err != nil {
	// 	return c.String(http.StatusBadRequest, "this was a bad request, please try again")
	// }

	VMID := uuid.NewV4()
	ROOTFS := filepath.Join(os.Getenv("ROOTFS_PATH"), "rootfs.ext4")
	_, err := copyFS(ROOTFS, "/tmp/rootfs-"+VMID.String()+".ext4")

	if err != nil {
		return c.String(http.StatusNotFound, "We are here"+err.Error())
	}

	cfg, err := vm.GetFirecrackerConfig(VMID.String())
	if err != nil {
		return c.String(http.StatusNotFound, "After config"+err.Error())
	}

	pwd, err := os.Getwd()

	stdoutPath := filepath.Join(pwd, "logs", VMID.String(), "stdout.log")

	stdout, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return c.String(http.StatusNotFound, "After config"+err.Error())
	}

	// stderr will be directed to this file
	stderrPath := filepath.Join(pwd, "logs", VMID.String(), "stderr.log")
	stderr, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to create stderr file: %v", err))
	}

	ctx := context.Background()

	cmd := firecracker.VMCommandBuilder{}.
		WithBin("firecracker").
		WithSocketPath(cfg.SocketPath).
		WithStdout(stdout).
		WithStderr(stderr).
		Build(ctx)

	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))

	if err != nil {
		panic(fmt.Errorf("failed to create new machine: %v", err))
	}
	if err := m.Start(ctx); err != nil {
		panic(fmt.Errorf("failed to initialize machine: %v", err))
	}

	defer func() {
		// m.StopVMM() //is this goes i go with it
		// os.Remove(cfg.SocketPath)
		// os.Remove(filepath.Join("/tmp/rootfs-" + VMID.String() + ".ext4"))
	}()
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

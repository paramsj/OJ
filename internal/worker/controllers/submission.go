package controllers

import (
	"agentV2/internal/worker/vm"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

type Controller struct {
	vms        *vm.Manager
	WarmVMPool chan *vm.VM
}

func NewController(m *vm.Manager) *Controller {
	return &Controller{
		vms:        m,
		WarmVMPool: make(chan *vm.VM, 20),
	}
}

func (h *Controller) Start(c *echo.Context) error {
	v, err := h.create()
	if err != nil {
		return c.String(http.StatusInternalServerError, "The VM was not created :(")
	}
	return c.String(http.StatusOK, v.IP.String())
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

// map[int] -> VM

func (h *Controller) Delete(c *echo.Context) error {
	VMID, err := echo.PathParam[string](c, "id")

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid VM ID")
	}

	v, ok := h.vms.Delete(VMID)
	if !ok {
		return c.String(http.StatusNotFound, "no such VM (already deleted?)")
	}
	delete(v)
	return c.NoContent(http.StatusNoContent)

}

func (h *Controller) Submit(c *echo.Context) error {
	var req SubmissionRequest
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Malformed request was sent\n")
	}
	var instance *vm.VM
	select {
	case instance <- main.WarmVMPool:

	}
}

func delete(v *vm.VM) {
	v.Machine.StopVMM()
	os.Remove(v.Cfg.SocketPath)
	os.Remove(filepath.Join("/tmp/rootfs-" + v.UUID + ".ext4"))
}

func (h *Controller) create() (*vm.VM, error) {
	VMID := uuid.NewV4()
	ROOTFS := filepath.Join(os.Getenv("ROOTFS_PATH"), "rootfs.ext4")
	vmfs := "/tmp/rootfs-" + VMID.String() + ".ext4"
	tStart := time.Now()
	_, err := copyFS(ROOTFS, "/tmp/rootfs-"+VMID.String()+".ext4")
	timeTaken := time.Since(tStart)

	fmt.Printf("time taken to copy fs %v\n", timeTaken)
	if err != nil {
		return nil, err
	}

	cfg, err := vm.GetFirecrackerConfig(VMID.String())
	if err != nil {
		os.Remove(vmfs)
		return nil, err
	}
	// comment this for serial console(up)
	pwd, err := os.Getwd()

	stdoutPath := filepath.Join(pwd, "logs", VMID.String(), "stdout.log")

	stdout, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		os.Remove(vmfs)
		return nil, err
	}

	// stderr will be directed to this file
	stderrPath := filepath.Join(pwd, "logs", VMID.String(), "stderr.log")
	stderr, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		os.Remove(vmfs)
		return nil, err
	}
	// all above
	ctx := context.Background()

	BinPath := filepath.Join(os.Getenv("FIRECRACKER_BINARY_PATH"), "firecracker")
	fmt.Println("firecracker binary is at", BinPath)
	cmd := firecracker.VMCommandBuilder{}.
		WithBin(BinPath).
		WithSocketPath(cfg.SocketPath).
		WithStdout(stdout).
		WithStderr(stderr).
		Build(ctx)

	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))

	if err != nil {
		os.Remove(cfg.SocketPath)
		os.Remove(filepath.Join("/tmp/rootfs-" + VMID.String() + ".ext4"))
		return nil, err
	}
	if err := m.Start(ctx); err != nil {
		os.Remove(cfg.SocketPath)
		os.Remove(filepath.Join("/tmp/rootfs-" + VMID.String() + ".ext4"))
		return nil, err
	}
	ipAddr := m.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr.IP
	v := &vm.VM{
		UUID:    VMID.String(),
		IP:      ipAddr,
		Machine: m,
		Cfg:     cfg,
	}
	h.vms.Add(v)

	return v, nil
}

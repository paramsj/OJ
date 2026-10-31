package controllers

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
)

type SubmissionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Input    string `json:"input"`
}

type SubmissionResponse struct {
	Stdout    string  `json:"stdout"`
	Stderr    string  `json:"stderr"`
	TimeTaken float64 `json:"timetaken"`
	MemUsage  int64   `json:"memusage"`
}

func CreateSubmission(c *echo.Context) error {
	var req SubmissionRequest

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	tmpDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "main.cpp")
	execPath := filepath.Join(tmpDir, "main.out")

	if err := os.WriteFile(filePath, []byte(req.Code), 0644); err != nil {
		return c.String(http.StatusInternalServerError, "could not create the code file")
	}

	compileCmd := exec.Command("g++", filePath, "-o", execPath)

	if compileOut, err := compileCmd.CombinedOutput(); err != nil {
		return c.JSON(http.StatusInternalServerError, SubmissionResponse{
			Stderr: string(compileOut),
		})
	}

	var execStdOut, execStdErr bytes.Buffer
	runCmd := exec.Command(execPath)

	runCmd.Stdout = &execStdOut
	runCmd.Stderr = &execStdErr
	runCmd.Stdin = strings.NewReader(req.Input)

	start := time.Now()
	err = runCmd.Run()
	elapsedTime := time.Since(start)
	stderrStr := execStdErr.String()

	if err != nil {
		if stderrStr != "" {
			stderrStr = stderrStr + "\n" + err.Error()
		} else {
			stderrStr = err.Error()
		}
	}

	return c.JSON(http.StatusOK, SubmissionResponse{
		Stdout:    execStdOut.String(),
		Stderr:    stderrStr,
		TimeTaken: elapsedTime.Seconds(),
		MemUsage:  runCmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss,
	})
}

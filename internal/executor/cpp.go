package executor

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func ExecuteCpp(code string) (Result, error) {
	// Create temporary directory
	dir, err := os.MkdirTemp("", "oj-")
	if err != nil {
		return Result{}, err
	}
	defer os.RemoveAll(dir)

	// Write C++ source
	codePath := filepath.Join(dir, "main.cpp")

	if err := os.WriteFile(codePath, []byte(code), 0644); err != nil {
		return Result{}, err
	}

	// Compile
	binaryPath := filepath.Join(dir, "main")

	cmd := exec.Command(
		"g++",
		codePath,
		"-o",
		binaryPath,
	)

	compileOutput, err := cmd.CombinedOutput()

	if err != nil {
		return Result{
			Stdout: "",
			Stderr: string(compileOutput),
		}, nil
	}

	// Run
	cmd = exec.Command(binaryPath)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return Result{
			Stdout: "",
			Stderr: string(output),
		}, nil
	}

	return Result{
		Stdout: string(output),
		Stderr: "",
	}, nil
}

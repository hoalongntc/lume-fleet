package lume

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// DeleteVM shells out to `lume delete <name>`.
func DeleteVM(name string) error {
	cmd := exec.Command("lume", "delete", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lume delete %q: %s: %w", name, out, err)
	}
	return nil
}

// CloneVM shells out to `lume clone <source> <dest>`.
func CloneVM(source, dest string) error {
	cmd := exec.Command("lume", "clone", source, dest)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lume clone %q -> %q: %s: %w", source, dest, out, err)
	}
	return nil
}

// ListVMsViaCLI shells out to `lume ls --format json`.
func ListVMsViaCLI() ([]VM, error) {
	cmd := exec.Command("lume", "ls", "--format", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("lume ls --format json: %s: %w", out, err)
	}

	var vms []VM
	if err := json.Unmarshal(out, &vms); err != nil {
		return nil, fmt.Errorf("decode lume ls output: %w", err)
	}
	return vms, nil
}

// CreateVMViaCLI shells out to `lume create` with translated request options.
func CreateVMViaCLI(req CreateRequest) error {
	args := buildCreateCommandArgs(req)
	cmd := exec.Command("lume", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lume %v: %s: %w", args, out, err)
	}
	return nil
}

func buildCreateCommandArgs(req CreateRequest) []string {
	args := []string{
		"create", req.Name,
		"--os", req.OS,
		"--cpu", strconv.Itoa(req.CPU),
		"--memory", req.Memory,
		"--disk-size", req.DiskSize,
		"--display", req.Display,
	}
	if req.IPSW != "" {
		args = append(args, "--ipsw", req.IPSW)
	}
	if req.Unattended != "" {
		args = append(args, "--unattended", req.Unattended)
	}
	if req.VNCPort != 0 {
		args = append(args, "--vnc-port", strconv.Itoa(req.VNCPort))
	}
	if req.Storage != "" {
		args = append(args, "--storage", req.Storage)
	}
	if req.Network != "" {
		args = append(args, "--network", req.Network)
	}
	return args
}

// RunVMViaCLI shells out to `lume run <name> --no-display` with optional flags.
func RunVMViaCLI(name, sharedDir, mountISO string) error {
	args := buildRunCommandArgs(name, sharedDir, mountISO)
	cmd := exec.Command("lume", args...)

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("open %s: %w", os.DevNull, err)
	}
	defer devNull.Close()

	cmd.Stdout = devNull
	cmd.Stderr = devNull
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("lume %v: %w", args, err)
	}

	// `lume run` is long-running. Treat a still-running process as success and
	// only fail if it exits immediately with an error.
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case runErr := <-done:
		if runErr != nil {
			return fmt.Errorf("lume %v: %w", args, runErr)
		}
		return nil
	case <-time.After(500 * time.Millisecond):
		return nil
	}
}

func buildRunCommandArgs(name, sharedDir, mountISO string) []string {
	args := []string{"run", name, "--no-display"}
	if sharedDir != "" {
		args = append(args, "--shared-dir", sharedDir)
	}
	if mountISO != "" {
		args = append(args, "--mount", mountISO)
	}
	return args
}

// StopVMViaCLI shells out to `lume stop <name>`.
func StopVMViaCLI(name string) error {
	cmd := exec.Command("lume", "stop", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lume stop %q: %s: %w", name, out, err)
	}
	return nil
}

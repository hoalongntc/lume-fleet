package lume

import (
	"reflect"
	"strconv"
	"testing"
)

func TestBuildRunCommandArgsWithMount(t *testing.T) {
	args := buildRunCommandArgs("test-linux", "~/Projects", "/tmp/ubuntu.iso")

	want := []string{"run", "test-linux", "--no-display", "--shared-dir", "~/Projects", "--mount", "/tmp/ubuntu.iso"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("buildRunCommandArgs() = %v, want %v", args, want)
	}
}

func TestBuildRunCommandArgsWithoutOptionalFlags(t *testing.T) {
	args := buildRunCommandArgs("test-linux", "", "")

	want := []string{"run", "test-linux", "--no-display"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("buildRunCommandArgs() = %v, want %v", args, want)
	}
}

func TestBuildCreateCommandArgsForMacOS(t *testing.T) {
	req := CreateRequest{
		Name:       "mac-vm",
		OS:         "macos",
		CPU:        4,
		Memory:     "8GB",
		DiskSize:   "50GB",
		Display:    "1024x768",
		IPSW:       "/tmp/macos.ipsw",
		Unattended: "tahoe",
		VNCPort:    5901,
		Storage:    "default",
		Network:    "nat",
	}

	args := buildCreateCommandArgs(req)
	want := []string{
		"create", "mac-vm",
		"--os", "macos",
		"--cpu", strconv.Itoa(4),
		"--memory", "8GB",
		"--disk-size", "50GB",
		"--display", "1024x768",
		"--ipsw", "/tmp/macos.ipsw",
		"--unattended", "tahoe",
		"--vnc-port", strconv.Itoa(5901),
		"--storage", "default",
		"--network", "nat",
	}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("buildCreateCommandArgs() = %v, want %v", args, want)
	}
}

func TestBuildCreateCommandArgsForLinux(t *testing.T) {
	req := CreateRequest{
		Name:     "linux-vm",
		OS:       "linux",
		CPU:      2,
		Memory:   "4GB",
		DiskSize: "50GB",
		Display:  "1024x768",
	}

	args := buildCreateCommandArgs(req)
	want := []string{
		"create", "linux-vm",
		"--os", "linux",
		"--cpu", strconv.Itoa(2),
		"--memory", "4GB",
		"--disk-size", "50GB",
		"--display", "1024x768",
	}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("buildCreateCommandArgs() = %v, want %v", args, want)
	}
}

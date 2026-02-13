package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hoalong/lume-fleet/fleet"
)

func TestBuildCreateRequestIncludesVNCPort(t *testing.T) {
	vm := fleet.ResolvedVM{
		Name:       "dev-mac",
		OS:         "macos",
		CPU:        4,
		Memory:     "8GB",
		DiskSize:   "50GB",
		Unattended: "tahoe",
		VNCPort:    5999,
		Image:      "/tmp/macos.ipsw",
	}

	req := buildCreateRequest(vm)

	if req.VNCPort != 5999 {
		t.Fatalf("buildCreateRequest() vncPort = %d, want 5999", req.VNCPort)
	}
	if req.IPSW != "/tmp/macos.ipsw" {
		t.Fatalf("buildCreateRequest() ipsw = %q, want /tmp/macos.ipsw", req.IPSW)
	}

	payload, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal(create request): %v", err)
	}

	if !strings.Contains(string(payload), `"vncPort":5999`) {
		t.Fatalf("create request payload missing vncPort: %s", payload)
	}
}

func TestBuildCreateRequestDefaultsMacOSIPSWToLatest(t *testing.T) {
	vm := fleet.ResolvedVM{
		Name:     "dev-mac",
		OS:       "macos",
		CPU:      4,
		Memory:   "8GB",
		DiskSize: "50GB",
	}

	req := buildCreateRequest(vm)
	if req.IPSW != "latest" {
		t.Fatalf("buildCreateRequest() ipsw = %q, want latest", req.IPSW)
	}
}

func TestBuildRunRequestDoesNotIncludeVNCPort(t *testing.T) {
	vm := fleet.ResolvedVM{
		SharedDir: "/tmp/share",
		VNCPort:   5999,
	}

	req := buildRunRequest(vm)

	payload, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal(run request): %v", err)
	}

	if strings.Contains(string(payload), "vncPort") {
		t.Fatalf("run request unexpectedly includes vncPort: %s", payload)
	}

	if req.SharedDir != "/tmp/share" {
		t.Fatalf("run request sharedDir = %q, want /tmp/share", req.SharedDir)
	}
}

func TestShouldUseISOMountOnCreate(t *testing.T) {
	vm := fleet.ResolvedVM{
		OS:    "linux",
		Image: "/tmp/ubuntu.iso",
	}

	if !shouldUseISOMountOnCreate(vm, fleet.ActionCreate) {
		t.Fatalf("expected ISO mount to be used for linux create")
	}
	if shouldUseISOMountOnCreate(vm, fleet.ActionStart) {
		t.Fatalf("did not expect ISO mount for start action")
	}
	if shouldUseISOMountOnCreate(fleet.ResolvedVM{OS: "macos", Image: "/tmp/macos.iso"}, fleet.ActionCreate) {
		t.Fatalf("did not expect ISO mount for macOS")
	}
	if shouldUseISOMountOnCreate(fleet.ResolvedVM{OS: "linux"}, fleet.ActionCreate) {
		t.Fatalf("did not expect ISO mount when ISO path is empty")
	}
}

func TestRunVMForActionUsesCLIOnlyForLinuxCreateWithImage(t *testing.T) {
	origCLI := runVMViaCLI
	defer func() {
		runVMViaCLI = origCLI
	}()

	var cliCalled bool
	var cliMount string

	runVMViaCLI = func(name, sharedDir, mountISO string) error {
		cliCalled = true
		cliMount = mountISO
		return nil
	}

	vm := fleet.ResolvedVM{
		Name:  "linux-vm",
		OS:    "linux",
		Image: "/tmp/ubuntu.iso",
	}
	if err := runVMForAction(vm, fleet.ActionCreate); err != nil {
		t.Fatalf("runVMForAction() returned error: %v", err)
	}

	if !cliCalled {
		t.Fatalf("expected CLI path")
	}
	if cliMount != "/tmp/ubuntu.iso" {
		t.Fatalf("CLI mount = %q, want /tmp/ubuntu.iso", cliMount)
	}
}

func TestRunVMForActionUsesCLIForStartWithoutMount(t *testing.T) {
	origCLI := runVMViaCLI
	defer func() {
		runVMViaCLI = origCLI
	}()

	var cliCalled bool
	var cliMount string

	runVMViaCLI = func(name, sharedDir, mountISO string) error {
		cliCalled = true
		cliMount = mountISO
		return nil
	}

	vm := fleet.ResolvedVM{
		Name:  "linux-vm",
		OS:    "linux",
		Image: "/tmp/ubuntu.iso",
	}
	if err := runVMForAction(vm, fleet.ActionStart); err != nil {
		t.Fatalf("runVMForAction() returned error: %v", err)
	}

	if !cliCalled {
		t.Fatalf("expected CLI path")
	}
	if cliMount != "" {
		t.Fatalf("expected no mount on start, got %q", cliMount)
	}
}

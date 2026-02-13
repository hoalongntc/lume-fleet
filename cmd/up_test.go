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
	}

	req := buildCreateRequest(vm)

	if req.VNCPort != 5999 {
		t.Fatalf("buildCreateRequest() vncPort = %d, want 5999", req.VNCPort)
	}

	payload, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal(create request): %v", err)
	}

	if !strings.Contains(string(payload), `"vncPort":5999`) {
		t.Fatalf("create request payload missing vncPort: %s", payload)
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

package fleet

import (
	"strings"
	"testing"
)

func TestResolveInheritsAndOverridesVNCPort(t *testing.T) {
	cfg := FleetConfig{
		Defaults: VMDefaults{VNCPort: 6100},
		VMs: map[string]VMSpec{
			"mac-default": {
				OS: "macos",
			},
			"mac-override": {
				OS:      "macos",
				VNCPort: 6200,
			},
		},
	}

	resolved, err := cfg.Resolve()
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	gotByName := map[string]ResolvedVM{}
	for _, vm := range resolved {
		gotByName[vm.Name] = vm
	}

	if gotByName["mac-default"].VNCPort != 6100 {
		t.Fatalf("mac-default VNCPort = %d, want 6100", gotByName["mac-default"].VNCPort)
	}

	if gotByName["mac-override"].VNCPort != 6200 {
		t.Fatalf("mac-override VNCPort = %d, want 6200", gotByName["mac-override"].VNCPort)
	}
}

func TestResolveInheritsAndExpandsImagePath(t *testing.T) {
	cfg := FleetConfig{
		Defaults: VMDefaults{
			Image: "~/Downloads/base.iso",
		},
		VMs: map[string]VMSpec{
			"linux-default": {
				OS: "linux",
			},
			"linux-override": {
				OS:    "linux",
				Image: "~/Downloads/custom.iso",
			},
			"mac-override": {
				OS:    "macos",
				Image: "~/Downloads/macos.ipsw",
			},
		},
	}

	resolved, err := cfg.Resolve()
	if err != nil {
		t.Fatalf("Resolve() returned error: %v", err)
	}

	gotByName := map[string]ResolvedVM{}
	for _, vm := range resolved {
		gotByName[vm.Name] = vm
	}

	if !strings.Contains(gotByName["linux-default"].Image, "/Downloads/base.iso") {
		t.Fatalf("linux-default image = %q, want expanded ~/Downloads/base.iso", gotByName["linux-default"].Image)
	}

	if !strings.Contains(gotByName["linux-override"].Image, "/Downloads/custom.iso") {
		t.Fatalf("linux-override image = %q, want expanded ~/Downloads/custom.iso", gotByName["linux-override"].Image)
	}

	if !strings.Contains(gotByName["mac-override"].Image, "/Downloads/macos.ipsw") {
		t.Fatalf("mac-override image = %q, want expanded ~/Downloads/macos.ipsw", gotByName["mac-override"].Image)
	}
}

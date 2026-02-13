package fleet

import "testing"

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

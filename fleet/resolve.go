package fleet

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolvedVM is a VMSpec with defaults applied and the name attached.
type ResolvedVM struct {
	Name       string
	OS         string
	CPU        int
	Memory     string
	DiskSize   string
	SharedDir  string
	Unattended string
	VNCPort    int
	Image      string
	Storage    string
	Tags       []string
	Autostart  bool
}

// Resolve merges defaults into each VM spec and returns a sorted list.
func (c *FleetConfig) Resolve() ([]ResolvedVM, error) {
	var vms []ResolvedVM

	for name, spec := range c.VMs {
		vm := ResolvedVM{
			Name:      name,
			OS:        coalesce(spec.OS, c.Defaults.OS, "macos"),
			CPU:       coalesceInt(spec.CPU, c.Defaults.CPU, 4),
			Memory:    coalesce(spec.Memory, c.Defaults.Memory, "8GB"),
			DiskSize:  coalesce(spec.DiskSize, c.Defaults.DiskSize, "50GB"),
			SharedDir: expandHome(spec.SharedDir),
			VNCPort:   coalesceInt(spec.VNCPort, c.Defaults.VNCPort, 0),
			Image:     expandHome(coalesce(spec.Image, c.Defaults.Image, "")),
			Storage:   coalesce(spec.Storage, c.Defaults.Storage, ""),
			Tags:      spec.Tags,
			Autostart: true,
		}

		// Only apply unattended default for macOS VMs
		if strings.EqualFold(vm.OS, "macos") {
			vm.Unattended = coalesce(spec.Unattended, c.Defaults.Unattended, "")
		}

		if spec.Autostart != nil {
			vm.Autostart = *spec.Autostart
		}

		// Validate memory and disk-size are parseable
		if _, err := ParseSize(vm.Memory); err != nil {
			return nil, fmt.Errorf("VM %q: invalid memory: %w", name, err)
		}
		if _, err := ParseSize(vm.DiskSize); err != nil {
			return nil, fmt.Errorf("VM %q: invalid disk-size: %w", name, err)
		}
		if vm.VNCPort < 0 || vm.VNCPort > 65535 {
			return nil, fmt.Errorf("VM %q: invalid vnc-port %d (must be 0-65535)", name, vm.VNCPort)
		}

		vms = append(vms, vm)
	}

	return vms, nil
}

// FilterByNames returns only VMs whose names are in the given list.
func FilterByNames(vms []ResolvedVM, names []string) []ResolvedVM {
	if len(names) == 0 {
		return vms
	}
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	var result []ResolvedVM
	for _, vm := range vms {
		if set[vm.Name] {
			result = append(result, vm)
		}
	}
	return result
}

// FilterByTag returns only VMs that have the given tag.
func FilterByTag(vms []ResolvedVM, tag string) []ResolvedVM {
	if tag == "" {
		return vms
	}
	var result []ResolvedVM
	for _, vm := range vms {
		for _, t := range vm.Tags {
			if t == tag {
				result = append(result, vm)
				break
			}
		}
	}
	return result
}

func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func coalesceInt(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

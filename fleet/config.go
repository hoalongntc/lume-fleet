package fleet

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// FleetConfig is the top-level fleet.yml structure.
type FleetConfig struct {
	Defaults VMDefaults        `yaml:"defaults"`
	VMs      map[string]VMSpec `yaml:"vms"`
}

// VMDefaults provides default values inherited by all VMs.
type VMDefaults struct {
	OS         string `yaml:"os"`
	CPU        int    `yaml:"cpu"`
	Memory     string `yaml:"memory"`
	DiskSize   string `yaml:"disk-size"`
	Unattended string `yaml:"unattended"`
	VNCPort    int    `yaml:"vnc-port"`
	Storage    string `yaml:"storage"`
}

// VMSpec is one VM entry in the fleet.
type VMSpec struct {
	OS         string   `yaml:"os,omitempty"`
	CPU        int      `yaml:"cpu,omitempty"`
	Memory     string   `yaml:"memory,omitempty"`
	DiskSize   string   `yaml:"disk-size,omitempty"`
	SharedDir  string   `yaml:"shared-dir,omitempty"`
	Unattended string   `yaml:"unattended,omitempty"`
	VNCPort    int      `yaml:"vnc-port,omitempty"`
	Storage    string   `yaml:"storage,omitempty"`
	Tags       []string `yaml:"tags,omitempty"`
	Autostart  *bool    `yaml:"autostart,omitempty"`
}

// LoadConfig reads and parses a fleet.yml file.
func LoadConfig(path string) (*FleetConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("fleet: read config %q: %w", path, err)
	}

	var cfg FleetConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("fleet: parse config %q: %w", path, err)
	}

	if len(cfg.VMs) == 0 {
		return nil, fmt.Errorf("fleet: no VMs defined in %q", path)
	}

	return &cfg, nil
}

// ParseSize converts a human-readable size like "8GB" or "512MB" to megabytes.
func ParseSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	s = strings.ToUpper(s)

	var multiplier int64
	var numStr string

	switch {
	case strings.HasSuffix(s, "TB"):
		multiplier = 1024 * 1024
		numStr = strings.TrimSuffix(s, "TB")
	case strings.HasSuffix(s, "GB"):
		multiplier = 1024
		numStr = strings.TrimSuffix(s, "GB")
	case strings.HasSuffix(s, "MB"):
		multiplier = 1
		numStr = strings.TrimSuffix(s, "MB")
	default:
		return 0, fmt.Errorf("unsupported size unit in %q (use MB, GB, or TB)", s)
	}

	n, err := strconv.ParseFloat(strings.TrimSpace(numStr), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number in size %q: %w", s, err)
	}

	return int64(n * float64(multiplier)), nil
}

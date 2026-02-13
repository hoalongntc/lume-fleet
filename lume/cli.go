package lume

import (
	"fmt"
	"os/exec"
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

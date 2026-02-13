package fleet

import (
	"strings"

	"github.com/hoalong/lume-fleet/lume"
)

// ActionType describes what needs to happen to a VM.
type ActionType int

const (
	ActionCreate  ActionType = iota // VM doesn't exist -> create + run
	ActionStart                     // VM exists, stopped -> run
	ActionNoop                      // VM exists, running -> skip
	ActionStop                      // stop a running VM
	ActionDestroy                   // delete the VM entirely
)

// Action represents a reconciliation step.
type Action struct {
	VM      ResolvedVM
	Type    ActionType
	Current *lume.VM // nil if VM doesn't exist yet
}

// PlanUp compares desired VMs against actual Lume state and returns actions.
func PlanUp(desired []ResolvedVM, actual []lume.VM) []Action {
	index := indexByName(actual)
	var actions []Action

	for _, vm := range desired {
		if !vm.Autostart {
			continue
		}

		current, exists := index[vm.Name]
		switch {
		case !exists:
			actions = append(actions, Action{VM: vm, Type: ActionCreate})
		case strings.EqualFold(current.Status, "stopped"):
			actions = append(actions, Action{VM: vm, Type: ActionStart, Current: &current})
		case strings.EqualFold(current.Status, "running"):
			actions = append(actions, Action{VM: vm, Type: ActionNoop, Current: &current})
		default:
			// provisioning or other state — skip
			actions = append(actions, Action{VM: vm, Type: ActionNoop, Current: &current})
		}
	}

	return actions
}

// PlanDown returns stop actions for matching VMs that are running.
func PlanDown(desired []ResolvedVM, actual []lume.VM) []Action {
	index := indexByName(actual)
	var actions []Action

	for _, vm := range desired {
		current, exists := index[vm.Name]
		if exists && strings.EqualFold(current.Status, "running") {
			actions = append(actions, Action{VM: vm, Type: ActionStop, Current: &current})
		}
	}

	return actions
}

// PlanDestroy returns destroy actions for matching VMs.
func PlanDestroy(desired []ResolvedVM, actual []lume.VM) []Action {
	index := indexByName(actual)
	var actions []Action

	for _, vm := range desired {
		if _, exists := index[vm.Name]; exists {
			current := index[vm.Name]
			actions = append(actions, Action{VM: vm, Type: ActionDestroy, Current: &current})
		}
	}

	return actions
}

// CountRunningMacOS counts how many macOS VMs are currently running.
func CountRunningMacOS(actual []lume.VM) int {
	count := 0
	for _, vm := range actual {
		if strings.EqualFold(vm.OS, "macos") && strings.EqualFold(vm.Status, "running") {
			count++
		}
	}
	return count
}

func indexByName(vms []lume.VM) map[string]lume.VM {
	m := make(map[string]lume.VM, len(vms))
	for _, vm := range vms {
		m[vm.Name] = vm
	}
	return m
}

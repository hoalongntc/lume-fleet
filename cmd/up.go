package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hoalong/lume-fleet/fleet"
	"github.com/hoalong/lume-fleet/lume"
	"github.com/spf13/cobra"
)

var upTag string
var (
	runVMViaCLI = func(name, sharedDir, mountISO string) error {
		return lume.RunVMViaCLI(name, sharedDir, mountISO)
	}
)

var upCmd = &cobra.Command{
	Use:   "up [vm1 vm2 ...]",
	Short: "Create and start VMs defined in fleet.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := fleet.LoadConfig(cfgFile)
		if err != nil {
			return err
		}

		resolved, err := cfg.Resolve()
		if err != nil {
			return err
		}

		resolved = fleet.FilterByNames(resolved, args)
		resolved = fleet.FilterByTag(resolved, upTag)
		if len(resolved) == 0 {
			fmt.Println("No VMs match the given filters.")
			return nil
		}

		actual, err := lume.ListVMsViaCLI()
		if err != nil {
			return fmt.Errorf("cannot list VMs via lume CLI: %w", err)
		}

		actions := fleet.PlanUp(resolved, actual)
		macosRunning := fleet.CountRunningMacOS(actual)
		failures := 0

		for _, a := range actions {
			switch a.Type {
			case fleet.ActionNoop:
				fmt.Printf("[ ] %s: already running\n", a.VM.Name)

			case fleet.ActionStart:
				if strings.EqualFold(a.VM.OS, "macos") && macosRunning >= 2 {
					fmt.Fprintf(os.Stderr, "[!] %s: skipped — macOS 2-VM concurrent limit reached\n", a.VM.Name)
					failures++
					continue
				}
				fmt.Printf("[>] %s: starting...\n", a.VM.Name)
				err := runVMForAction(a.VM, fleet.ActionStart)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[x] %s: start failed: %v\n", a.VM.Name, err)
					failures++
					continue
				}
				if strings.EqualFold(a.VM.OS, "macos") {
					macosRunning++
				}
				fmt.Printf("[+] %s: running\n", a.VM.Name)

			case fleet.ActionCreate:
				if strings.EqualFold(a.VM.OS, "macos") && macosRunning >= 2 {
					fmt.Fprintf(os.Stderr, "[!] %s: skipped — macOS 2-VM concurrent limit reached\n", a.VM.Name)
					failures++
					continue
				}
				fmt.Printf("[>] %s: creating (this may take several minutes)...\n", a.VM.Name)

				createReq := buildCreateRequest(a.VM)

				if err := lume.CreateVMViaCLI(createReq); err != nil {
					fmt.Fprintf(os.Stderr, "[x] %s: create failed: %v\n", a.VM.Name, err)
					failures++
					continue
				}

				fmt.Printf("[>] %s: starting...\n", a.VM.Name)
				err := runVMForAction(a.VM, fleet.ActionCreate)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[x] %s: start failed: %v\n", a.VM.Name, err)
					failures++
					continue
				}
				if strings.EqualFold(a.VM.OS, "macos") {
					macosRunning++
				}
				fmt.Printf("[+] %s: running\n", a.VM.Name)
			}
		}

		if failures > 0 {
			return fmt.Errorf("%d VM(s) failed", failures)
		}
		return nil
	},
}

func init() {
	upCmd.Flags().StringVar(&upTag, "tag", "", "filter VMs by tag")
	rootCmd.AddCommand(upCmd)
}

func buildCreateRequest(vm fleet.ResolvedVM) lume.CreateRequest {
	req := lume.CreateRequest{
		Name:       vm.Name,
		OS:         vm.OS,
		CPU:        vm.CPU,
		Memory:     vm.Memory,
		DiskSize:   vm.DiskSize,
		Display:    "1024x768",
		Unattended: vm.Unattended,
		VNCPort:    vm.VNCPort,
	}
	if strings.EqualFold(vm.OS, "macos") {
		req.IPSW = vm.Image
		if req.IPSW == "" {
			req.IPSW = "latest"
		}
	}
	return req
}

func buildRunRequest(vm fleet.ResolvedVM) lume.RunRequest {
	return lume.RunRequest{
		NoDisplay: true,
		SharedDir: vm.SharedDir,
	}
}

func shouldUseISOMountOnCreate(vm fleet.ResolvedVM, actionType fleet.ActionType) bool {
	return actionType == fleet.ActionCreate && strings.EqualFold(vm.OS, "linux") && vm.Image != ""
}

func runVMForAction(vm fleet.ResolvedVM, actionType fleet.ActionType) error {
	mountISO := ""
	if shouldUseISOMountOnCreate(vm, actionType) {
		mountISO = vm.Image
	}
	return runVMViaCLI(vm.Name, vm.SharedDir, mountISO)
}

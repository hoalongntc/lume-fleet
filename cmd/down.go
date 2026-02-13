package cmd

import (
	"fmt"
	"os"

	"github.com/hoalong/lume-fleet/fleet"
	"github.com/hoalong/lume-fleet/lume"
	"github.com/spf13/cobra"
)

var downTag string

var downCmd = &cobra.Command{
	Use:   "down [vm1 vm2 ...]",
	Short: "Stop running VMs",
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
		resolved = fleet.FilterByTag(resolved, downTag)
		if len(resolved) == 0 {
			fmt.Println("No VMs match the given filters.")
			return nil
		}

		actual, err := lume.ListVMsViaCLI()
		if err != nil {
			return fmt.Errorf("cannot list VMs via lume CLI: %w", err)
		}

		actions := fleet.PlanDown(resolved, actual)
		if len(actions) == 0 {
			fmt.Println("No running VMs to stop.")
			return nil
		}

		failures := 0
		for _, a := range actions {
			fmt.Printf("[>] %s: stopping...\n", a.VM.Name)
			if err := lume.StopVMViaCLI(a.VM.Name); err != nil {
				fmt.Fprintf(os.Stderr, "[x] %s: stop failed: %v\n", a.VM.Name, err)
				failures++
				continue
			}
			fmt.Printf("[+] %s: stopped\n", a.VM.Name)
		}

		if failures > 0 {
			return fmt.Errorf("%d VM(s) failed to stop", failures)
		}
		return nil
	},
}

func init() {
	downCmd.Flags().StringVar(&downTag, "tag", "", "filter VMs by tag")
	rootCmd.AddCommand(downCmd)
}

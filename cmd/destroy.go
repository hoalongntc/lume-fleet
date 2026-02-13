package cmd

import (
	"fmt"
	"os"

	"github.com/hoalong/lume-fleet/fleet"
	"github.com/hoalong/lume-fleet/lume"
	"github.com/spf13/cobra"
)

var (
	destroyTag   string
	destroyForce bool
)

var destroyCmd = &cobra.Command{
	Use:   "destroy [vm1 vm2 ...]",
	Short: "Delete VMs entirely",
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
		resolved = fleet.FilterByTag(resolved, destroyTag)
		if len(resolved) == 0 {
			fmt.Println("No VMs match the given filters.")
			return nil
		}

		actual, err := lume.ListVMsViaCLI()
		if err != nil {
			return fmt.Errorf("cannot list VMs via lume CLI: %w", err)
		}

		actions := fleet.PlanDestroy(resolved, actual)
		if len(actions) == 0 {
			fmt.Println("No existing VMs to destroy.")
			return nil
		}

		if !destroyForce {
			fmt.Printf("About to destroy %d VM(s):\n", len(actions))
			for _, a := range actions {
				fmt.Printf("  - %s (%s)\n", a.VM.Name, a.Current.Status)
			}
			fmt.Print("\nThis is irreversible. Use --force to confirm.\n")
			return nil
		}

		failures := 0
		for _, a := range actions {
			// Stop running VMs before deleting
			if a.Current != nil && a.Current.Status == "running" {
				fmt.Printf("[>] %s: stopping before delete...\n", a.VM.Name)
				if err := lume.StopVMViaCLI(a.VM.Name); err != nil {
					fmt.Fprintf(os.Stderr, "[x] %s: stop failed: %v\n", a.VM.Name, err)
					failures++
					continue
				}
			}

			fmt.Printf("[>] %s: deleting...\n", a.VM.Name)
			if err := lume.DeleteVM(a.VM.Name); err != nil {
				fmt.Fprintf(os.Stderr, "[x] %s: delete failed: %v\n", a.VM.Name, err)
				failures++
				continue
			}
			fmt.Printf("[+] %s: deleted\n", a.VM.Name)
		}

		if failures > 0 {
			return fmt.Errorf("%d VM(s) failed to destroy", failures)
		}
		return nil
	},
}

func init() {
	destroyCmd.Flags().StringVar(&destroyTag, "tag", "", "filter VMs by tag")
	destroyCmd.Flags().BoolVar(&destroyForce, "force", false, "skip confirmation")
	rootCmd.AddCommand(destroyCmd)
}

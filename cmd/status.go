package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hoalong/lume-fleet/fleet"
	"github.com/hoalong/lume-fleet/lume"
	"github.com/hoalong/lume-fleet/ui"
	"github.com/spf13/cobra"
)

var (
	statusTag  string
	statusJSON bool
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show fleet VM status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := fleet.LoadConfig(cfgFile)
		if err != nil {
			return err
		}

		resolved, err := cfg.Resolve()
		if err != nil {
			return err
		}

		resolved = fleet.FilterByTag(resolved, statusTag)
		if len(resolved) == 0 {
			fmt.Println("No VMs match the given filters.")
			return nil
		}

		client := lume.NewClient("")
		actual, err := client.ListVMs()
		if err != nil {
			return fmt.Errorf("cannot reach Lume API at localhost:7777. Is 'lume serve' running?\n%w", err)
		}

		if statusJSON {
			return printJSON(resolved, actual)
		}

		rows := ui.BuildStatusRows(resolved, actual)
		macosRunning := fleet.CountRunningMacOS(actual)
		fmt.Println(ui.RenderStatusTable(rows, macosRunning))
		return nil
	},
}

func init() {
	statusCmd.Flags().StringVar(&statusTag, "tag", "", "filter VMs by tag")
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "output as JSON")
	rootCmd.AddCommand(statusCmd)
}

func printJSON(resolved []fleet.ResolvedVM, actual []lume.VM) error {
	rows := ui.BuildStatusRows(resolved, actual)
	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		return err
	}
	fmt.Println()
	return nil
}
